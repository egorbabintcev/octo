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
WHERE id = @webhookID;

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
WHERE id = @webhookID