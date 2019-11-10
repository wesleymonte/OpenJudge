package pkg

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

const DatabaseAddress = "DATABASE_ADDRESS"
const DatabaseName = "DATABASE_NAME"
const ProblemCollection = "PROBLEM_COLLECTION"
const SubmissionCollection = "SUBMISSION_COLLECTION"
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
	err := CreateFolder(SubmissionsDirName)
	if err != nil {
		log.Fatal("Error while create submission folder: " + err.Error())
	}
}

func CreateFolder(dirName string) error {
	_, err := os.Stat(dirName)
	if os.IsExist(err) {
		log.Println("Directory [" + dirName + "] already exists ")
		return nil
	}
	if os.IsNotExist(err) {
		log.Println("Creating directory [" + dirName + "]")
		if err := os.MkdirAll(dirName, os.ModePerm); err != nil {
			log.Println("Error while create directory [" + dirName + "]")
			return err
		}
		return nil
	}
	return err
}