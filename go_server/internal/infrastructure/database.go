package infrastructure

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func NewDBConnection() (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable not set")
	}

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func MigrateDB(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS builds (
			upload_id TEXT PRIMARY KEY,
			bundle_id TEXT NOT NULL,
			version TEXT NOT NULL,
			build_number TEXT NOT NULL,
			title TEXT NOT NULL,
			icon TEXT,
			description TEXT,
			file_size BIGINT NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL,
			platform TEXT NOT NULL
		)
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	return nil
}
