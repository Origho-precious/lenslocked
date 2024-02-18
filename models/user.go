package models

import (
	"database/sql"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           uint
	Email        string
	PasswordHash string
}

type UserService struct {
	DB *sql.DB
}

func (us *UserService) Create(email, password string) (*User, error) {
	email = strings.ToLower(email)

	hashByte, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost,
	)

	if err != nil {
		return nil, fmt.Errorf("creating new user: %w", err)
	}

	passwordHash := string(hashByte)

	row := us.DB.QueryRow(`
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2) RETURNING id
	`, email, passwordHash,
	)

	newUser := User{
		Email:        email,
		PasswordHash: passwordHash,
	}

	err = row.Scan(&newUser.Id)

	if err != nil {
		return nil, fmt.Errorf("creating new user: %w", err)
	}

	return &newUser, nil
}
