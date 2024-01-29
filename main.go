package main

import (
	"fmt"
	"net/http"
)

func homeHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
}

func contactHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Contact Page</h1><p>To get in touch, email me at <a href=\"mailto:origho9@gmail.com\">origho9@gmail.com</a>.</p>")
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

type Router struct{}

func (router Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		homeHandlerFunc(w, r)
	case "/contact":
		contactHandlerFunc(w, r)
	default:
		notFoundHandler(w, r)
	}
}

func main() {
	var router Router

	fmt.Println("Starting the server on :5500...")
	http.ListenAndServe(":5500", router)
}
