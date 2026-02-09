package env

import (
	"os"
	"strconv"
)

var SHARDING_COUNT = getEnvInt("SHARDING_COUNT", 36)
var URL = getEnv("URL", "0.0.0.0")
var PORT = getEnvInt("PORT", 2001)
var FILE_PATH = getEnv("FILE_PATH", "/app/database/lesto.db")

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		i, _ := strconv.Atoi(value)
		return i
	}
	return fallback
}
