package database

import (
	"database/sql"
	"fmt"

	"gamification-service/migrations"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func NewPostgres(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	if err := runMigrations(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return db, nil
}

func runMigrations(db *sql.DB) error {
	// Читаем именно init.sql (как у тебя на скриншоте)
	content, err := migrations.FS.ReadFile("init.sql")
	if err != nil {
		return fmt.Errorf("read migration file: %w", err)
	}

	if _, err := db.Exec(string(content)); err != nil {
		return fmt.Errorf("execute migration: %w", err)
	}

	return nil
}
