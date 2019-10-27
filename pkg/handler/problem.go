package handler

import (
	"github.com/gorilla/mux"
	"log"
	"pss/pkg"
	"net/http"
	"encoding/json"
)

func RegisterProblem( w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var p pkg.Problem

	err := json.NewDecoder(r.Body).Decode(&p)

	if err != nil {
		log.Println("Error while decode body to problem")
	}

	res, err := pkg.SaveProblem(p)

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

	problem, err := pkg.RetrieveProblem(problemId)

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