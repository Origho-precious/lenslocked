package main

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgressConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
	SSLMode  string
}

func (cfg PostgressConfig) String() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName, cfg.SSLMode,
	)
}

func main() {
	cfg := PostgressConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "pilolo",
		Password: "sredev",
		DbName:   "lenslocked",
		SSLMode:  "disable",
	}

	db, err := sql.Open("pgx", cfg.String())

	if err != nil {
		panic(err)
	}

	defer db.Close()
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Database connection established!")

	// Create tables
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT,
			email TEXT UNIQUE NOT NULL
		);

		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			amount INT,
			description TEXT
		);
	`)

	if err != nil {
		panic(err)
	}

	fmt.Println("Tables created!")

	// Create records
	name := "Bruce Wayne"
	email := "bruce@wayne.com"

	_, err = db.Exec(`INSERT INTO users (name, email)  VALUES ($1, $2);`,
		name, email,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("Records added!")
}
