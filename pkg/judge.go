package pkg

import (
	"fmt"
	dclient "github.com/docker/docker/client"
	"log"
	"os"
)

type Judge struct {
	Name string
	in chan Submission
	out chan Status
}

type Status struct {
	SubmissionId string
	Result string
}

func New(name string, in chan Submission, out chan Status) *Judge {
	j := new(Judge)
	j.Name = name
	j.in = in
	j.out = out
	return j
}

func (j *Judge) Run() {
	for {
		select {
		case s := <-j.in:
			status, err := j.submit(&s)
			if err != nil {
				log.Println("Error while judge submission [" + s.ID.Hex() + "]: " + err.Error())
			} else {
				j.out <- status
			}
		}
	}
}

func (j *Judge) submit(s *Submission) (status Status, err error){
	cli, _ := dclient.NewEnvClient()
	var problem *Problem
	problem, err = RetrieveProblem(s.ProblemId)
	if err != nil {
		log.Println("Error while retrieve problem [" + s.ProblemId + "]: " + err.Error())
	}
	if err = j.loadProblem(problem); err != nil {return}
	if err = j.start(cli, problem.ID.Hex(), "ubuntu"); err != nil {return}
	if err = j.sendScript(s.ID.Hex()); err != nil {return}
	if result, err := j.runScript(problem.ID.Hex(), s.ID.Hex()); err == nil {
		status = Status{
			SubmissionId: s.ID.Hex(),
			Result:       result,
		}
	}
	if err = j.stop(cli); err != nil {return}
	return
}

func (j *Judge) loadProblem(problem *Problem) (err error){
	dir := fmt.Sprintf("./problems/%s", problem.ID.Hex())
	inputDir := dir + "/in"
	outputDir := dir + "/out"
	if err = CreateFolders(inputDir, outputDir); err != nil {
		log.Println("Error while creating input/output directory")
		return
	}
	for i, t := range problem.TestCases {
		var inputFileName = fmt.Sprintf("%d.in", i + 1)
		var inputFilePath = inputDir + "/" + inputFileName

		var outputFileName = fmt.Sprintf("%d.out", i + 1)
		var outputFilePath = outputDir + "/" + outputFileName

		file, _ := os.Create(inputFilePath)
		_, _ = file.Write([]byte(t.In))
		file, _ = os.Create(outputFilePath)
		_, _ = file.Write([]byte(t.Out))
	}
	return
}

func (j *Judge) start(cli *dclient.Client, problemId string, image string) (err error){
	mount := NewProblemMount(problemId)
	spec := Spec{
		Name:   j.Name,
		Image:  image,
		Mounts: []ProblemMount{mount},
	}
	if err = Start(cli, spec); err != nil {
		log.Println("Error while starting judge")
	}
	return
}

func (j *Judge) sendScript(submissionId string) (err error) {
	var scriptPath = "./" + SubmissionsDirName + "/" + "submission-" + submissionId + ".py"
	if err = Mkdir(j.Name, SubmissionsDirName); err != nil {
		log.Println("Error while creating submissions directory to judge [" + j.Name + "]")
	} else {
		if err = Send(j.Name, scriptPath, SubmissionsDirName + "/"); err != nil {
			log.Println("Error while send script file to submissions directory from judge [" + j.Name + "]")
		}
	}
	return
}

func (j *Judge) runScript(problemId, submissionId string) (result string, err error) {
	var bresult []byte
	bresult, err = Run(j.Name, problemId, submissionId)
	result = string(bresult)
	return
}

func (j *Judge) stop(cli *dclient.Client) (err error) {
	err = Stop(cli, j.Name)
	return
}