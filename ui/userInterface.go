package ui

import (
	g "github.com/AllenDang/giu"
	"image"
	"math"
)

//const drawReferenceLines = true

func MainUiLoop() {
	window := g.SingleWindow()
	wx32, wy32 := window.CurrentPosition()
	windowX := float64(wx32)
	windowY := float64(wy32)
	//fmt.Printf("Window position = %g,%g\n", windowX, windowY)

	w32, h32 := window.CurrentSize()
	width := float64(w32)
	height := float64(h32)
	//fmt.Printf("Window size: %dx%d\n", int(width), int(height))

	squareDimension := math.Min(width, height)
	xPadding := int(math.Round((width - squareDimension) / 2))
	yPadding := int(math.Round((height - squareDimension) / 2))
	//fmt.Printf("Window position = (%g,%g), size = (%g,%g). Square image is %g x %g,  x padding %d, y padding %d\n",
	//	windowX, windowY,
	//	width, height,
	//	squareDimension, squareDimension,
	//	xPadding, yPadding)
	imageMin := image.Pt(int(windowX)+xPadding, int(windowY)+yPadding)
	imageMax := image.Pt(imageMin.X+int(squareDimension), imageMin.Y+int(squareDimension))
	//fmt.Printf("image min %d, max %d\n", imageMin, imageMax)

	SetDartboardDimensions(squareDimension, imageMin, imageMax)

	window.Layout(
		g.Custom(DartboardCustomFunc),
	)

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
