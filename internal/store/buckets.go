package store

import (
	"aliceSpeakerService/pkg/db"
)

var (
	AccountsBucket   = []byte("Accounts")
	RoomsBucket      = []byte("Rooms")
	DevicesBucket    = []byte("Devices")
	ScenariosBucket  = []byte("Scenarios")
	DeliveriesBucket = []byte("Deliveries")
)

var repository *db.Repository

func InitStore(repo *db.Repository) error {
	repository = repo
	return repo.EnsureBuckets([][]byte{
		AccountsBucket,
		RoomsBucket,
		DevicesBucket,
		ScenariosBucket,
		DeliveriesBucket,
	})
}

func GetRepository() *db.Repository {
	return repository
}
