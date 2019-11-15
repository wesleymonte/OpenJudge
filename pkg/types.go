package pkg

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//ADD VALIDATE
type Problem struct {
	ID 			primitive.ObjectID 	`json:"id" bson:"_id,omitempty"`
	TimeLimit   int64      `json:"time_limit" bson:"time_limit"`
	MemoryLimit int64      `json:"memory_limit" bson:"memory_limit"`
	TestCases   []TestCase `json:"test_cases" bson:"test_cases"`
}

type TestCase struct {
	In  string `json:"input" bson:"input"`
	Out string `json:"output" bson:"output"`
}

type Submission struct {
	ID 			primitive.ObjectID 	`json:"id" bson:"_id"`
	ProblemId 	string				`json:"problem_id" bson:"problem_id"`
	State 		string				`json:"state" bson:"state"`
	Result 		string				`json:"result" bson:"result"`
}

type Version struct {
	Tag string `json:"tag" bson:"tag"`
}
