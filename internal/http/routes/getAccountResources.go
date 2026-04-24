package routes

import (
	"aliceSpeakerService/internal/model"
	"aliceSpeakerService/internal/store"
	"net/http"

	"github.com/go-www/silverlining"
)

func defaultAliceVoices() []model.VoiceOption {
	return []model.VoiceOption{
		{Value: "omazh", Label: "Omazh"},
		{Value: "dasha", Label: "Dasha"},
		{Value: "jane", Label: "Jane"},
		{Value: "alena", Label: "Alena"},
		{Value: "julia", Label: "Julia"},
		{Value: "zahar", Label: "Zahar"},
		{Value: "ermil", Label: "Ermil"},
		{Value: "filipp", Label: "Filipp"},
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
