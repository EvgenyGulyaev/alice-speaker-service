package routes

import (
	"aliceSpeakerService/internal/store"
	"net/http"

	"github.com/go-www/silverlining"
)

func GetAccountResources(ctx *silverlining.Context, accountID string) {
	resources, err := store.GetResourceRepository().GetResources(accountID)
	if err != nil {
		GetError(ctx, &Error{Message: err.Error(), Status: http.StatusInternalServerError})
		return
	}
	_ = ctx.WriteJSON(http.StatusOK, resources)
}
