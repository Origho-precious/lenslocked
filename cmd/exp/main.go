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
	// _, err = db.Exec(`
	// 	CREATE TABLE IF NOT EXISTS users (
	// 		id SERIAL PRIMARY KEY,
	// 		name TEXT,
	// 		email TEXT UNIQUE NOT NULL
	// 	);

	// 	CREATE TABLE IF NOT EXISTS orders (
	// 		id SERIAL PRIMARY KEY,
	// 		user_id INT NOT NULL,
	// 		amount INT,
	// 		description TEXT
	// 	);
	// `)

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Tables created!")

	// Create records
	// name := "Bruce Wayne"
	// email := "bruce@wayne.com"

	// row := db.QueryRow(`
	// 	INSERT INTO users (name, email)
	// 	VALUES ($1, $2) RETURNING id;`,
	// 	name, email,
	// )

	// var id int

	// err = row.Scan(&id)

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("id:", id)

	// fmt.Println("Records added!")

	// Querying a record
	// id := 1

	// row := db.QueryRow(`
	// 	SELECT name, email
	// 	FROM users
	// 	WHERE id=$1;`, id,
	// )

	// var name, email string

	// err = row.Scan(&name, &email)

	// if err == sql.ErrNoRows {
	// 	fmt.Println("Error, no rows!")
	// }

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("User: name=%s, email=%s\n", name, email)

	// Creating associated data

	userId := 1

	for i := 1; i <= 5; i++ {
		amount := i * 100
		desc := fmt.Sprintf("Fake order #%d", i)

		_, err := db.Exec(`
			INSERT INTO orders(user_id, amount, description)
			VALUES ($1, $2, $3);`, userId, amount, desc)

		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Created fake orders!")
}
