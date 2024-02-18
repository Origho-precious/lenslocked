package models

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
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
