package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"pss/pkg"
)

const CurrentVersion = "0.0.1"

func GetVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(pkg.Version{Tag: CurrentVersion}); err != nil {
		log.Println(err.Error())
	}
}

func RunProcessor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	spec := pkg.Spec{
		Name:  "judge-1",
		Image: "ubuntu",
		Mounts: []pkg.Mount{
			{
				Source:   "/home/wesley/go/src/pss/submissions",
				Target:   "/submissions",
				ReadOnly: false,
			},
		},
	}

	if err := pkg.Start(spec); err != nil {
		log.Println("Error while start container: " + err.Error())
	}

	if err := pkg.Exec(spec.Name, "ls"); err != nil {
		log.Println("Error while exec command: " + err.Error())
	}

	if err := pkg.Stop(spec.Name); err != nil {
		log.Println("Error while stop [" + spec.Name + "]: " + err.Error())
	}
	if err := json.NewEncoder(w).Encode(pkg.Version{Tag: CurrentVersion}); err != nil {
		log.Println(err.Error())
	}
}