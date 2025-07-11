package database

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/pressly/goose/v3"
)

var (
	ErrMigrate = errors.New("failed to migrate database")
)

//go:embed migrations
var migrations embed.FS

func Migrate(db *sql.DB) error {
	goose.SetBaseFS(migrations)

	if err := goose.SetDialect("sqlite"); err != nil {
		return fmt.Errorf("%w: %v", ErrMigrate, err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return fmt.Errorf("%w: %v", ErrMigrate, err)
	}

	return nil
}
