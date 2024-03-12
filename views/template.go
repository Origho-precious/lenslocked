package views

import (
	"bytes"
	"errors"
	"fmt"
	appcontext "github/Origho-precious/lenslocked/context"
	"github/Origho-precious/lenslocked/models"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/csrf"
)

type public interface {
	Public() string
}

type Template struct {
	htmlTpl *template.Template
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func ParseFS(fs fs.FS, patterns ...string) (Template, error) {
	// tpl := template.New(patterns[0])
	tpl := template.New(filepath.Base(patterns[0]))

	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() (template.HTML, error) {
				return "", fmt.Errorf("csrfField not implemented")
			},
			"currentUser": func() (*models.User, error) {
				return nil, fmt.Errorf("currentUser not implemented")
			},
			"errors": func() []string {
				return nil
			},
		},
	)

	tpl, err := tpl.ParseFS(fs, patterns...)

	if err != nil {
		return Template{}, fmt.Errorf("ParseFS:- parsing tempate: %w", err)
	}

	return Template{htmlTpl: tpl}, nil
}

// func Parse(path string) (Template, error) {
// 	tpl, err := template.ParseFiles(path)

// 	if err != nil {
// 		return Template{}, fmt.Errorf("parsing tempate: %w", err)
// 	}

// 	return Template{htmlTpl: tpl}, nil
// }

func (t Template) Execute(
	w http.ResponseWriter, r *http.Request, data any, errs ...error,
) {
	tpl, err := t.htmlTpl.Clone()

	if err != nil {
		log.Printf("Error cloning template: %v", err)
		http.Error(w,
			"Something went wrong while rendering page",
			http.StatusInternalServerError,
		)

		return
	}

	// Call the errMessages func before the closures.
	errMsgs := errMessages(errs...)

	tpl = tpl.Funcs(
		template.FuncMap{
			"csrfField": func() template.HTML {
				return csrf.TemplateField(r)
			},
			"currentUser": func() *models.User {
				return appcontext.User(r.Context())
			},
			"errors": func() []string {
				return errMsgs
			},
		},
	)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var buf bytes.Buffer

	err = tpl.Execute(&buf, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(
			w,
			"Something went wrong while executing the template",
			http.StatusInternalServerError,
		)
	}

	io.Copy(w, &buf)
}

func errMessages(errs ...error) []string {
	var msgs []string
	for _, err := range errs {
		var pubErr public
		if errors.As(err, &pubErr) {
			msgs = append(msgs, pubErr.Public())
		} else {
			fmt.Println(err)
			msgs = append(msgs, "Something went wrong.")
		}
	}
	return msgs
}
