package routes

import (
	"aliceSpeakerService/internal/model"
	"aliceSpeakerService/internal/store"
	"aliceSpeakerService/internal/transport"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-www/silverlining"
)

type announceScenarioBody struct {
	AccountID      string `json:"account_id"`
	DeviceID       string `json:"device_id"`
	ScenarioID     string `json:"scenario_id"`
	InitiatorEmail string `json:"initiator_email"`
	RecipientEmail string `json:"recipient_email"`
	ConversationID string `json:"conversation_id"`
	MessageID      string `json:"message_id"`
	Text           string `json:"text"`
	Voice          string `json:"voice"`
}

var announceTransport = transport.NewManager(
	transport.NewOfficialScenarioTransport(yandexClient),
	transport.NewUnofficialTransport(yandexClient),
)

func PostAnnounceScenario(ctx *silverlining.Context, body []byte) {
	var payload announceScenarioBody
	if err := json.Unmarshal(body, &payload); err != nil {
		GetError(ctx, &Error{Message: err.Error(), Status: http.StatusBadRequest})
		return
	}
	if strings.TrimSpace(payload.AccountID) == "" {
		GetError(ctx, &Error{Message: "account_id is required", Status: http.StatusBadRequest})
		return
	}

	account, err := store.GetAccountRepository().FindByID(payload.AccountID)
	if err != nil {
		GetError(ctx, &Error{Message: err.Error(), Status: http.StatusNotFound})
		return
	}

	result, err := announceTransport.Announce(account, transport.Request{
		AccountID:      payload.AccountID,
		DeviceID:       strings.TrimSpace(payload.DeviceID),
		ScenarioID:     payload.ScenarioID,
		InitiatorEmail: strings.TrimSpace(payload.InitiatorEmail),
		RecipientEmail: strings.TrimSpace(payload.RecipientEmail),
		ConversationID: strings.TrimSpace(payload.ConversationID),
		MessageID:      strings.TrimSpace(payload.MessageID),
		Text:           strings.TrimSpace(payload.Text),
		Voice:          strings.TrimSpace(payload.Voice),
	})
	delivery := model.Delivery{
		ID:             payload.AccountID + ":" + payload.ScenarioID + ":" + time.Now().UTC().Format(time.RFC3339Nano),
		AccountID:      payload.AccountID,
		ScenarioID:     payload.ScenarioID,
		InitiatorEmail: strings.TrimSpace(payload.InitiatorEmail),
		RecipientEmail: strings.TrimSpace(payload.RecipientEmail),
		ConversationID: strings.TrimSpace(payload.ConversationID),
		MessageID:      strings.TrimSpace(payload.MessageID),
		Status:         "sent",
		RequestID:      result.RequestID,
		CreatedAt:      time.Now().UTC(),
	}
	if err != nil {
		delivery.Status = "failed"
		delivery.ErrorMessage = err.Error()
		_ = store.GetDeliveryRepository().Save(delivery)
		GetError(ctx, &Error{Message: err.Error(), Status: http.StatusBadGateway})
		return
	}
	if err := store.GetDeliveryRepository().Save(delivery); err != nil {
		GetError(ctx, &Error{Message: err.Error(), Status: http.StatusInternalServerError})
		return
	}

	_ = ctx.WriteJSON(http.StatusOK, map[string]string{
		"status":      result.Status,
		"request_id":  result.RequestID,
		"delivery_id": delivery.ID,
	})
}
