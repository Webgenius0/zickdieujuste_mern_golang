package dto

// VerifyReceiptRequest is the body for POST /api/v1/subscriptions/verify.
type VerifyReceiptRequest struct {
	Store          string `json:"store" example:"APPLE"           validate:"required,oneof=APPLE GOOGLE"`
	ReceiptPayload string `json:"receipt_payload" example:"MIIT7wYJKoZIhvcNAQcCoIIT4D..." validate:"required"` // Raw receipt string / token from the store
}

// WebhookRequest is the body for POST /api/v1/subscriptions/webhook.
// The exact shape depends on the store — kept generic here.
type WebhookRequest struct {
	Store   string                 `json:"store" example:"APPLE"`
	Payload map[string]interface{} `json:"payload"`
}
