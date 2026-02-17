package env

import (
	"log"
	"os"
	"strconv"
)

func String(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	return val
}

func StringOrDefault(key string, defaultValue string) string {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		return defaultValue
	}
	return val
}

func Int(key string) int {
	val := String(key)
	i, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("invalid int env var %s=%q", key, val)
	}
	return i
}
