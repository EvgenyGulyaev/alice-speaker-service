package transport

import "aliceSpeakerService/internal/model"

type Request struct {
	AccountID      string
	DeviceID       string
	ScenarioID     string
	InitiatorEmail string
	RecipientEmail string
	ConversationID string
	MessageID      string
	Text           string
	AudioURL       string
}

type Result struct {
	Status    string
	RequestID string
}

type SpeakerTransport interface {
	Name() string
	Announce(account model.Account, request Request) (Result, error)
}
