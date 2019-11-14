package pkg

import (
	"fmt"
	"log"
)

type Processor struct {
	In chan Submission
	out chan Status
	jc chan Submission
	Judges []Judge
}

func NewProcessor(judges int) (p *Processor) {
	var in =  make(chan Submission)
	var out = make(chan Status)
	var js []Judge
	for i := 0; i < judges; i++ {
		log.Println(fmt.Sprintf("Creating judge-%d", i + 1))
		j := &Judge{
			Name: fmt.Sprintf("judge-%d", i + 1),
		}
		log.Println("Starting judge [" + j.Name + "]")
		go j.Run(in, out)
		js = append(js, *j)
	}
	p = &Processor{
		In:     in,
		out:    out,
		Judges: js,
	}
	go updateRoutine(p.out)
	return
}

func updateRoutine(out <- chan Status) {
	for s:= range out {
		_, err := UpdateStateSubmission(s.SubmissionId, s.Result)
		if err != nil {
			log.Println("Error while update submission state [" + s.SubmissionId + "]")
		}
	}
}


