package views

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type Template struct {
	htmlTpl *template.Template
}

func Must(t Template, err error) Template {
	if err != nil {
		panic(err)
	}
	return t
}

func Parse(path string) (Template, error) {
	footerPath := filepath.Join("templates", "footer.gohtml")

	tpl, err := template.ParseFiles(path, footerPath)

	if err != nil {
		return Template{}, fmt.Errorf("parsing tempate: %w", err)
	}

	return Template{htmlTpl: tpl}, nil
}

func (t Template) Execute(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := t.htmlTpl.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(
			w,
			"Something went wrong while executing the template",
			http.StatusInternalServerError,
		)
	}
}
