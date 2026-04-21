package routes

import (
	"log"
	"net/http"

	"github.com/go-www/silverlining"
)

func NotFound(ctx *silverlining.Context) {
	if err := ctx.WriteJSON(http.StatusNotFound, map[string]int{"error": http.StatusNotFound}); err != nil {
		log.Print(err)
	}
}
