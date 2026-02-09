package commands

import (
	"fmt"
	"hash/fnv"
	"log"
	"time"

	"github.com/FraMan97/lestodb/internal/database"
	"github.com/FraMan97/lestodb/internal/env"
	"github.com/FraMan97/lestodb/internal/storage"
)

type CommandInterface interface {
	ExecCommand() (string, error)
}

type SetCommand struct {
	Key   string
	Value string
	Ttl   int
}

type GetCommand struct {
	Key string
}

type DelCommand struct {
	Key string
}

type BackupCommand struct {
	Key string
}

type RestoreCommand struct {
	Key string
	Ttl int
}

func (s *SetCommand) ExecCommand() (string, error) {
	log.Printf("SET: %s = %s (%d)\n", s.Key, s.Value, s.Ttl)
	index := getShardIndex(s.Key, env.SHARDING_COUNT)
	storage.DB.Sharding[index].Lock()
	defer storage.DB.Sharding[index].Unlock()
	storage.DB.Data[index][s.Key] = &storage.Entry{Value: s.Value, Ttl: s.Ttl, ExpiredAt: time.Now().Unix() + int64(s.Ttl)}
	return "OK", nil
}

func (g *GetCommand) ExecCommand() (string, error) {
	log.Printf("GET: %s\n", g.Key)
	index := getShardIndex(g.Key, env.SHARDING_COUNT)
	storage.DB.Sharding[index].RLock()
	defer storage.DB.Sharding[index].RUnlock()
	result, ok := storage.DB.Data[index][g.Key]
	if !ok || result == nil {
		return "KO", fmt.Errorf("key %s not found", g.Key)
	}
	return result.Value, nil
}

func (d *DelCommand) ExecCommand() (string, error) {
	log.Printf("DEL: %s\n", d.Key)
	index := getShardIndex(d.Key, env.SHARDING_COUNT)
	storage.DB.Sharding[index].Lock()
	defer storage.DB.Sharding[index].Unlock()
	delete(storage.DB.Data[index], d.Key)
	return "OK", nil
}

func (b *BackupCommand) ExecCommand() (string, error) {
	log.Printf("BACKUP: %s\n", b.Key)
	if b.Key == "ALL" {
		for index, shard := range storage.DB.Sharding {
			shard.Lock()
			for key, value := range storage.DB.Data[index] {
				if storage.DB.Data[index][key] == nil {
					continue
				}
				err := storage.DB.BackupDatabase.Save(&database.EntryRecord{Key: key, Value: value.Value})
				if err != nil {
					continue
				}
			}
			shard.Unlock()
		}
		return "OK", nil
	} else {
		index := getShardIndex(b.Key, env.SHARDING_COUNT)
		storage.DB.Sharding[index].Lock()
		if storage.DB.Data[index][b.Key] == nil {
			return "KO", fmt.Errorf("key %s not found", b.Key)
		}
		err := storage.DB.BackupDatabase.Save(&database.EntryRecord{Key: b.Key, Value: storage.DB.Data[index][b.Key].Value})
		if err != nil {
			storage.DB.Sharding[index].Unlock()
			return "KO", fmt.Errorf("error saving key %s", b.Key)
		}
		storage.DB.Sharding[index].Unlock()
		return "OK", nil
	}
}

func (r *RestoreCommand) ExecCommand() (string, error) {
	log.Printf("RESTORE: %s\n", r.Key)
	if r.Key == "ALL" {
		for index, shard := range storage.DB.Sharding {
			shard.Lock()
			for key, _ := range storage.DB.Data[index] {
				data, err := storage.DB.BackupDatabase.Get(key)
				if err != nil {
					continue
				}
				if storage.DB.Data[index][r.Key] == nil {
					storage.DB.Data[index][key] = &storage.Entry{Value: data.Value}
				}
				storage.DB.Data[index][key].Value = data.Value
				storage.DB.Data[index][key].ExpiredAt = time.Now().Unix() + int64(r.Ttl)
				storage.DB.Data[index][key].Ttl = r.Ttl

			}
			shard.Unlock()
		}
		return "OK", nil
	} else {
		index := getShardIndex(r.Key, env.SHARDING_COUNT)
		storage.DB.Sharding[index].Lock()
		data, err := storage.DB.BackupDatabase.Get(r.Key)
		if err != nil {
			storage.DB.Sharding[index].Unlock()
			return "KO", fmt.Errorf("error get key %s", r.Key)
		}
		if storage.DB.Data[index][r.Key] == nil {
			storage.DB.Data[index][r.Key] = &storage.Entry{Value: data.Value}
		}
		storage.DB.Data[index][r.Key].Value = data.Value
		storage.DB.Sharding[index].Unlock()
		return "OK", nil
	}
}

func getShardIndex(key string, numShards int) int {
	h := fnv.New32a()
	h.Write([]byte(key))
	hashTotale := h.Sum32()
	return int(hashTotale) % numShards
}
