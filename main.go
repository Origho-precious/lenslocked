package main

import (
	"fmt"
	"github/Origho-precious/lenslocked/controllers"
	"github/Origho-precious/lenslocked/templates"
	"github/Origho-precious/lenslocked/views"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

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

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// Home route
	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "home.gohtml")),
	))

	// Contact route
	r.Get("/contact", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "contact.gohtml")),
	))

	// FAQ route
	r.Get("/faqs", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "faqs.gohtml", "faq.gohtml")),
	))

	r.NotFound(notFoundHandler)

	fmt.Println("Starting the server on :5500...")
	http.ListenAndServe(":5500", r)
}
