package main

import (
	"context"
	"flag"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"pss/pkg"
	"pss/pkg/handler"
	"time"
)

const GetVersionEndpoint = "/version"
const RegisterProblemEndpoint = "/problems"
const GetProblemEndpoint = "/problems/{id}"
const SubmitProblemEndpoint = "/problems/{id}/submissions"
const GetSubmissionEndpoint = "/problems/{id}/submissions/{s_id}"

func init() {
	log.Println("Starting pss...")
	pkg.ValidateEnv()
	pkg.CreateSubmissionsFolder()
	handler.Init()
}

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful_timeout", time.Second*15, "the duration for which the server "+
		"gracefully wait for existing connections to finish - e.g. 15s or 1m")

	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	pkg.SetUp(ctx)

	router := mux.NewRouter()

	router.HandleFunc(GetVersionEndpoint, handler.GetVersion).Methods("GET")
	router.HandleFunc(RegisterProblemEndpoint, handler.RegisterProblem).Methods("POST")
	router.HandleFunc(GetProblemEndpoint, handler.RetrieveProblem).Methods("GET")
	router.HandleFunc(SubmitProblemEndpoint, handler.SubmitProblem).Methods("POST")

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
