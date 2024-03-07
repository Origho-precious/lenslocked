package main

import (
	"fmt"
	"github/Origho-precious/lenslocked/controllers"
	"github/Origho-precious/lenslocked/migrations"
	"github/Origho-precious/lenslocked/models"
	"github/Origho-precious/lenslocked/templates"
	"github/Origho-precious/lenslocked/views"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
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
	// Database config
	config := models.DefaultPostgresConfig()
	db, err := models.Open(config)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Services
	userService := models.UserService{
		DB: db,
	}

	sessionService := models.SessionService{
		DB: db,
	}

	// Middlewares
	userMiddleware := controllers.UserMiddleware{
		SessionService: &sessionService,
	}

	csrfAuthKey := "gFvi45R4fy5xNBlnEeZtQbfAVCYEIAUX"

	csrfMiddleware := csrf.Protect(
		[]byte(csrfAuthKey),
		csrf.Secure(false), // TODO: update this before deploying
	)

	// User controller
	usersController := controllers.Users{
		UserService:    &userService,
		SessionService: &sessionService,
	}

	usersController.Templates.New = views.Must(
		views.ParseFS(templates.FS, "signup.gohtml", "tailwind.gohtml"),
	)

	usersController.Templates.Signin = views.Must(
		views.ParseFS(templates.FS, "signin.gohtml", "tailwind.gohtml"),
	)

	usersController.Templates.ForgotPassword = views.Must(
		views.ParseFS(templates.FS, "forgot-pw.gohtml", "tailwind.gohtml"),
	)

	// Router and Routes
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(csrfMiddleware)
	r.Use(userMiddleware.SetUser)

	r.Get("/", controllers.StaticHandler(
		views.Must(views.ParseFS(templates.FS, "home.gohtml", "tailwind.gohtml")),
	))

	r.Get("/contact", controllers.StaticHandler(
		views.Must(
			views.ParseFS(templates.FS, "contact.gohtml", "tailwind.gohtml"),
		),
	))

	r.Get("/faq", controllers.FAQ(
		views.Must(views.ParseFS(templates.FS, "faq.gohtml", "tailwind.gohtml")),
	))

	r.Get("/signup", usersController.New)
	r.Post("/users", usersController.Create)
	r.Get("/signin", usersController.Signin)
	r.Post("/signin", usersController.ProcessSignin)
	r.Get("/forgot-pw", usersController.ForgotPassword)
	r.Post("/forgot-pw", usersController.ProcessForgotPassword)

	r.Route("/users/me", func(r chi.Router) {
		r.Use(userMiddleware.RequireUser)
		r.Get("/", usersController.CurrentUser)
	})

	r.Post("/signout", usersController.ProcessSignOut)

	r.NotFound(notFoundHandler)

	// Start server
	fmt.Println("Starting the server on :5500...")
	http.ListenAndServe(":5500", r)
}
