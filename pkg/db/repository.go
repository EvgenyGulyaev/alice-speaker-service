package db

import "go.etcd.io/bbolt"

type Repository struct {
	db *Db
}

func GetRepository(filename string) *Repository {
	return &Repository{db: Init(filename)}
}

func (r *Repository) EnsureBuckets(buckets [][]byte) error {
	for _, bucket := range buckets {
		if err := r.db.EnsureBucket(bucket); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) Update(fn func(*bbolt.Tx) error) error {
	return r.db.DB.Update(fn)
}

func (r *Repository) View(fn func(*bbolt.Tx) error) error {
	return r.db.DB.View(fn)
}
