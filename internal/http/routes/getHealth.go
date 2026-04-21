package routes

import (
	"net/http"

	"github.com/go-www/silverlining"
)

func GetHealth(ctx *silverlining.Context) {
	_ = ctx.WriteJSON(http.StatusOK, map[string]string{"status": "ok"})
}
