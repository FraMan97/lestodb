package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"strings"
	"time"

	boltdatabase "github.com/FraMan97/lestodb/internal/database/bolt-database"
	"github.com/FraMan97/lestodb/internal/env"
	"github.com/FraMan97/lestodb/internal/storage"
	"github.com/FraMan97/lestodb/internal/validator"
)

func main() {
	err := boltdatabase.Init(env.FILE_PATH)
	if err != nil {
		log.Fatal("error opening DB:", err)
	}
	defer boltdatabase.Close()

	repo := boltdatabase.NewBoltEntryRepository(boltdatabase.DB)

	storage.NewStorage(repo, env.SHARDING_COUNT)

	log.Println("Server is listening on " + fmt.Sprintf("%s:%d", env.URL, env.PORT))
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", env.PORT))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	go cleanTTLValues()
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		go handleCommand(conn)
	}
}

func handleCommand(connection net.Conn) {
	defer connection.Close()
	reader := bufio.NewReader(connection)
	for {
		commandsString, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		commandsString = strings.TrimSpace(commandsString)

		if commandsString == "" {
			continue
		}

		commands := strings.Split(commandsString, ";")

		var results []string
		for _, cmd := range commands {
			cmd = strings.TrimSpace(cmd)
			command, err := validator.ValidateCommand(cmd)
			if err != nil {
				results = append(results, err.Error())
				continue
			}
			result, err := command.ExecCommand()
			if err != nil {
				results = append(results, err.Error())
				continue
			}
			results = append(results, result)
		}
		connection.Write([]byte("Results: [" + strings.Join(results, "; ") + "]\n"))
	}
}

func cleanTTLValues() {
	for {
		shard := rand.IntN(env.SHARDING_COUNT)
		now := time.Now().Unix()
		storage.DB.Sharding[shard].Lock()
		for key, val := range storage.DB.Data[shard] {
			if now >= val.ExpiredAt {
				delete(storage.DB.Data[shard], key)
			}
		}
		storage.DB.Sharding[shard].Unlock()
		time.Sleep(time.Second * 1)
	}
}
