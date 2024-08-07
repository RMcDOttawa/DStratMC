package ui

import (
	boardgeo "DStratMC/board-geometry"
	"DStratMC/simulation"
	"fmt"
	g "github.com/AllenDang/giu"
	"image"
	"math"
)

const (
	RadioExactScore = iota
	RadioOneAvgScore
	RadioGroupAvgScore
)

var radioValue int

const LeftToolbarMinimumWidth = 200

var accuracyModel simulation.AccuracyModel
var DartboardTexture *g.Texture
var dartboard Dartboard

func UserInterfaceSetup(theAccuracyModel simulation.AccuracyModel) {
	accuracyModel = theAccuracyModel
	radioValue = RadioOneAvgScore
	dartboard = NewDartboard(dartboardClickCallback)
}

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
	//fmt.Printf("Window position = (%g,%g), size = (%g,%g). Square image is %g x %g\n",
	//	windowX, windowY,
	//	windowWidth, windowHeight,
	//	squareDimension, squareDimension)
	dartboardImageMin := image.Pt(int(windowX)+leftToolbarWidth, int(windowY))
	dartboardImageMax := image.Pt(dartboardImageMin.X+int(squareDimension), dartboardImageMin.Y+int(squareDimension))
	//fmt.Printf("image min %d, max %d\n", imageMin, imageMax)

	dartboard.SetInfo(window, DartboardTexture, squareDimension, dartboardImageMin, dartboardImageMax)

	window.Layout(
		g.RadioButton("One Exact", radioValue == RadioExactScore).OnChange(func() { radioValue = RadioExactScore }),
		g.RadioButton("One Statistical", radioValue == RadioOneAvgScore).OnChange(func() { radioValue = RadioOneAvgScore }),
		g.RadioButton("Group Statistical", radioValue == RadioGroupAvgScore).OnChange(func() { radioValue = RadioGroupAvgScore }),
		g.Custom(dartboard.DrawFunction),
	)

}

func dartboardClickCallback(dartboard Dartboard, position boardgeo.BoardPosition) {
	//fmt.Printf("Dartboard clicked at radius %g, angle %g\n", position.Radius, position.Angle)
	if position.Radius <= 1.0 {
		if radioValue == RadioExactScore {
			//markHitPoint(polarRadius, thetaDegrees)
			_, score, description := boardgeo.DescribeBoardPoint(position)
			fmt.Printf("Single, exact: %s: %d points\n", description, score)
		} else if radioValue == RadioOneAvgScore {
			oneStatisticalThrow(dartboard, position, accuracyModel)
		} else if radioValue == RadioGroupAvgScore {
			fmt.Println("STUB group statistical score")
		} else {
			fmt.Println("Invalid radio button value")
		}
	}
}

func oneStatisticalThrow(dartboard Dartboard, position boardgeo.BoardPosition, model simulation.AccuracyModel) {
	//fmt.Printf("oneStatisticalThrow STUB %v,\n", position)

	//	Un-draw any previous markers or annotations
	dartboard.RemoveThrowMarkers()

	//  Draw a marker to record where we clicked
	dartboard.DrawTargetMarker(position)

	//	Draw a circle showing the accuracy radius around the clicked point
	accuracyRadius := model.GetAccuracyRadius()
	dartboard.DrawAccuracyCircle(position, accuracyRadius)

	//	Get a modeled hit within the accuracy
	//	Draw the hit within this circle
	//	Calculate the hit score
}
