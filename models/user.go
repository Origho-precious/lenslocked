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

func (us UserService) Authenticate(email, password string) (*User, error) {
	email = strings.ToLower(email)

	// Get user
	row := us.DB.QueryRow(`
		SELECT id, password_hash FROM users
		WHERE email = $1`, email,
	)

	user := User{
		Email: email,
	}

	err := row.Scan(&user.Id, &user.PasswordHash)

	if err != nil {
		return nil, fmt.Errorf("authenticating user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash), []byte(password),
	)

	if err != nil {
		return nil, fmt.Errorf("authenticating user: %w", err)
	}

	return &user, nil
}
