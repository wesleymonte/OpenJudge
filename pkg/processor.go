package pkg

import (
	"fmt"
	"log"
	"time"
)

type Processor struct {
	In chan Submission
	out chan Status
	jc chan Submission
	Judges []Judge
}

func NewProcessor(judges int) (p *Processor) {
	var In, jc = make(chan Submission), make(chan Submission)
	var out = make(chan Status)
	var js []Judge
	for i := 0; i < judges; i++ {
		j := Judge{
			Name: fmt.Sprintf("judge-%d", i + 1),
			in:   jc,
			out:  out,
		}
		js = append(js, j)
	}
	p = &Processor{
		In:     In,
		out:    out,
		jc:     jc,
		Judges: js,
	}
	return
}

func (p *Processor) Run() {
	go p.wakeUpJudges()
	go p.runSubmitRoutine()
	go p.runUpdateRoutine()
}

func (p *Processor) wakeUpJudges() {
	for _, j := range p.Judges {
		log.Println("Starting judge [" + j.Name + "]")
		j.Run()
		log.Println("Started judge [" + j.Name + "]")
	}
}

func (p *Processor) runSubmitRoutine() {
	for {
		select {
		case s := <- p.In:
			p.jc <- s
		default:
			log.Println("Empty queue.")
			time.Sleep(5 * time.Second)
		}
	}
}

func (p *Processor) runUpdateRoutine() {
	for {
		select {
		case s := <- p.out:
			_, err := UpdateStateSubmission(s.SubmissionId, s.Result)
			if err != nil {
				log.Println("Error while update submission state [" + s.SubmissionId + "]")
			}
		default:
			log.Println("Empty result queue.")
			time.Sleep(5 * time.Second)
		}
	}
}


