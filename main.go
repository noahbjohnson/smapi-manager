package main

import (
	"github.com/leaanthony/mewn"
	"github.com/noahbjohnson/smapi-manager/backend"
	"github.com/wailsapp/wails"
	"log"
)

// TODO: customize or randomize
const addr = ":53494"
const title = "Stardew Mod Manager"

func newWebview() (app *wails.App, err error) {
	js := mewn.String("./frontend/build/static/js/main.js")
	css := mewn.String("./frontend/build/static/css/main.css")
	app = wails.CreateApp(&wails.AppConfig{
		Width:  1024,
		Height: 768,
		Title:  title,
		JS:     js,
		CSS:    css,
		Colour: "#131313",
	})
	return app, err
}

func openSmapiInstall() error {
	return backend.OpenSmapiInstall()
}

func hasSmapi() (bool, error) {
	return backend.HasSMAPI()
}

func main() {
	// Load defaults and/or config file
	backend.InitializeConfig()

	// Start the API server in the background
	err := backend.StartAPI(addr)
	if err != nil {
		log.Fatalln("Failed to start API", err)
	}

	// Create the webview window
	app, err := newWebview()
	if err != nil {
		log.Fatalln("Failed to bind functions to frontend", err)
	}

	app.Bind(openSmapiInstall)
	app.Bind(hasSmapi)

	err = app.Run()
}
