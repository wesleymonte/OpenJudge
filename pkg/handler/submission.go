package handler

import (
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"pss/pkg"
)

func SubmitProblem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//10 MB
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("script")

	if err != nil {
		//SPECIFY HTTP CODE
		log.Println("Error while retrieving the File: " + err.Error())
		return
	}

	defer file.Close()
	log.Printf("Uploaded File: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)
	log.Printf("MIME Header: %+v\n", handler.Header)

	err = writeScriptFile("id", file)

	if err != nil {
		log.Println("Error while write script file")
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{ "message": "Successful upload" }`)); err != nil {
			log.Println(err.Error())
		}
	}
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
