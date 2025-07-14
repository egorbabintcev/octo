-- name: GetAllWebhooks :many
SELECT
  event_type,
  context_type,
  context_id,
  endpoint,
  description
FROM webhooks;

-- name: GetWebhookByID :one
SELECT
  event_type,
  context_type,
  context_id,
  endpoint,
  description
FROM webhooks
WHERE id = @ID;

-- name: GetWebhookByParams :one
SELECT
  event_type,
  context_type,
  context_id,
  endpoint,
  description
FROM webhooks WHERE
  event_type = @eventType AND
  context_type = @contextType AND
  context_id = @contextID AND
  endpoint = @endpoint;

-- name: CreateWebhook :one
INSERT INTO webhooks (
  event_type,
  context_type,
  context_id,
  endpoint,
  passcode,
  description
) VALUES (
  @eventType,
  @contextType,
  @contextID,
  @endpoing,
  @passcode,
  sqlc.narg('description')
) RETURNING "id";

-- name: DeleteWebhook :exec
DELETE from webhooks
WHERE id = @ID