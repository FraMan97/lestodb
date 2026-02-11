package main

import (
	"bufio"
	"fmt"
	"log"
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

	log.Println("Loading data from BoltDB...")
	if err := storage.DB.LoadFromDatabase(); err != nil {
		log.Printf("Warning: failed to load initial data: %v\n", err)
	} else {
		log.Println("Initial data loaded successfully")
	}

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
	currentShard := 0

	for {
		now := time.Now().Unix()

		shardIndex := currentShard % env.SHARDING_COUNT

		storage.DB.Sharding[shardIndex].Lock()
		for key, val := range storage.DB.Data[shardIndex] {
			if val.ExpiredAt != 0 && now >= val.ExpiredAt {
				delete(storage.DB.Data[shardIndex], key)
			}
		}
		storage.DB.Sharding[shardIndex].Unlock()

		currentShard++

		if currentShard >= env.SHARDING_COUNT {
			currentShard = 0
		}

		time.Sleep(time.Millisecond * 100)
	}
}
