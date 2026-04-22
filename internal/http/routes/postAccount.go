package routes

import (
	"aliceSpeakerService/internal/model"
	"aliceSpeakerService/internal/store"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-www/silverlining"
)

type postAccountBody struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	OAuthToken string `json:"oauth_token"`
	Transport  string `json:"transport"`
	IsActive   bool   `json:"is_active"`
}

func PostAccount(ctx *silverlining.Context, body []byte) {
	var payload postAccountBody
	if err := json.Unmarshal(body, &payload); err != nil {
		GetError(ctx, &Error{Message: err.Error(), Status: http.StatusBadRequest})
		return
	}
	if strings.TrimSpace(payload.ID) == "" || strings.TrimSpace(payload.Title) == "" || strings.TrimSpace(payload.OAuthToken) == "" {
		GetError(ctx, &Error{Message: "id, title and oauth_token are required", Status: http.StatusBadRequest})
		return
	}

	account := model.Account{
		ID:         strings.TrimSpace(payload.ID),
		Title:      strings.TrimSpace(payload.Title),
		Provider:   "yandex",
		Transport:  model.NormalizeTransport(strings.TrimSpace(payload.Transport)),
		OAuthToken: strings.TrimSpace(payload.OAuthToken),
		IsActive:   payload.IsActive,
		CreatedAt:  time.Now().UTC(),
	}
	if err := store.GetAccountRepository().Save(account); err != nil {
		GetError(ctx, &Error{Message: err.Error(), Status: http.StatusInternalServerError})
		return
	}
	_ = ctx.WriteJSON(http.StatusCreated, account)
}
