package pkg

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

func (j Judge) Submit(s Submission) {
	//problem, err := RetrieveProblem(s.ProblemId)
	//if err != nil {
	//	log.Println("Error while retrieve problem [" + s.ProblemId + "]: " + err.Error())
	//}
}