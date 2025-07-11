package figma

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

const BASE_URL = "https://api.figma.com/v2"

type Client struct {
	logger     *slog.Logger
	httpClient http.Client
	token      string
}

func NewClient(l *slog.Logger, token string) *Client {
	return &Client{
		logger:     l,
		httpClient: *http.DefaultClient,
		token:      token,
	}
}

func (c *Client) CreateWebhook(ctx context.Context, r CreateWebhookRequest) (CreateWebhookResponse, error) {
	reqJson, err := json.Marshal(r)
	if err != nil {
		return CreateWebhookResponse{}, fmt.Errorf("%w: %v", ErrUnknown, err)
	}

	reqBody := bytes.NewReader(reqJson)
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/webhooks", BASE_URL), reqBody)
	if err != nil {
		return CreateWebhookResponse{}, fmt.Errorf("%w: %v", ErrUnknown, err)
	}

	req.Header.Set("X-FIGMA-TOKEN", c.token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return CreateWebhookResponse{}, fmt.Errorf("%w: %v", ErrUnknown, err)
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return CreateWebhookResponse{}, fmt.Errorf("%w: %v", ErrUnknown, err)
	}

	if res.StatusCode == http.StatusOK {
		var wh CreateWebhookResponse
		if err := json.Unmarshal(resBody, &wh); err != nil {
			return CreateWebhookResponse{}, nil
		}

		return wh, nil
	}

	var errRes errResponse
	if err := json.Unmarshal(resBody, &errRes); err != nil {
		return CreateWebhookResponse{}, fmt.Errorf("%w: %v", ErrUnknown, err)
	}

	switch res.StatusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		return CreateWebhookResponse{}, fmt.Errorf("%w: %v", ErrInvalidCredentials, errRes.Message)
	case http.StatusBadRequest:
		return CreateWebhookResponse{}, fmt.Errorf("%w: %v", ErrInvalidRequest, errRes.Message)
	default:
		return CreateWebhookResponse{}, ErrUnknown
	}
}

func (c *Client) DeleteWebhook(ctx context.Context, r DeleteWebhookRequest) (DeleteWebhookResponse, error) {
	if r.ID == "" {
		return DeleteWebhookResponse{}, fmt.Errorf("%w: %v", ErrInvalidRequest, "webhook ID must be a non-empty string")
	}

	req, err := http.NewRequestWithContext(ctx, "DELETE", fmt.Sprintf("%s/webhooks/%s", BASE_URL, r.ID), nil)
	if err != nil {
		return DeleteWebhookResponse{}, fmt.Errorf("%w: %v", ErrUnknown, err)
	}

	req.Header.Set("X-FIGMA-TOKEN", c.token)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return DeleteWebhookResponse{}, fmt.Errorf("%w: %v", ErrUnknown, err)
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		return DeleteWebhookResponse{}, nil
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return DeleteWebhookResponse{}, fmt.Errorf("%w: %v", ErrUnknown, err)
	}

	var errRes errResponse
	if err := json.Unmarshal(resBody, &errRes); err != nil {
		return DeleteWebhookResponse{}, fmt.Errorf("%w: %v", ErrUnknown, err)
	}

	switch res.StatusCode {
	case http.StatusUnauthorized, http.StatusForbidden:
		return DeleteWebhookResponse{}, fmt.Errorf("%w: %v", ErrInvalidCredentials, errRes.Message)
	case http.StatusBadRequest:
		return DeleteWebhookResponse{}, fmt.Errorf("%w: %v", ErrInvalidRequest, errRes.Message)
	case http.StatusNotFound:
		return DeleteWebhookResponse{}, ErrNotFound
	default:
		return DeleteWebhookResponse{}, ErrUnknown
	}
}
