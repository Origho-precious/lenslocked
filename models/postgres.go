package models

import (
	"database/sql"
	"fmt"
	"io/fs"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DbName   string
	SSLMode  string
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName, cfg.SSLMode,
	)
}

func DefaultPostgresConfig() PostgresConfig {
	config := PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "pilolo",
		Password: "sredev",
		DbName:   "lenslocked",
		SSLMode:  "disable",
	}

	return config
}

func Open(config PostgresConfig) (*sql.DB, error) {
	db, err := sql.Open("pgx", config.String())

	if err != nil {
		return nil, fmt.Errorf("opening SQL DB connection: %w", err)
	}

	return db, nil
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}

func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
	// In case the dir is an empty string, they probably meant the current directory and goose wants a period for that.
	if dir == "" {
		dir = "."
	}
	goose.SetBaseFS(migrationsFS)
	defer func() {
		// Ensure that we remove the FS on the off chance some other part of our app uses goose for migrations and doesn't want to use our FS.
		goose.SetBaseFS(nil)
	}()
	return Migrate(db, dir)
}
