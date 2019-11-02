package pkg

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var client *mongo.Client

func SetUp(ctx context.Context) {
	log.Println("Database Address: " + os.Getenv(DatabaseAddress))
	clientOpt := options.Client().ApplyURI(os.Getenv(DatabaseAddress))
	client, _ = mongo.Connect(ctx, clientOpt)
	err := client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal("Error to connect with db: ", err)
	}

	log.Println("Connected with the storage")
}

func checkCollections() {
	// CHECK IF EXISTS PROBLEM COLLECTION AND SUBMISSION COLLECTION
}


func SaveProblem(p Problem) (interface{}, error){
	collection := client.Database(os.Getenv(DatabaseName)).Collection(os.Getenv(ProblemCollection))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return collection.InsertOne(ctx, &p)
}

func RetrieveProblem(problemId string) (*Problem, error) {
	log.Println("Retriving problem [" + problemId + "]")
	pId, err := primitive.ObjectIDFromHex(problemId)

	if err != nil {
		log.Printf("Problem Id [%s] with wrong shape: %s", problemId, err.Error())
		return nil, err
	}

	filter := bson.M{"_id":pId}

	var p Problem

	collection := client.Database(os.Getenv(DatabaseName)).Collection(os.Getenv(ProblemCollection))

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	e := collection.FindOne(ctx, filter).Decode(&p)

	return &p, e

}

func SaveSubmission(submission Submission) (*mongo.InsertOneResult, error) {
	collection := client.Database(os.Getenv(DatabaseName)).Collection(os.Getenv(SubmissionCollection))
	if collection == nil {
		log.Println("Collection is nil!")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return collection.InsertOne(ctx, &submission)
}

