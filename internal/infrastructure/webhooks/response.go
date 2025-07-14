package webhooks

type Webhook struct {
	EventType   string
	ContextType string
	ContextID   string
	Endpoint    string
	Description string
}

type GetAllResponse struct {
	Webhooks []Webhook
}

type GetByIdResponse struct {
	Webhook Webhook
}

type GetByParamsResponse struct {
	Webhook Webhook
}

type CreateResponse struct {
	ID int
}

type DeleteResponse struct{}
