package routes

import (
	"aliceSpeakerService/internal/store"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-www/silverlining"
)

type postAccountCleanupScenariosBody struct {
	DeviceID string `json:"device_id"`
}

func PostAccountCleanupScenarios(ctx *silverlining.Context, accountID string, body []byte) {
	account, err := store.GetAccountRepository().FindByID(accountID)
	if err != nil {
		GetError(ctx, &Error{Message: err.Error(), Status: http.StatusNotFound})
		return
	}

	var payload postAccountCleanupScenariosBody
	if len(body) > 0 {
		if err := json.Unmarshal(body, &payload); err != nil {
			GetError(ctx, &Error{Message: err.Error(), Status: http.StatusBadRequest})
			return
		}
	}

	deleted, err := yandexClient.CleanupCloudTTSScenarios(account.UnofficialXToken, strings.TrimSpace(payload.DeviceID))
	if err != nil {
		GetError(ctx, &Error{Message: err.Error(), Status: http.StatusBadGateway})
		return
	}

	_ = ctx.WriteJSON(http.StatusOK, map[string]any{
		"status":     "ok",
		"account_id": account.ID,
		"device_id":  strings.TrimSpace(payload.DeviceID),
		"deleted":    deleted,
	})
}
