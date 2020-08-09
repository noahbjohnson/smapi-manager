package main

import (
	"github.com/noahbjohnson/smapi-manager/backend"
	_ "github.com/noahbjohnson/smapi-manager/statik"
	"github.com/webview/webview"
	"log"
)

// TODO: customize or randomize
const addr = ":53494"

func main() {
	backend.StartAPI(addr)

	debug := true
	w := webview.New(debug)
	defer w.Destroy()
	err := backend.BindFunctions(w)
	if err != nil {
		log.Fatalln("Failed to bind functions to frontend", err)
	}
	w.SetTitle("Minimal webview example")
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("http://localhost" + addr)
	w.Run()
}
