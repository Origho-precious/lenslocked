package main

import (
	"fmt"
	"github/Origho-precious/lenslocked/models"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	// email := models.Email{
	// 	From:      "test@lenslocked.com",
	// 	To:        "jon@calhoun.io",
	// 	Subject:   "This is a test email",
	// 	Plaintext: "This is the body of the email",
	// 	HTML:      `<h1>Hello there buddy!</h1><p>This is the email</p><p>Hope you enjoy it</p>`,
	// }

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic(err)
	}

	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")

	es := models.NewEmailService(models.SMTPConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	})

	err = es.ForgotPassword("jon@calhoun.io", "https://lenslocked.com/reset-pw?token=abc123")

	if err != nil {
		panic(err)
	}

	fmt.Println("Email sent!")
}
