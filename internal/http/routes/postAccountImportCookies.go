package routes

import (
	"aliceSpeakerService/internal/model"
	"aliceSpeakerService/internal/store"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-www/silverlining"
)

type postAccountImportCookiesBody struct {
	Cookies string `json:"cookies"`
}

func PostAccountImportCookies(ctx *silverlining.Context, accountID string, body []byte) {
	account, err := store.GetAccountRepository().FindByID(accountID)
	if err != nil {
		GetError(ctx, &Error{Message: err.Error(), Status: http.StatusNotFound})
		return
	}

	var payload postAccountImportCookiesBody
	if err := json.Unmarshal(body, &payload); err != nil {
		GetError(ctx, &Error{Message: err.Error(), Status: http.StatusBadRequest})
		return
	}
	if strings.TrimSpace(payload.Cookies) == "" {
		GetError(ctx, &Error{Message: "cookies are required", Status: http.StatusBadRequest})
		return
	}

	xToken, err := yandexClient.ExtractXTokenFromCookies(payload.Cookies)
	if err != nil {
		GetError(ctx, &Error{Message: err.Error(), Status: http.StatusBadGateway})
		return
	}

	account.UnofficialXToken = strings.TrimSpace(xToken)
	account.Transport = model.TransportUnofficial
	if err := store.GetAccountRepository().Save(account); err != nil {
		GetError(ctx, &Error{Message: err.Error(), Status: http.StatusInternalServerError})
		return
	}

	_ = ctx.WriteJSON(http.StatusOK, map[string]any{
		"status":     "ok",
		"account_id": account.ID,
		"transport":  account.Transport,
	})
}
