package model

import "time"

type Household struct {
	ID        string    `json:"id"`
	AccountID string    `json:"account_id"`
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Room struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"account_id"`
	HouseholdID string    `json:"household_id"`
	Name        string    `json:"name"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Device struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"account_id"`
	RoomID      string    `json:"room_id"`
	HouseholdID string    `json:"household_id"`
	Name        string    `json:"name"`
	Type        string    `json:"type"`
	Online      bool      `json:"online"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Scenario struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"account_id"`
	HouseholdID string    `json:"household_id"`
	Name        string    `json:"name"`
	IsActive    bool      `json:"is_active"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type VoiceOption struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type Resources struct {
	Households []Household   `json:"households"`
	Rooms      []Room        `json:"rooms"`
	Devices    []Device      `json:"devices"`
	Scenarios  []Scenario    `json:"scenarios"`
	Voices     []VoiceOption `json:"voices"`
}
