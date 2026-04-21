package model

import "time"

type Delivery struct {
	ID             string    `json:"id"`
	AccountID      string    `json:"account_id"`
	ScenarioID     string    `json:"scenario_id"`
	InitiatorEmail string    `json:"initiator_email"`
	RecipientEmail string    `json:"recipient_email"`
	ConversationID string    `json:"conversation_id"`
	MessageID      string    `json:"message_id"`
	Status         string    `json:"status"`
	RequestID      string    `json:"request_id"`
	ErrorMessage   string    `json:"error_message"`
	CreatedAt      time.Time `json:"created_at"`
}
