package store

import (
	"aliceSpeakerService/internal/model"
	"encoding/json"
	"time"

	"go.etcd.io/bbolt"
)

type DeliveryRepository struct{}

func GetDeliveryRepository() *DeliveryRepository { return &DeliveryRepository{} }

func (r *DeliveryRepository) Save(delivery model.Delivery) error {
	if delivery.CreatedAt.IsZero() {
		delivery.CreatedAt = time.Now().UTC()
	}
	payload, err := json.Marshal(delivery)
	if err != nil {
		return err
	}
	return repository.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(DeliveriesBucket).Put([]byte(delivery.ID), payload)
	})
}
