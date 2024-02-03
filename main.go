package main

import (
	"fmt"
	"github/Origho-precious/lenslocked/views"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func executeTemplate(w http.ResponseWriter, path string, data interface{}) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tpl, err := views.Parse(path, data)
	if err != nil {
		log.Printf("%v", err)
		http.Error(
			w,
			"Something went wrong while parsing the template",
			http.StatusInternalServerError,
		)

		return
	}

	tpl.Execute(w, data)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("templates", "home.gohtml")

	executeTemplate(w, path, nil)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("templates", "contact.gohtml")

	executeTemplate(w, path, nil)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// w.WriteHeader(http.StatusNotFound)
	// fmt.Fprint(w, "<h1>Page not found</h1>")
	// OR
	// http.Error(w, "Page not found", http.StatusNotFound)
	// OR
	// http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	http.Error(w, http.StatusText(404), http.StatusNotFound) // Not found
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("templates", "faqs.gohtml")

	executeTemplate(w, path, nil)
}

func singleFaqHandler(w http.ResponseWriter, r *http.Request) {
	faqID := chi.URLParam(r, "id")

	path := filepath.Join("templates", "faq.gohtml")

	data := struct {
		PageTitle string
		Id        string
	}{
		PageTitle: "FAQ title",
		Id:        faqID,
	}

	executeTemplate(w, path, data)

}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faqs", faqHandler)
	r.Get("/faqs/{id}", singleFaqHandler)
	r.NotFound(notFoundHandler)

	fmt.Println("Starting the server on :5500...")
	http.ListenAndServe(":5500", r)
}
