package ui

import (
	boardgeo "DStratMC/board-geometry"
	g "github.com/AllenDang/giu"
	"image"
	"image/color"
)

const drawReferenceLines = true

var DartboardInfo struct {
	window          *g.WindowWidget
	Texture         *g.Texture
	squareDimension float64
	imageMin        image.Point
	imageMax        image.Point
	clickCallback   func(position boardgeo.BoardPosition)
	//sema            sync.Mutex
}

func SetDartboardDimensions(windowWidget *g.WindowWidget,
	squareDimension float64,
	imageMin image.Point, imageMax image.Point) {
	DartboardInfo.window = windowWidget
	DartboardInfo.squareDimension = squareDimension
	DartboardInfo.imageMin = imageMin
	DartboardInfo.imageMax = imageMax
}

func SetDartboardClickCallback(callback func(position boardgeo.BoardPosition)) {
	DartboardInfo.clickCallback = callback
}

func DartboardCustomFunc() {
	d := DartboardInfo
	if d.squareDimension == 0 {
		//fmt.Println("Squaredimension 0, returning")
		return
	}
	if d.imageMin.X < 0 || d.imageMin.Y < 0 || d.imageMax.X < 0 || d.imageMax.Y < 0 {
		//fmt.Println("imageMin or Max 0, returning")
		return
	}
	//fmt.Println("DartboardCustomFunc")
	//fmt.Printf("imageMin = %#v, imageMax = %#v\n", d.imageMin, d.imageMax)
	//fmt.Println("Square dimension:", d.squareDimension)
	canvas := g.GetCanvas()

	//	Basic test: draw a centred circle
	//radius := d.squareDimension / 2 * .5
	//stubCentre := image.Pt((d.imageMin.X+d.imageMax.X)/2, (d.imageMin.Y+d.imageMax.Y)/2)
	//stubRadius := float32(radius)
	//stubColour := color.RGBA{200, 0, 0, 255}
	//canvas.AddCircleFilled(stubCentre, stubRadius, stubColour)

	//	Position an invisible button on top of this image to detect clicks
	//	Remember and then restore drawing cursor so image comes out on top of this
	savedCsp := g.GetCursorScreenPos()
	g.SetCursorScreenPos(d.imageMin)
	g.InvisibleButton().Size(float32(d.squareDimension), float32(d.squareDimension)).
		OnClick(dartboardClicked).
		Build()
	g.SetCursorScreenPos(savedCsp)

	// Display dartboard image
	canvas.AddImage(d.Texture, d.imageMin, d.imageMax)

	if drawReferenceLines {
		//	For testing, draw a semitransparent circle on the centre
		xCentre := (d.imageMin.X + d.imageMax.X) / 2
		yCentre := (d.imageMin.Y + d.imageMax.Y) / 2
		testCirclePosition := image.Pt(xCentre, yCentre)
		testCircleColour := color.RGBA{R: 0, G: 0, B: 255, A: 128}
		radius := float32(d.squareDimension / 8.0)
		canvas.AddCircle(testCirclePosition, radius, testCircleColour, 0, 2)

		//	And add centred vertical and horizontal lines to help calibrate angle measurement
		crossHairColour := color.RGBA{R: 150, G: 150, B: 150, A: 255}

		verticalFrom := image.Pt(xCentre, yCentre-int(d.squareDimension/2))
		verticalTo := image.Pt(xCentre, yCentre+int(d.squareDimension/2))
		canvas.AddLine(verticalFrom, verticalTo, crossHairColour, 1)

		horizontalFrom := image.Pt(xCentre-int(d.squareDimension/2), yCentre)
		horizontalTo := image.Pt(xCentre+int(d.squareDimension/2), yCentre)
		canvas.AddLine(horizontalFrom, horizontalTo, crossHairColour, 1)
	}
}

func dartboardClicked() {
	//fmt.Println("dartboard clicked")
	if DartboardInfo.clickCallback == nil {
		//fmt.Println("  No callback function")
	} else {
		position := boardgeo.CalcMousePolarPosition(DartboardInfo.window)
		DartboardInfo.clickCallback(position)
	}
}
