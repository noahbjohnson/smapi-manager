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

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	_ = r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("zip")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	var fileUrl = filepath.Join(getConfigDir(), handler.Filename)
	fmt.Println("Writing to: " + fileUrl)
	fileBytes, err := ioutil.ReadAll(file)
	errr := afero.WriteFile(AppFs, fileUrl, fileBytes, 0777)
	if errr != nil {
		panic(errr)
	}

	_, _ = fmt.Fprintf(w, "Successfully Uploaded File\n")
	return
}
