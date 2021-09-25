package main

import (
	_ "embed"
	"fmt"
	"os/exec"
	"strings"

	"github.com/wailsapp/wails"
)

// TODO
// - i think that u could just add cdn files (jquery, bootstrap, etc..) here like this 2 below, but idk how.. yet

//go:embed frontend/build/main.js
var js string

//go:embed frontend/build/main.css
var css string

func main() {

	fmt.Println("MAAAAIN")
	//get max window size
	out, err := exec.Command("cmd", "/C", "wmic desktopmonitor get screenheight,screenwidth").Output()
	if err != nil {
		fmt.Println(err)
	} else {
		info := strings.Split(string(out), "")
		fmt.Println(len(info))
		fmt.Println(info)
	}
	fmt.Println(string(out))
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
