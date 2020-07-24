package main

import (
	"github.com/leaanthony/mewn"
	"github.com/wailsapp/wails"
	"smapi-manager/backend"
)

func main() {

	js := mewn.String("./frontend/build/static/js/main.js")
	css := mewn.String("./frontend/build/static/css/main.css")

	app := wails.CreateApp(&wails.AppConfig{
		Width:  1024,
		Height: 768,
		Title:  "SMAPI Mod Manager",
		JS:     js,
		CSS:    css,
		Colour: "#131313",
	})
	app.Bind(backend.Basic)
	app.Bind(backend.Initialize)
	err := app.Run()
	if err != nil {
		panic(err)
	}
}
