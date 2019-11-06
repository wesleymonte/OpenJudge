package pkg

import (
	"log"
	"time"
)

type Processor struct {
	In chan Submission
	out chan Submission
	jc chan Submission
	Judges []Judge
}

func (p *Processor) Run() {
	go p.runSubmitRoutine()
	go p.runUpdateRoutine()
}

func (p *Processor) startJudge(judgeName string) {
	for _, j := range p.Judges {
		name := j.Name
		if judgeName == name {
			//if id, err := j.Start(Spec{Image:"ubuntu"}); err != nil {
			//	log.Println("Error while start judge: " + err.Error())
		} else {
			//log.Println("Started judge [" + id + "]")
		}

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
			_, err := UpdateStateSubmission(s)
			if err != nil {
				log.Println("Error while update submission state [" + s.ProblemId + "]")
			}
		default:
			log.Println("Empty result queue.")
			time.Sleep(5 * time.Second)
		}
	}
}


