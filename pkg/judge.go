package pkg

import (
	"fmt"
	"log"
	"os"
	dclient "github.com/docker/docker/client"
)

type Judge struct {
	Name string
	in chan Submission
}

func New(name string, in chan Submission) *Judge {
	j := new(Judge)
	j.Name = name
	j.in = in
	return j
}

func (j Judge) Submit(s *Submission) {
	//problem, err := RetrieveProblem(s.ProblemId)
	//if err != nil {
	//	log.Println("Error while retrieve problem [" + s.ProblemId + "]: " + err.Error())
	//}
}

func (j *Judge) LoadProblem(problem *Problem) {
	dir := fmt.Sprintf("./problems/%s", problem.ID.Hex())
	inputDir := dir + "/in"
	outputDir := dir + "/out"
	err := CreateFolder(inputDir)
	err = CreateFolder(outputDir)
	if err != nil {

	}
	for i, t := range problem.TestCases {
		var inputFileName string = fmt.Sprintf("%d.in", i + 1)
		var inputFilePath string = inputDir + "/" + inputFileName

		var outputFileName string = fmt.Sprintf("%d.out", i + 1)
		var outputFilePath string = outputDir + "/" + outputFileName

		file, _ := os.Create(inputFilePath)
		_, _ = file.Write([]byte(t.In))
		file, _ = os.Create(outputFilePath)
		_, _ = file.Write([]byte(t.Out))
	}
}

func (j *Judge) start(problemId string, image string) {
	cli, _ := dclient.NewEnvClient()
	mount := NewProblemMount(problemId)
	spec := Spec{
		Name:   j.Name,
		Image:  image,
		Mounts: []ProblemMount{mount},
	}
	err := Start(cli, spec)
	if err != nil {
		log.Println("Error while starting judge: " + err.Error())
		return
	}
}

func (j *Judge) SendScript(submissionId string) {
	var scriptPath string = "./" + SubmissionsDirName + "/" + "submission-" + submissionId + ".py"
	Mkdir(j.Name, SubmissionsDirName)
	Send(j.Name, scriptPath, SubmissionsDirName + "/")
}