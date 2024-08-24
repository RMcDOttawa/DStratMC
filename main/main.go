package main

import (
	"DStratMC/ui"
	"bytes"
	_ "embed"
	g "github.com/AllenDang/giu"
	"image"
)

//go:embed style.css
var cssStyle []byte

//go:embed DartboardIllustration.png
var imageBytes []byte

func main() {
	wnd := g.NewMasterWindow("Dartboard", ui.MasterWindowWidth, ui.MasterWindowHeight, 0)
	wnd.SetSizeLimits(ui.MasterWindowWidth, ui.MasterWindowHeight, 8000, 8000)
	if err := g.ParseCSSStyleSheet(cssStyle); err != nil {
		panic(err)
	}
	//fmt.Println("Embedded # bytes", len(imageBytes))
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		panic(err)
	}
	//fmt.Println("Converted to image:", img.Bounds())
	loadedImage := g.ImageToRgba(img)
	userInterface := ui.NewUserInterface(loadedImage)
	wnd.Run(userInterface.MainUiLoop)

}
