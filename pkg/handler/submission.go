package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"pss/pkg"
)

const ScriptFileNamePattern = "submission-"
const PythonExtension = ".py"
//const CPlusPlusExtension = ".cpp"

var DefaultProcessor = pkg.NewProcessor(10)

func SubmitProblem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	problemId := mux.Vars(r)["id"]
	submissionId :=  primitive.NewObjectID()

	submission := pkg.Submission{
		ID:        submissionId,
		ProblemId: problemId, State: "CREATED"}

	_, err := pkg.SaveSubmission(submission)
	if err != nil {
		log.Println("Error while save submission")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var file multipart.File
	file, err = loadScriptFile(r)
	go submitToProcessor(file, submission)

	log.Println("Created submission [" + submissionId.Hex() + "] to problem [" + problemId + "]")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`{ "message": "Successful upload", "submission_id":"` +  submissionId.Hex() + `" }`)); err != nil {
		log.Println(err.Error())
	}
}

func RetrieveSubmission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	submissionId := params["id"]

	submission, err := pkg.RetrieveSubmission(submissionId)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte(`{ "message": "` + err.Error() + `" }`)); err != nil {
			log.Println(err.Error())
		}
	} else {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(submission); err != nil {
			log.Println(err.Error())
		}
	}
}

func submitToProcessor(file multipart.File, submission pkg.Submission) {
	err := writeScriptFile(file, submission.ID.Hex())
	if err != nil {
		log.Println("Error while write script file")
		_, _ = pkg.UpdateStateSubmission(submission.ID.Hex(), "ERROR")
	}
	DefaultProcessor.In <- submission
	_, _ = pkg.UpdateStateSubmission(submission.ID.Hex(), "SUBMITTED")
}

func writeScriptFile(multiPartFile multipart.File, submissionId string) error {
	var fileName string = ScriptFileNamePattern + submissionId + PythonExtension
	var filePath string = pkg.SubmissionsDirName + "/" + fileName
	file, err := os.Create(filePath)

	if err != nil {
		return err
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(multiPartFile)
	_, err = file.Write(fileBytes)
	return err
}

func loadScriptFile(r *http.Request) (multipart.File, error) {
	//10 MB
	_ = r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("script")
	if err != nil {
		return nil, err
	}
	log.Printf("Uploaded File: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)
	log.Printf("MIME Header: %+v\n", handler.Header)
	return file, nil
}
