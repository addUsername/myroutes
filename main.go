package main

import (
	_ "embed"
	//"encoding/json"
	"github.com/wailsapp/wails"
)

// TODO
// - i think that u could just add cdn files (jquery, bootstrap, etc..) here like this 2 below
// - add counter and .json
// - export to .xlsx

//go:embed frontend/build/main.js
var js string

//go:embed frontend/build/main.css
var css string

func main() {

	app := wails.CreateApp(&wails.AppConfig{
		Width:  1024,
		Height: 768,
		Title:  "myroutes",
		JS:     js,
		CSS:    css,
		Colour: "#131313",
	})
	app.Bind(&Handler{})
	app.Run()
}
