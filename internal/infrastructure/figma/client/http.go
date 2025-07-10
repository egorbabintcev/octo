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
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.figma.com/v2/webhooks", reqBody)

	req.Header.Set("X-FIGMA-TOKEN", c.token)
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return CreateWebhookResponse{}, fmt.Errorf("%w: %v", ErrUnknown, err)
	}

	res, err := c.httpClient.Do(req)

	if err != nil {
		return CreateWebhookResponse{}, fmt.Errorf("%w: %v", ErrUnknown, err)
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		return CreateWebhookResponse{}, fmt.Errorf("%w: %v", ErrUnknown, err)
	}

	if res.StatusCode != http.StatusOK {
		errRes := errResponse{}

		if err := json.Unmarshal(resBody, &errRes); err != nil {
			return CreateWebhookResponse{}, fmt.Errorf("%w: %v", ErrUnknown, err)
		}

		if res.StatusCode == http.StatusUnauthorized || res.StatusCode == http.StatusForbidden {
			return CreateWebhookResponse{}, fmt.Errorf("%w: %v", ErrInvalidCredentials, errRes.Message)
		}

		if res.StatusCode == http.StatusBadRequest {
			return CreateWebhookResponse{}, fmt.Errorf("%w: %v", ErrInvalidRequest, errRes.Message)
		}

		return CreateWebhookResponse{}, ErrUnknown
	}

	wh := CreateWebhookResponse{}

	if err := json.Unmarshal(resBody, &wh); err != nil {
		return CreateWebhookResponse{}, nil
	}

	return wh, nil
}
