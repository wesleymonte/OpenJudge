package pkg

import (
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"log"
	"os"
)

const DatabaseAddress = "DATABASE_ADDRESS"
const DatabaseName = "DATABASE_NAME"
const ProblemCollection = "PROBLEM_COLLECTION"
const SubmissionsDirName = "submissions"

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

func CreateSubmissionsFolder() {
	_, err := os.Stat(SubmissionsDirName)
	if os.IsNotExist(err) {
		log.Println("Not found submissions folder")
		if err := os.Mkdir(SubmissionsDirName, os.ModePerm); err != nil {
			log.Fatal("Error while create submissions dir: " + err.Error())
		}
		return
	}
	if err != nil {
		log.Fatal("Error while create submission folder: " + err.Error())
	}
}

func GetRandomUUID() string {
	var id uuid.UUID = uuid.New()
	return id.String()
}