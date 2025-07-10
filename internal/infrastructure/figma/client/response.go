package figma

type errResponse struct {
	Message string `json:"message"`
}

type CreateWebhookResponse struct {
	ID string `json:"id"`
}
