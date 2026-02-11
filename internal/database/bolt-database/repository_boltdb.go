package boltdatabase

import (
	"fmt"

	"github.com/FraMan97/lestodb/internal/database"
	"go.etcd.io/bbolt"
)

type BoltEntryRepository struct {
	db *bbolt.DB
}

func NewBoltEntryRepository(db *bbolt.DB) *BoltEntryRepository {
	return &BoltEntryRepository{db: db}
}

func (r *BoltEntryRepository) Save(entry *database.EntryRecord) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(DefaultBucket))
		return b.Put([]byte(entry.Key), []byte(entry.Value))
	})
}

func (r *BoltEntryRepository) Delete(key string) error {
	return r.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(DefaultBucket))
		return b.Delete([]byte(key))
	})
}

func (r *BoltEntryRepository) Get(key string) (*database.EntryRecord, error) {
	var entry *database.EntryRecord

	err := r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(DefaultBucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", DefaultBucket)
		}

		v := b.Get([]byte(key))
		if v == nil {
			return nil
		}
		entry = &database.EntryRecord{
			Key:   key,
			Value: string(v),
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return entry, nil
}

func (r *BoltEntryRepository) GetAll() ([]database.EntryRecord, error) {
	var entries []database.EntryRecord

	err := r.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(DefaultBucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", DefaultBucket)
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			entries = append(entries, database.EntryRecord{
				Key:   string(k),
				Value: string(v),
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return entries, nil
}
