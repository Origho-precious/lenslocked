package models

import "database/sql"

type User struct {
	Id           uint
	Email        string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
}
