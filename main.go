package main

import (
	"fmt"
	"github/Origho-precious/lenslocked/controllers"
	"github/Origho-precious/lenslocked/migrations"
	"github/Origho-precious/lenslocked/models"
	"github/Origho-precious/lenslocked/templates"
	"github/Origho-precious/lenslocked/views"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
)

type config struct {
	PSQL models.PostgresConfig
	SMTP models.SMTPConfig
	CSRF struct {
		Key    string
		Secure bool
	}
	Server struct {
		Address string
	}
}

func loadEnvConfig() (config, error) {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		return cfg, err
	}
	cfg.PSQL = models.PostgresConfig{
		Host:     os.Getenv("PSQL_HOST"),
		Port:     os.Getenv("PSQL_PORT"),
		User:     os.Getenv("PSQL_USER"),
		SSLMode:  os.Getenv("PSQL_SSLMODE"),
		Password: os.Getenv("PSQL_PASSWORD"),
		DbName:   os.Getenv("PSQL_DATABASE"),
	}
	if cfg.PSQL.Host == "" && cfg.PSQL.Port == "" {
		return cfg, fmt.Errorf("no PSQL Config provided")
	}

	cfg.SMTP.Host = os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")

	cfg.SMTP.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return cfg, err
	}
	cfg.SMTP.Username = os.Getenv("SMTP_USERNAME")
	cfg.SMTP.Password = os.Getenv("SMTP_PASSWORD")

	cfg.CSRF.Key = os.Getenv("CSRF_KEY")
	cfg.CSRF.Secure = os.Getenv("CSRF_SECURE") == "true"

	cfg.Server.Address = os.Getenv("SERVER_ADDRESS")

	return cfg, nil
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

func main() {
	cfg, err := loadEnvConfig()
	if err != nil {
		panic(err)
	}

	// Database config
	db, err := models.Open(cfg.PSQL)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = models.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}

	// Services
	userService := &models.UserService{
		DB: db,
	}

	sessionService := &models.SessionService{
		DB: db,
	}

	passwordResetService := &models.PasswordResetService{
		DB: db,
	}

	emailService := models.NewEmailService(cfg.SMTP)

	gallaryService := &models.GalleryService{
		DB: db,
	}

	// Middlewares
	userMiddleware := controllers.UserMiddleware{
		SessionService: sessionService,
	}

	csrfMiddleware := csrf.Protect(
		[]byte(cfg.CSRF.Key),
		csrf.Secure(cfg.CSRF.Secure),
		csrf.Path("/"),
	)

	// User controller
	usersController := controllers.Users{
		UserService:          userService,
		EmailService:         emailService,
		SessionService:       sessionService,
		PasswordResetService: passwordResetService,
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

	usersController.Templates.CheckYourEmail = views.Must(
		views.ParseFS(templates.FS, "check-your-email.gohtml", "tailwind.gohtml"),
	)

	usersController.Templates.ResetPassword = views.Must(
		views.ParseFS(templates.FS, "reset-pw.gohtml", "tailwind.gohtml"),
	)

	// Gallery Controller
	galleryController := controllers.Galleries{
		GalleryService: gallaryService,
	}

	galleryController.Templates.New = views.Must(
		views.ParseFS(templates.FS, "galleries/new.gohtml", "tailwind.gohtml"),
	)

	galleryController.Templates.Edit = views.Must(
		views.ParseFS(templates.FS, "galleries/edit.gohtml", "tailwind.gohtml"),
	)

	galleryController.Templates.Index = views.Must(
		views.ParseFS(templates.FS, "galleries/index.gohtml", "tailwind.gohtml"),
	)

	galleryController.Templates.Show = views.Must(
		views.ParseFS(templates.FS, "galleries/show.gohtml", "tailwind.gohtml"),
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

	r.Route("/users/me", func(r chi.Router) {
		r.Use(userMiddleware.RequireUser)
		r.Get("/", usersController.CurrentUser)
	})

	r.Post("/signout", usersController.ProcessSignOut)

	r.Get("/forgot-pw", usersController.ForgotPassword)
	r.Post("/forgot-pw", usersController.ProcessForgotPassword)
	r.Get("/reset-pw", usersController.ResetPassword)
	r.Post("/reset-pw", usersController.ProcessResetPassword)

	r.Route("/galleries", func(r chi.Router) {
		r.Get("/{id}", galleryController.Show)
		r.Group(func(r chi.Router) {
			r.Use(userMiddleware.RequireUser)
			r.Get("/", galleryController.Index)
			r.Get("/new", galleryController.New)
			r.Post("/", galleryController.Create)
			r.Post("/{id}", galleryController.Update)
			r.Get("/{id}/edit", galleryController.Edit)
			r.Post("/{id}/delete", galleryController.Delete)
			r.Post("/{id}/images", galleryController.UploadImage)
			r.Get("/{id}/images/{filename}", galleryController.Image)
			r.Post("/{id}/images/{filename}/delete", galleryController.DeleteImage)
		})
	})

	r.NotFound(notFoundHandler)

	// Start server
	fmt.Printf("Starting the server on %s...", cfg.Server.Address)

	err = http.ListenAndServe(cfg.Server.Address, r)
	if err != nil {
		panic(err)
	}
}
