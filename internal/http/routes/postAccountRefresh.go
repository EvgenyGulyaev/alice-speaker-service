package routes

import (
	"aliceSpeakerService/internal/store"
	"aliceSpeakerService/internal/yandex"
	"net/http"
	"time"

	"github.com/go-www/silverlining"
)

var yandexClient = yandex.NewClient()

func PostAccountRefresh(ctx *silverlining.Context, accountID string) {
	account, err := store.GetAccountRepository().FindByID(accountID)
	if err != nil {
		GetError(ctx, &Error{Message: err.Error(), Status: http.StatusNotFound})
		return
	}

	resources, err := yandexClient.LoadResources(account.OAuthToken)
	if err != nil {
		GetError(ctx, &Error{Message: err.Error(), Status: http.StatusBadGateway})
		return
	}
	if err := store.GetResourceRepository().ReplaceResources(accountID, resources); err != nil {
		GetError(ctx, &Error{Message: err.Error(), Status: http.StatusInternalServerError})
		return
	}

	account.LastSyncedAt = resourcesLatestUpdate(resources)
	if err := store.GetAccountRepository().Save(account); err != nil {
		GetError(ctx, &Error{Message: err.Error(), Status: http.StatusInternalServerError})
		return
	}

	_ = ctx.WriteJSON(http.StatusOK, resources)
}

func resourcesLatestUpdate(_ any) time.Time {
	return time.Now().UTC()
}
