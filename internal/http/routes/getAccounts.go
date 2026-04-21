package routes

import (
	"aliceSpeakerService/internal/store"
	"net/http"

	"github.com/go-www/silverlining"
)

func GetAccounts(ctx *silverlining.Context) {
	accounts, err := store.GetAccountRepository().List()
	if err != nil {
		GetError(ctx, &Error{Message: err.Error(), Status: http.StatusInternalServerError})
		return
	}
	_ = ctx.WriteJSON(http.StatusOK, map[string]any{"items": accounts})
}
