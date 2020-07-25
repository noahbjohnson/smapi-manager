package main

import (
	"github.com/webview/webview"
	"net/http"
	"smapi-manager/backend"
)

func setupRoutes() {
	http.HandleFunc("/upload", backend.UploadFile)
	http.Handle("/", http.FileServer(http.Dir("./frontend/build")))
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
