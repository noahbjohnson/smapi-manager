package main

import (
	"github.com/noahbjohnson/smapi-manager/backend"
	_ "github.com/noahbjohnson/smapi-manager/statik"
	"github.com/rakyll/statik/fs"
	"github.com/webview/webview"
	"log"
	"net/http"
)

func setupRoutes() {
	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/upload", backend.UploadFile)
	http.HandleFunc("/smapi", backend.GetSMAPI)
	http.HandleFunc("/mods", backend.EnumerateMods)
	http.Handle("/", http.FileServer(statikFS))
}

func main() {
	const addr = ":53494"
	backend.Initialize()
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
