package main

import (
	"DStratMC/ui"
	_ "embed"
	"fmt"
	g "github.com/AllenDang/giu"
)

//go:embed style.css
var cssStyle []byte

func main() {
	wnd := g.NewMasterWindow("Dartboard", ui.MasterWindowWidth, ui.MasterWindowHeight, 0)
	wnd.SetSizeLimits(ui.MasterWindowWidth, ui.MasterWindowHeight, 8000, 8000)
	if err := g.ParseCSSStyleSheet(cssStyle); err != nil {
		panic(err)
	}
	loadedImage, err := g.LoadImage("./Dartboard Illustration.png")
	if err != nil {
		fmt.Println("Unable to load dartboard image:", err)
		return
	}
	userInterface := ui.NewUserInterface(loadedImage)
	wnd.Run(userInterface.MainUiLoop)

}
