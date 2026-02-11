package storage

import (
	"hash/fnv"
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

func (s *Storage) LoadFromDatabase() error {
	records, err := s.BackupDatabase.GetAll()
	if err != nil {
		return err
	}

	for _, record := range records {
		h := fnv.New32a()
		h.Write([]byte(record.Key))
		index := int(h.Sum32()) % len(s.Data)

		s.Sharding[index].Lock()
		s.Data[index][record.Key] = &Entry{
			Value:     record.Value,
			Ttl:       0,
			ExpiredAt: 0,
		}
		s.Sharding[index].Unlock()
	}
	return nil
}
