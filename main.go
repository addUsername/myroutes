package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/wailsapp/wails"
)

//go:embed frontend/build/main.js
var js string

//go:embed frontend/build/main.css
var css string

var DAYS int = 30
var routes_json string = "\\routes.json"
var ext_excel string = "\\PROGRAMACION-UPDATED.xlsx"
var int_excel string = "\\PROGRAMACION.xlsx"
var path string

func main() {

	fmt.Println("MAAAAIN")
	// get path()
	path, _ = os.Getwd()

	app := wails.CreateApp(&wails.AppConfig{
		MinWidth:  1500,
		MinHeight: 768,
		MaxHeight: 1080,
		MaxWidth:  1920,
		Title:     "myroutes",
		JS:        js,
		CSS:       css,
		Resizable: true,
		Colour:    "#131313",
	})
	app.Bind(&Handler{})
	app.Run()
}
