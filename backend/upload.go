package backend

import (
	"archive/zip"
	"fmt"
	"github.com/spf13/afero"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const zipsFolder = "zips/"
const unzipsFolder = "unzipped/"

func parseUpload(r *http.Request) (fileUrl string, err error) {
	_ = r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("zip")
	if err != nil {
		return fileUrl, err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()
	Sugar.Debugf("Uploaded File: %+v\n", handler.Filename)
	Sugar.Debugf("File Size: %+v\n", handler.Size)
	Sugar.Debugf("MIME Header: %+v\n", handler.Header)

	fileUrl = filepath.Join(getConfigPathString(), zipsFolder, handler.Filename)
	Sugar.Info("Writing to: " + fileUrl)

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return fileUrl, err
	}

	err = afero.WriteFile(AppFs, fileUrl, fileBytes, 0777)
	if err != nil {
		return fileUrl, err
	}

	return fileUrl, nil
}

/*
1. Upload file to {config}/zips/ (UploadFile)
*/
func UploadFile(w http.ResponseWriter, r *http.Request) {
	Sugar.Info("File Upload Endpoint Hit")
	(w).Header().Set("Access-Control-Allow-Origin", "*")

	fileUrl, err := parseUpload(r)
	if err != nil {
		Sugar.Fatal(err)
	}

	err = unzipMod(fileUrl)
	if err != nil {
		Sugar.Fatal(err)
	}
	_, _ = fmt.Fprintf(w, "Successfully Uploaded File\n")
	return
}

func unzipMod(zipUrl string) error {
	zipReader, err := zip.OpenReader(zipUrl)
	if err != nil {
		return err
	}
	defer func() {
		if err := zipReader.Close(); err != nil {
			panic(err)
		}
	}()

	// check if there's a mod at the root of the zip
	//hasNoRootDir := false
	//for _, f := range zipReader.File {
	//	if !f.FileInfo().IsDir() {
	//		if _, filename := filepath.Split(f.Name); filename == manifestFileName ||
	//			filepath.Ext(f.Name) == "dll" {
	//			hasNoRootDir = true
	//		}
	//	}
	//}

	_, filename := filepath.Split(zipUrl)
	filename = strings.TrimRight(filename, ".zip")
	outDir := filepath.Join(getConfigPathString(), unzipsFolder, filename)

	err = os.MkdirAll(outDir, 0755)
	if err != nil {
		return err
	}

	// Closure to address file descriptors issue with all the deferred .Close() methods
	// from https://stackoverflow.com/questions/20357223/easy-way-to-unzip-file-with-golang
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(outDir, f.Name)

		if f.FileInfo().IsDir() {
			err = os.MkdirAll(path, 0755) // not using f.Mode() since it can cause a permission error
			if err != nil {
				return err
			}
		} else {
			err = os.MkdirAll(filepath.Dir(path), 0755)
			if err != nil {
				return err
			}
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()
			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range zipReader.File {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}
