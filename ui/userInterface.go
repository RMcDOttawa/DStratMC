package ui

import (
	boardgeo "DStratMC/board-geometry"
	"fmt"
	g "github.com/AllenDang/giu"
	"image"
	"math"
)

// const drawReferenceLines = true
const LeftToolbarMinimumWidth = 200

func MainUiLoop() {
	window := g.SingleWindow()
	wx32, wy32 := window.CurrentPosition()
	windowX := float64(wx32)
	windowY := float64(wy32)
	//fmt.Printf("Window position = %g,%g\n", windowX, windowY)

	w32, h32 := window.CurrentSize()
	windowWidth := float64(w32)
	windowHeight := float64(h32)
	leftToolbarWidth := int(math.Max(windowWidth-windowHeight, float64(LeftToolbarMinimumWidth)))
	dartboardWidth := int(windowWidth) - leftToolbarWidth
	//fmt.Printf("Window size: %dx%d\n", int(width), int(height))

	// There is a left toolbar with buttons and messages, and the dartboard occupies a square
	// in the remaining window to the right of this

	squareDimension := math.Min(float64(dartboardWidth), windowHeight)
	//xPadding := leftToolbarWidth
	//yPadding := 0 // dartboard is at top of window - no padding above it
	//fmt.Printf("Window position = (%g,%g), size = (%g,%g). Square image is %g x %g,  x padding %d, y padding %d\n",
	//	windowX, windowY,
	//	width, height,
	//	squareDimension, squareDimension,
	//	xPadding, yPadding)
	dartboardImageMin := image.Pt(int(windowX)+leftToolbarWidth, int(windowY))
	dartboardImageMax := image.Pt(dartboardImageMin.X+int(squareDimension), dartboardImageMin.Y+int(squareDimension))
	//fmt.Printf("image min %d, max %d\n", imageMin, imageMax)

	SetDartboardDimensions(window, squareDimension, dartboardImageMin, dartboardImageMax)
	SetDartboardClickCallback(dartboardClickCallback)

	window.Layout(
		g.Custom(DartboardCustomFunc),
	)

}

func dartboardClickCallback(position boardgeo.BoardPosition) {
	//fmt.Printf("Dartboard clicked at radius %g, angle %g\n", position.Radius, position.Angle)
	if position.Radius <= 1.0 {
		//markHitPoint(polarRadius, thetaDegrees)
		_, score, description := boardgeo.DescribeBoardPoint(position)
		fmt.Printf("%s: %d points\n", description, score)
	}
}

//func MainUiLoop() {
//	window := g.SingleWindow()
//	width, height := window.CurrentSize()
//	//fmt.Printf("Window size: %dx%d\n", int(width), int(height))
//	squareDimension := math.Min(float64(width), float64(height))
//	//fmt.Printf("Square image is %d x %[1]d\n", int(squareDimension))
//	window.Layout(
//		g.Align(g.AlignCenter).To(
//			g.ImageWithFile("Dartboard Illustration.png").
//				OnClick(func() {
//					polarRadius, thetaDegrees := boardgeo.CalcMousePolarPosition(squareDimension)
//					position := boardgeo.BoardPosition{
//						Radius: polarRadius,
//						Angle:  thetaDegrees,
//					}
//					if polarRadius <= 1.0 {
//						//markHitPoint(polarRadius, thetaDegrees)
//						description, score := boardgeo.DescribeBoardPoint(position)
//						fmt.Printf("%s: %d points\n", description, score)
//					}
//				}).
//				Size(float32(squareDimension), float32(squareDimension)),
//		),
//	)
//}
