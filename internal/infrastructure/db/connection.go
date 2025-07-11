package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
)

var (
	ErrOpen = errors.New("failed to open database connection")
)

type Connection struct {
	logger *slog.Logger
	DB     *sql.DB
}

func NewConnection(l *slog.Logger, dsn string) (*Connection, error) {
	l = l.With(slog.String("component", "database"))

	l.Info("Attempting to open database connection")
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrOpen, err)
	}

	return &Connection{
		logger: l,
		DB:     db,
	}, nil
}

func (c *Connection) Close() {
	c.logger.Info("Attempting to close database connection")

	if err := c.DB.Close(); err != nil {
		c.logger.Error("Failed to close database connection")
	}
}
