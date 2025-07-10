package figma

type CreateWebhookRequest struct {
	EventType   string `json:"event_type"`
	ContextType string `json:"context"`
	ContextID   string `json:"context_id"`
	Endpoint    string `json:"endpoint"`
	Passcode    string `json:"passcode"`
	Description string `json:"description,omitempty"`
}

type DeleteWebhookRequest struct {
	ID string
}
