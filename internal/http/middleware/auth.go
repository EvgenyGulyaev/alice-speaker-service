package middleware

import (
	"aliceSpeakerService/internal/config"
	"strings"

	"github.com/go-www/silverlining"
)

func RequireServiceToken(next func(c *silverlining.Context)) func(c *silverlining.Context) {
	return func(c *silverlining.Context) {
		expected := strings.TrimSpace(config.LoadConfig().ServiceToken)
		if expected == "" {
			next(c)
			return
		}

		authorization, _ := c.RequestHeaders().Get("Authorization")
		token := strings.TrimSpace(strings.TrimPrefix(string(authorization), "Bearer "))
		if token != expected {
			_ = c.WriteJSON(401, map[string]string{"message": "unauthorized"})
			return
		}
		next(c)
	}
}
