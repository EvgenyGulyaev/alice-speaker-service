package routes

import (
	"log"

	"github.com/go-www/silverlining"
)

type Error struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func GetError(ctx *silverlining.Context, value *Error) {
	if err := ctx.WriteJSON(value.Status, value); err != nil {
		log.Print(err)
	}
}
