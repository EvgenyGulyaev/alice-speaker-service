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

type accountRecord struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Provider     string    `json:"provider"`
	Transport    string    `json:"transport"`
	OAuthToken   string    `json:"oauth_token"`
	IsActive     bool      `json:"is_active"`
	LastSyncedAt time.Time `json:"last_synced_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func accountRecordFromModel(account model.Account) accountRecord {
	return accountRecord{
		ID:           account.ID,
		Title:        account.Title,
		Provider:     account.Provider,
		Transport:    model.NormalizeTransport(account.Transport),
		OAuthToken:   account.OAuthToken,
		IsActive:     account.IsActive,
		LastSyncedAt: account.LastSyncedAt,
		CreatedAt:    account.CreatedAt,
		UpdatedAt:    account.UpdatedAt,
	}
}

func (record accountRecord) toModel() model.Account {
	return model.Account{
		ID:           record.ID,
		Title:        record.Title,
		Provider:     record.Provider,
		Transport:    model.NormalizeTransport(record.Transport),
		OAuthToken:   record.OAuthToken,
		IsActive:     record.IsActive,
		LastSyncedAt: record.LastSyncedAt,
		CreatedAt:    record.CreatedAt,
		UpdatedAt:    record.UpdatedAt,
	}
}

func (r *AccountRepository) Save(account model.Account) error {
	now := time.Now().UTC()
	if account.CreatedAt.IsZero() {
		account.CreatedAt = now
	}
	account.UpdatedAt = now
	if account.Provider == "" {
		account.Provider = "yandex"
	}
	account.Transport = model.NormalizeTransport(account.Transport)

	payload, err := json.Marshal(accountRecordFromModel(account))
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
		var record accountRecord
		if err := json.Unmarshal(raw, &record); err != nil {
			return err
		}
		account = record.toModel()
		return nil
	})
	return account, err
}

func (r *AccountRepository) List() ([]model.Account, error) {
	result := make([]model.Account, 0)
	err := repository.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(AccountsBucket).ForEach(func(_, value []byte) error {
			var record accountRecord
			if err := json.Unmarshal(value, &record); err != nil {
				return err
			}
			result = append(result, record.toModel())
			return nil
		})
	})
	sort.Slice(result, func(i, j int) bool {
		return result[i].Title < result[j].Title
	})
	return result, err
}
