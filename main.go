package main

import (
	"github.com/rakyll/statik/fs"
	"github.com/webview/webview"
	"log"
	"net/http"
	"smapi-manager/backend"
	_ "smapi-manager/statik"
)

func setupRoutes() {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/upload", backend.UploadFile)
	http.Handle("/", http.FileServer(statikFS))
}

func main() {
	const addr = ":53494"
	setupRoutes()

	go func() {
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			panic(err)
		}
	}()

	debug := true
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("Minimal webview example")
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("http://localhost" + addr)
	w.Run()
}
