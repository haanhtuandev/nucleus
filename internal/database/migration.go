package database

import (
	"database/sql"
	"embed"
	"io/fs"

	"github.com/pressly/goose/v3"
)

//go:embed sql/schema/*.sql
var embedMigrations embed.FS

func RunMigrations(db *sql.DB) error {
	// Get the sub-directory containing migrations
	migrationsFS, err := fs.Sub(embedMigrations, "sql/schema")
	if err != nil {
		return err
	}

	goose.SetBaseFS(migrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	return goose.Up(db, ".")
}
