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
	wnd := g.NewMasterWindow("Dartboard", 1000+ui.LeftToolbarMinimumWidth, 1000, 0)
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
