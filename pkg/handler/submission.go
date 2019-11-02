package handler

import (
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

func SubmitProblem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	problemId := mux.Vars(r)["id"]
	submissionId :=  primitive.NewObjectID()

	err := writeScriptFile(r, submissionId.Hex())

	if err != nil {
		log.Println("Error while write script file")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	submission := pkg.Submission{
		ID:        submissionId,
		ProblemId: problemId, State: "CREATED"}

	_, err = pkg.SaveSubmission(submission)

	if err != nil {
		log.Println("Error while save submission")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println("Created submission [" + submissionId.Hex() + "] to problem [" + problemId + "]")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`{ "message": "Successful upload", "submission_id":"` +  submissionId.Hex() + `" }`)); err != nil {
		log.Println(err.Error())
	}
}

func writeScriptFile(r *http.Request, submissionId string) error {
	multiPartFile, err := loadScriptFile(r)

	if err != nil {
		return err
	}

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
