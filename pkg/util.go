package pkg

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

const DatabaseAddress = "DATABASE_ADDRESS"
const DatabaseName = "DATABASE_NAME"
const ProblemCollection = "PROBLEM_COLLECTION"

func ValidateEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file specified")
	}

	if _, exists := os.LookupEnv(DatabaseAddress); !exists {
		log.Fatal("No storage address on the env")
	} else if _, exists := os.LookupEnv(DatabaseName); !exists {
		log.Fatal("No storage name on the env")
	} else {
		log.Println("Environment loaded with success")
	}
}