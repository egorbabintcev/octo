package webhooks

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	database "octo/internal/infrastructure/db/gen"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type Storage struct {
	logger  *slog.Logger
	db      *sql.DB
	queries *database.Queries
}

func NewStorage(l *slog.Logger, db *sql.DB) *Storage {
	return &Storage{
		logger:  l,
		db:      db,
		queries: database.New(db),
	}
}

func (s *Storage) GetAll(ctx context.Context) (GetAllResponse, error) {
	rows, err := s.queries.GetAllWebhooks(ctx)
	if err != nil {
		return GetAllResponse{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	res := GetAllResponse{Webhooks: make([]Webhook, 0, len(rows))}
	for _, row := range rows {
		res.Webhooks = append(res.Webhooks, Webhook{
			EventType:   row.EventType,
			ContextType: row.ContextType,
			ContextID:   row.ContextID,
			Endpoint:    row.Endpoint,
			Description: row.Description.String,
		})
	}

	return res, nil
}

func (s *Storage) GetByID(ctx context.Context, id int) (GetByIdResponse, error) {
	row, err := s.queries.GetWebhookByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return GetByIdResponse{}, ErrNotFound
		}

		return GetByIdResponse{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	w := Webhook{
		EventType:   row.EventType,
		ContextType: row.ContextType,
		ContextID:   row.ContextID,
		Endpoint:    row.Endpoint,
		Description: row.Description.String,
	}

	return GetByIdResponse{
		Webhook: w,
	}, nil
}

func (s *Storage) GetByParams(ctx context.Context, eventType, contextType, contextID, endpoint string) (GetByParamsResponse, error) {
	row, err := s.queries.GetWebhookByParams(ctx, database.GetWebhookByParamsParams{
		EventType:   eventType,
		ContextType: contextType,
		ContextID:   contextID,
		Endpoint:    endpoint,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return GetByParamsResponse{}, ErrNotFound
		}

		return GetByParamsResponse{}, fmt.Errorf("%w: %v", ErrInternal, err)
	}

	w := Webhook{
		EventType:   row.EventType,
		ContextType: row.ContextType,
		ContextID:   row.ContextID,
		Endpoint:    row.Endpoint,
		Description: row.Description.String,
	}

	return GetByParamsResponse{
		Webhook: w,
	}, nil
}

func (s *Storage) Create(ctx context.Context, eventType, contextType, contextID, endpoint string, passcode []byte, description string) (CreateResponse, error) {
	id, err := s.queries.CreateWebhook(ctx, database.CreateWebhookParams{
		EventType:   eventType,
		ContextType: contextType,
		ContextID:   contextID,
		Endpoing:    endpoint,
		Passcode:    passcode,
		Description: sql.NullString{
			Valid:  description != "",
			String: description,
		},
	})
	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) {
		if sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return CreateResponse{}, fmt.Errorf("%w: %v", ErrConflict, sqliteErr.Error())
		}
	} else if err != nil {
		return CreateResponse{}, fmt.Errorf("%w: %v", ErrUknown, err)
	}

	return CreateResponse{
		ID: id,
	}, nil
}

func (s *Storage) Delete(ctx context.Context, id int) (DeleteResponse, error) {
	if _, err := s.GetByID(ctx, id); errors.Is(err, ErrNotFound) {
		return DeleteResponse{}, ErrNotFound
	}

	if err := s.queries.DeleteWebhook(ctx, id); err != nil {
		return DeleteResponse{}, err
	}

	return DeleteResponse{}, nil
}
