package store

import (
	"aliceSpeakerService/internal/model"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"go.etcd.io/bbolt"
)

type ResourceRepository struct{}

func GetResourceRepository() *ResourceRepository { return &ResourceRepository{} }

func resourceKey(accountID, entityID string) []byte {
	return []byte(fmt.Sprintf("%s:%s", accountID, entityID))
}

func (r *ResourceRepository) ReplaceResources(accountID string, resources model.Resources) error {
	now := time.Now().UTC()
	for i := range resources.Rooms {
		resources.Rooms[i].AccountID = accountID
		resources.Rooms[i].UpdatedAt = now
	}
	for i := range resources.Devices {
		resources.Devices[i].AccountID = accountID
		resources.Devices[i].UpdatedAt = now
	}
	for i := range resources.Scenarios {
		resources.Scenarios[i].AccountID = accountID
		resources.Scenarios[i].UpdatedAt = now
	}

	return repository.Update(func(tx *bbolt.Tx) error {
		if err := clearAccountBucketItems(tx.Bucket(RoomsBucket), accountID); err != nil {
			return err
		}
		if err := clearAccountBucketItems(tx.Bucket(DevicesBucket), accountID); err != nil {
			return err
		}
		if err := clearAccountBucketItems(tx.Bucket(ScenariosBucket), accountID); err != nil {
			return err
		}

		for _, room := range resources.Rooms {
			payload, err := json.Marshal(room)
			if err != nil {
				return err
			}
			if err := tx.Bucket(RoomsBucket).Put(resourceKey(accountID, room.ID), payload); err != nil {
				return err
			}
		}
		for _, device := range resources.Devices {
			payload, err := json.Marshal(device)
			if err != nil {
				return err
			}
			if err := tx.Bucket(DevicesBucket).Put(resourceKey(accountID, device.ID), payload); err != nil {
				return err
			}
		}
		for _, scenario := range resources.Scenarios {
			payload, err := json.Marshal(scenario)
			if err != nil {
				return err
			}
			if err := tx.Bucket(ScenariosBucket).Put(resourceKey(accountID, scenario.ID), payload); err != nil {
				return err
			}
		}
		return nil
	})
}

func clearAccountBucketItems(bucket *bbolt.Bucket, accountID string) error {
	keys := make([][]byte, 0)
	err := bucket.ForEach(func(key, _ []byte) error {
		if strings.HasPrefix(string(key), accountID+":") {
			copyKey := append([]byte(nil), key...)
			keys = append(keys, copyKey)
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, key := range keys {
		if err := bucket.Delete(key); err != nil {
			return err
		}
	}
	return nil
}

func (r *ResourceRepository) GetResources(accountID string) (model.Resources, error) {
	result := model.Resources{
		Rooms:     make([]model.Room, 0),
		Devices:   make([]model.Device, 0),
		Scenarios: make([]model.Scenario, 0),
	}

	err := repository.View(func(tx *bbolt.Tx) error {
		if err := tx.Bucket(RoomsBucket).ForEach(func(_, value []byte) error {
			var room model.Room
			if err := json.Unmarshal(value, &room); err != nil {
				return err
			}
			if room.AccountID == accountID {
				result.Rooms = append(result.Rooms, room)
			}
			return nil
		}); err != nil {
			return err
		}
		if err := tx.Bucket(DevicesBucket).ForEach(func(_, value []byte) error {
			var device model.Device
			if err := json.Unmarshal(value, &device); err != nil {
				return err
			}
			if device.AccountID == accountID {
				result.Devices = append(result.Devices, device)
			}
			return nil
		}); err != nil {
			return err
		}
		return tx.Bucket(ScenariosBucket).ForEach(func(_, value []byte) error {
			var scenario model.Scenario
			if err := json.Unmarshal(value, &scenario); err != nil {
				return err
			}
			if scenario.AccountID == accountID {
				result.Scenarios = append(result.Scenarios, scenario)
			}
			return nil
		})
	})

	sort.Slice(result.Rooms, func(i, j int) bool { return result.Rooms[i].Name < result.Rooms[j].Name })
	sort.Slice(result.Devices, func(i, j int) bool { return result.Devices[i].Name < result.Devices[j].Name })
	sort.Slice(result.Scenarios, func(i, j int) bool { return result.Scenarios[i].Name < result.Scenarios[j].Name })
	return result, err
}
