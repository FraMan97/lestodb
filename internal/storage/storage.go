package storage

import (
	"sync"

	"github.com/FraMan97/lestodb/internal/database"
)

type Entry struct {
	Value     string
	Ttl       int
	ExpiredAt int64
}

type Storage struct {
	BackupDatabase database.EntryRecordRepository
	Data           map[int]map[string]*Entry
	Sharding       map[int]*sync.RWMutex
}

var DB *Storage

func NewStorage(bd database.EntryRecordRepository, shardsCount int) {
	DB = &Storage{
		BackupDatabase: bd,
		Data:           make(map[int]map[string]*Entry),
		Sharding:       make(map[int]*sync.RWMutex),
	}

	for i := 0; i < shardsCount; i++ {
		DB.Data[i] = make(map[string]*Entry)
		DB.Sharding[i] = &sync.RWMutex{}
	}
}
