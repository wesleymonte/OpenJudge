package main

import (
	"context"
	"flag"
	"github.com/gorilla/mux"
	"github.com/wesleymonte/openjudge/api/handler"
	"github.com/wesleymonte/openjudge/openjudge"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const GetVersionEndpoint = "/version"
const RegisterProblemEndpoint = "/problems"
const ProblemEndpoint = "/problems/{id}"
const GetSubmissionEndpoint = "/submissions/{id}"

func init() {
	log.Println("Starting pss...")
	openjudge.ValidateEnv()
	openjudge.CreateSubmissionsFolder()
}

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful_timeout", time.Second*15, "the duration for which the server "+
		"gracefully wait for existing connections to finish - e.g. 15s or 1m")

	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	openjudge.SetUp(ctx)

	router := mux.NewRouter()

	router.HandleFunc(GetVersionEndpoint, handler.GetVersion).Methods("GET")
	router.HandleFunc(RegisterProblemEndpoint, handler.RegisterProblem).Methods("POST")
	router.HandleFunc(ProblemEndpoint, handler.RetrieveProblem).Methods("GET")
	router.HandleFunc(ProblemEndpoint, handler.SubmitProblem).Methods("POST")
	router.HandleFunc(GetSubmissionEndpoint, handler.RetrieveSubmission).Methods("GET")

	server := &http.Server{
		Addr:         ":8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err.Error())
		}
	}()

	log.Println("Service available")

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	if err := server.Shutdown(ctx); err != nil {
		log.Println(err.Error())
	}

	log.Println("Shutting down service")

	os.Exit(1)
}
