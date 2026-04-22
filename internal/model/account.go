package model

import "time"

const (
	TransportOfficial   = "official"
	TransportUnofficial = "unofficial"
)

type Account struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Provider     string    `json:"provider"`
	Transport    string    `json:"transport"`
	OAuthToken   string    `json:"-"`
	IsActive     bool      `json:"is_active"`
	LastSyncedAt time.Time `json:"last_synced_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NormalizeTransport(value string) string {
	switch value {
	case TransportUnofficial:
		return TransportUnofficial
	default:
		return TransportOfficial
	}
}
