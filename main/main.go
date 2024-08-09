package main

import (
	"DStratMC/ui"
	"fmt"
	g "github.com/AllenDang/giu"
)

func main() {
	wnd := g.NewMasterWindow("Dartboard", 1000+ui.LeftToolbarMinimumWidth, 1000, 0)
	loadedImage, err := g.LoadImage("./Dartboard Illustration.png")
	if err != nil {
		fmt.Println("Unable to load dartboard image:", err)
		return
	}
	ui.UserInterfaceSetup(loadedImage)
	wnd.Run(ui.MainUiLoop)

}
