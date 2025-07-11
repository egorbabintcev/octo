-- +goose Up
-- +goose StatementBegin
CREATE TABLE webhooks (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  event_type TEXT NOT NULL,
  context_type TEXT NOT NULL,
  context_id TEXT NOT NULL,
  endpoint TEXT NOT NULL,
  passcode BLOB NOT NULL,
  description TEXT,
  UNIQUE (event_type, context_type, context_id, endpoint),
  UNIQUE (passcode)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE webhooks;
-- +goose StatementEnd
