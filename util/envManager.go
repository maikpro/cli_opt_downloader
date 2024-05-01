package util

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvString(str string) (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file:", err)
		return "", err
	}

	envString := os.Getenv(str)
	if envString == "" {
		log.Printf("%s environment variable is not set or cannot be found", str)
		return "", errors.New("environment variable is not set or cannot be found")
	}
	return envString, nil
}
