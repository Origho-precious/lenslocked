package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           int
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
		// See if we can use this error as a PgError
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			// This is a PgError, so see if it matches a unique violation.
			if pgError.Code == pgerrcode.UniqueViolation {
				// If this is true, it has to be an email violation since this is the
				// only way to trigger this type of violation with our SQL.
				return nil, ErrEmailTaken
			}
		}
		return nil, fmt.Errorf("create user: %w", err)
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

func (us *UserService) UpdatePassword(userID int, password string) error {
	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost,
	)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	passwordHash := string(hashedBytes)
	_, err = us.DB.Exec(`
		UPDATE users
		SET password_hash = $2
		WHERE id = $1;`, userID, passwordHash,
	)

	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	return nil
}
