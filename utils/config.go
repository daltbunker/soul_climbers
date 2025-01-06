package utils

import (
	"log"
	"os"
)

func ValidateEnv() {
	requiredEnvVariables := []string{"DB_URL", "PORT", "SESSION_KEY", "MODE"}
	for i := 0; i < len(requiredEnvVariables); i++ {
		_, valExists := os.LookupEnv(requiredEnvVariables[i])
		if !valExists {
			log.Fatal("Missing env variable: " + requiredEnvVariables[i])
		}
	}
}