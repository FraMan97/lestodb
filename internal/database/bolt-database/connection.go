package boltdatabase

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"go.etcd.io/bbolt"
)

var DB *bbolt.DB

const DefaultBucket = "LestoData"

func Init(filePath string) error {
	var err error

	opts := &bbolt.Options{Timeout: 1 * time.Second}
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creating directory %s: %w", dir, err)
	}
	DB, err = bbolt.Open(filePath, 0600, opts)
	if err != nil {
		return fmt.Errorf("error opening connection to boltdb: %w", err)
	}

	err = DB.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(DefaultBucket))
		return err
	})

	if err != nil {
		return fmt.Errorf("creation bucket error: %w", err)
	}

	log.Println("BoltDB connected on:", filePath)
	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
