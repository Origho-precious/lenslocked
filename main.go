package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, "<h1>Welcome to my awesome site!</h1>")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
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

func faqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `
		<!DOCTYPE html>
		<html lang="en">
		<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>FAQ</title>
				<style>
						body {
								font-family: Arial, sans-serif;
								margin: 20px;
						}

						h1 {
							font-size: 20px;
						}

						.faq-item {
								margin-bottom: 20px;
						}

						.question {
							font-size: 18px;
								font-weight: bold;
						}

						.answer {
							font-size: 16px;
							margin-top: 5px;
						}
				</style>
		</head>
		<body>
		<h1>FAQ Page</h1>

		<div class="faq-item">
				<p class="question">Is there a free version?</p>
				<p class="answer">Yes! We offer a free trial for 30 days on any paid plan.</p>
		</div>

		<div class="faq-item">
				<p class="question">What are your support hours?</p>
				<p class="answer">We have a support staff answering 24/7, though response times may be a bit slower on weekends.</p>
		</div>

		<div class="faq-item">
				<p class="question">How do I contact support?</p>
				<p class="answer">Email us - <a href="mailto:support@lenslocked.com">support@lenslocked.com</a></p>
		</div>

		</body>
		</html>
	`)
}

func singleFaqHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	faqID := chi.URLParam(r, "id")

	fmt.Fprintf(w, "<p>The FAQ Id: %v</p>", faqID)
}

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Get("/", homeHandler)
	r.Get("/contact", contactHandler)
	r.Get("/faq", faqHandler)
	r.Get("/faq/{id}", singleFaqHandler)
	r.NotFound(notFoundHandler)

	fmt.Println("Starting the server on :5500...")
	http.ListenAndServe(":5500", r)
}
