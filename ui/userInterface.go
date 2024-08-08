package ui

import (
	boardgeo "DStratMC/board-geometry"
	"DStratMC/simulation"
	"fmt"
	g "github.com/AllenDang/giu"
	"image"
	"math"
	"strconv"
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

var scoreDisplay string
var messageDisplay string

var throwTotal int64
var throwCount int64
var throwAverage float64

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
		g.RadioButton("One Exact", radioValue == RadioExactScore).OnChange(func() { radioValue = RadioExactScore; radioChanged() }),
		g.RadioButton("One Statistical", radioValue == RadioOneAvgScore).OnChange(func() { radioValue = RadioOneAvgScore; radioChanged() }),
		g.RadioButton("Group Statistical", radioValue == RadioGroupAvgScore).OnChange(func() { radioValue = RadioGroupAvgScore; radioChanged() }),
		g.Label(""),
		g.Button("Reset").OnClick(radioChanged),
		g.Label(""),
		g.Label(messageDisplay),
		g.Label(scoreDisplay),
		g.Condition(throwCount > 0,
			g.Layout{
				g.Label(""),
				g.Label("Throws: " + strconv.Itoa(int(throwCount))),
				g.Label("Total: " + strconv.Itoa(int(throwTotal))),
				g.Label("Average: " + strconv.FormatFloat(throwAverage, 'f', 1, 64)),
			},
			nil),
		g.Custom(dartboard.DrawFunction),
	)

}

func radioChanged() {
	scoreDisplay = ""
	messageDisplay = ""
	throwTotal = 0
	throwCount = 0
	throwAverage = 0
	dartboard.RemoveThrowMarkers()
}

func dartboardClickCallback(dartboard Dartboard, position boardgeo.BoardPosition) {
	//fmt.Printf("Dartboard clicked at radius %g, angle %g\n", position.Radius, position.Angle)
	if position.Radius <= 1.0 {
		messageDisplay = ""
		scoreDisplay = ""
		dartboard.RemoveThrowMarkers()
		if radioValue == RadioExactScore {
			//markHitPoint(polarRadius, thetaDegrees)
			dartboard.DrawTargetMarker(position)
			_, score, description := boardgeo.DescribeBoardPoint(position)
			messageDisplay = description
			scoreDisplay = strconv.Itoa(score) + " points"
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
	//fmt.Printf("oneStatisticalThrow  %v,\n", position)

	//  Draw a marker to record where we clicked
	dartboard.DrawTargetMarker(position)

	//	Draw a circle showing the accuracy radius around the clicked point
	accuracyRadius := model.GetAccuracyRadius()
	dartboard.DrawAccuracyCircle(position, accuracyRadius)

	//	Get a modeled hit within the accuracy
	hit, err := accuracyModel.GetThrow(position,
		dartboard.GetScoringRadiusPixels(),
		dartboard.GetSquareDimension(),
		dartboard.GetImageMinPoint())
	if err != nil {
		fmt.Printf("Error getting throw %v", err)
		return
	}
	//fmt.Printf("Hit: %#v \n", hit)

	//	Draw the hit within this circle
	dartboard.AddHitMarker(hit)

	//	Calculate the hit score
	_, score, description := boardgeo.DescribeBoardPoint(hit)
	messageDisplay = description
	scoreDisplay = strconv.Itoa(score) + " points"
	throwCount++
	throwTotal += int64(score)
	throwAverage = float64(throwTotal) / float64(throwCount)
	g.Update()

	//	Add a second hit just to see if it works
	//hit, err = accuracyModel.GetThrow(position,
	//	dartboard.GetScoringRadiusPixels(),
	//	dartboard.GetSquareDimension(),
	//	dartboard.GetImageMinPoint())
	//if err != nil {
	//	fmt.Printf("Error getting throw %v", err)
	//	return
	//}
	////fmt.Printf("Hit: %#v \n", hit)
	//
	////	Draw the hit within this circle
	//dartboard.AddHitMarker(hit)
	//g.Update()

}
