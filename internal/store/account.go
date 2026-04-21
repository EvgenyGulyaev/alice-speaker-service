package store

import (
	"aliceSpeakerService/internal/model"
	"encoding/json"
	"errors"
	"sort"
	"time"

	"go.etcd.io/bbolt"
)

type AccountRepository struct{}

func GetAccountRepository() *AccountRepository { return &AccountRepository{} }

func (r *AccountRepository) Save(account model.Account) error {
	now := time.Now().UTC()
	if account.CreatedAt.IsZero() {
		account.CreatedAt = now
	}
	account.UpdatedAt = now
	if account.Provider == "" {
		account.Provider = "yandex"
	}

	payload, err := json.Marshal(account)
	if err != nil {
		return err
	}

	return repository.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(AccountsBucket).Put([]byte(account.ID), payload)
	})
}

func (r *AccountRepository) FindByID(id string) (model.Account, error) {
	var account model.Account
	err := repository.View(func(tx *bbolt.Tx) error {
		raw := tx.Bucket(AccountsBucket).Get([]byte(id))
		if len(raw) == 0 {
			return errors.New("account not found")
		}
		return json.Unmarshal(raw, &account)
	})
	return account, err
}

func (r *AccountRepository) List() ([]model.Account, error) {
	result := make([]model.Account, 0)
	err := repository.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(AccountsBucket).ForEach(func(_, value []byte) error {
			var account model.Account
			if err := json.Unmarshal(value, &account); err != nil {
				return err
			}
			result = append(result, account)
			return nil
		})
	})
	sort.Slice(result, func(i, j int) bool {
		return result[i].Title < result[j].Title
	})
	return result, err
}
