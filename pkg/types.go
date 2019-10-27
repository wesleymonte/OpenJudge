package pkg

//ADD VALIDATE
type Problem struct {
	TimeLimit   int64      `json:"time_limit" bson:"time_limit"`
	MemoryLimit int64      `json:"memory_limit" bson:"memory_limit"`
	TestCases   []TestCase `json:"test_cases" bson:"test_cases"`
}

type TestCase struct {
	In  string `json:"in" bson:"in"`
	Out string `json:"out" bson:"out"`
}

type Version struct {
	Tag string `json:"tag" bson:"tag"`
}
