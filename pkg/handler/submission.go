package handler

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"pss/pkg"
)

func SubmitProblem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)

	problemId := params["id"]
	submissionId := pkg.GetRandomUUID()

	file, err := loadScriptFile(r)
	if err != nil {
		log.Println("Error while retrieving the File: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	err = writeScriptFile(submissionId, file)

	if err != nil {
		log.Println("Error while write script file")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		log.Println("Created submission [" + submissionId + "] to problem [" + problemId + "]")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{ "message": "Successful upload", "submission_id":"` +  submissionId + `" }`)); err != nil {
			log.Println(err.Error())
		}
	}
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

const ScriptFileNamePattern = "submission-"
const PythonExtension = ".py"
//const CPlusPlusExtension = ".cpp"

func writeScriptFile(submissionId string, f multipart.File) error {
	var fileName string = ScriptFileNamePattern + submissionId + PythonExtension
	var filePath string = pkg.SubmissionsDirName + "/" + fileName
	file, err := os.Create(filePath)

	if err != nil {
		return err
	}

	fileBytes, err := ioutil.ReadAll(f)
	file.Write(fileBytes)
	return err
}
