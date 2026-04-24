package routes

import (
	"aliceSpeakerService/internal/model"
	"aliceSpeakerService/internal/store"
	"net/http"

	"github.com/go-www/silverlining"
)

func defaultAliceVoices() []model.VoiceOption {
	return []model.VoiceOption{
		{Value: "jane", Label: "Jane"},
		{Value: "oksana", Label: "Oksana"},
		{Value: "zahar", Label: "Zahar"},
		{Value: "ermil", Label: "Ermil"},
	}
}

func GetAccountResources(ctx *silverlining.Context, accountID string) {
	resources, err := store.GetResourceRepository().GetResources(accountID)
	if err != nil {
		GetError(ctx, &Error{Message: err.Error(), Status: http.StatusInternalServerError})
		return
	}
	resources.Voices = defaultAliceVoices()
	_ = ctx.WriteJSON(http.StatusOK, resources)
}
