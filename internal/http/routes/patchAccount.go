package routes

import (
	"aliceSpeakerService/internal/model"
	"aliceSpeakerService/internal/store"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-www/silverlining"
)

type patchAccountBody struct {
	Title            *string `json:"title"`
	OAuthToken       *string `json:"oauth_token"`
	UnofficialXToken *string `json:"unofficial_x_token"`
	Transport        *string `json:"transport"`
	IsActive         *bool   `json:"is_active"`
}

func PatchAccount(ctx *silverlining.Context, accountID string, body []byte) {
	account, err := store.GetAccountRepository().FindByID(accountID)
	if err != nil {
		GetError(ctx, &Error{Message: err.Error(), Status: http.StatusNotFound})
		return
	}

	var payload patchAccountBody
	if err := json.Unmarshal(body, &payload); err != nil {
		GetError(ctx, &Error{Message: err.Error(), Status: http.StatusBadRequest})
		return
	}
	if payload.Title != nil {
		account.Title = strings.TrimSpace(*payload.Title)
	}
	if payload.OAuthToken != nil {
		account.OAuthToken = strings.TrimSpace(*payload.OAuthToken)
	}
	if payload.UnofficialXToken != nil {
		account.UnofficialXToken = strings.TrimSpace(*payload.UnofficialXToken)
	}
	if payload.Transport != nil {
		account.Transport = model.NormalizeTransport(strings.TrimSpace(*payload.Transport))
	}
	if payload.IsActive != nil {
		account.IsActive = *payload.IsActive
	}
	if err := store.GetAccountRepository().Save(account); err != nil {
		GetError(ctx, &Error{Message: err.Error(), Status: http.StatusInternalServerError})
		return
	}
	_ = ctx.WriteJSON(http.StatusOK, account)
}
