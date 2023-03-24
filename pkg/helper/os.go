package helper

import (
	"log"
	"os"
)

func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetEnvOtherwisePanic(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Panicf("missing environment variable %s", key)
	}
	return value
}
