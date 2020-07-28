package backend

import (
	"fmt"
	"github.com/spf13/afero"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	_ = r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("zip")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	Sugar.Debugf("Uploaded File: %+v\n", handler.Filename)
	Sugar.Debugf("File Size: %+v\n", handler.Size)
	Sugar.Debugf("MIME Header: %+v\n", handler.Header)

	var fileUrl = filepath.Join(getConfigPathString(), "zips/", handler.Filename)
	fmt.Println("Writing to: " + fileUrl)
	fileBytes, err := ioutil.ReadAll(file)
	errr := afero.WriteFile(AppFs, fileUrl, fileBytes, 0777)
	if errr != nil {
		panic(errr)
	}

	_, _ = fmt.Fprintf(w, "Successfully Uploaded File\n")
	return
}
