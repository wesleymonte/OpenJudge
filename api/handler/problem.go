package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/wesleymonte/openjudge/openjudge"
	"log"
	"net/http"
)

func RegisterProblem( w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var p openjudge.Problem

	err := json.NewDecoder(r.Body).Decode(&p)

	if err != nil {
		log.Println("Error while decode body to problem")
	}

	res, err := openjudge.SaveProblem(p)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(`{ "message": "` + err.Error() + `" }`)); err != nil {
			log.Println(err.Error())
		}
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println(err.Error())
	}
}

func RetrieveProblem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	problemId := params["id"]

	problem, err := openjudge.RetrieveProblem(problemId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(`{ "message": "` + err.Error() + `" }`)); err != nil {
			log.Println(err.Error())
		}
	} else {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(problem); err != nil {
			log.Println(err.Error())
		}
	}
}