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
	RadioExact = iota
	RadioOneAvg
	RadioMultiAvg
	RadioOneNormal
	RadioMultiNormal
)

const LeftToolbarMinimumWidth = 200
const singleHitMarkerRadius = 5
const multipleHitMarkerRadius = 1

const ThrowsAtOneTarget = 1_000
const numThrowsTextWidth = 120

const uniformCEPRadius = 0.3
const normalCEPRadius = 0.25
const stubStandardDeviation = normalCEPRadius * 2

var radioValue int
var DartboardTexture *g.Texture
var dartboard Dartboard
var AccuracyModel simulation.AccuracyModel
var scoreDisplay string
var messageDisplay string
var throwTotal int64
var throwCount int64
var throwAverage float64
var numThrowsField int32 = ThrowsAtOneTarget
var drawReferenceLinesCheckbox = true

var drawOneSigma bool
var drawTwoSigma bool
var drawThreeSigma bool

func UserInterfaceSetup(loadedImage *image.RGBA) {
	radioValue = RadioOneNormal
	drawOneSigma = false
	drawTwoSigma = true
	drawThreeSigma = false
	g.EnqueueNewTextureFromRgba(loadedImage, func(t *g.Texture) {
		DartboardTexture = t
	})
	dartboard = NewDartboard(dartboardClickCallback)
	dartboard.SetDrawRefLines(drawReferenceLinesCheckbox)
}

func MainUiLoop() {
	window := setUpWindow()

	window.Layout(
		leftToolbarLayout(),
		g.Custom(dartboard.DrawFunction),
	)

}

func leftToolbarLayout() g.Widget {
	AccuracyModel = getAccuracyModel(radioValue)
	return g.Layout{
		//	Checkbox controlling whether crosshairs are drawn
		g.Checkbox("Reference Lines", &drawReferenceLinesCheckbox).OnChange(func() { dartboard.SetDrawRefLines(drawReferenceLinesCheckbox) }),

		// Radio buttons to select the type of interaction and model
		g.Label(""),
		g.RadioButton("One Exact", radioValue == RadioExact).OnChange(func() { radioValue = RadioExact; AccuracyModel = getAccuracyModel(radioValue); radioChanged() }),
		g.RadioButton("One Throw Uniform", radioValue == RadioOneAvg).OnChange(func() { radioValue = RadioOneAvg; AccuracyModel = getAccuracyModel(radioValue); radioChanged() }),
		g.RadioButton("Multi Throw Uniform", radioValue == RadioMultiAvg).OnChange(func() { radioValue = RadioMultiAvg; AccuracyModel = getAccuracyModel(radioValue); radioChanged() }),
		g.RadioButton("One Throw Normal", radioValue == RadioOneNormal).OnChange(func() { radioValue = RadioOneNormal; AccuracyModel = getAccuracyModel(radioValue); radioChanged() }),
		g.RadioButton("Multi Throw Normal", radioValue == RadioMultiNormal).OnChange(func() { radioValue = RadioMultiNormal; AccuracyModel = getAccuracyModel(radioValue); radioChanged() }),

		// A reset button resets counters and displays
		g.Label(""),
		g.Button("Reset").OnClick(radioChanged),

		//	Display a generic message and throw score, if present
		g.Condition(messageDisplay != "" || scoreDisplay != "",
			g.Layout{
				g.Label(""),
				g.Label(messageDisplay),
				g.Label(scoreDisplay),
			}, nil),

		// If we are doing multiple throws, allow the user to set the number of throws
		g.Condition(radioValue == RadioMultiAvg || radioValue == RadioMultiNormal,
			g.Layout{
				g.Label(""),
				g.InputInt(&numThrowsField).Label("# Throws").
					Size(numThrowsTextWidth).
					StepSize(1).
					StepSizeFast(100),
			}, nil),

		// When we are doing normal distribution (and only then) offer 3 checkboxes for drawing reference
		// circles at 1, 2, and 3 standard deviations
		g.Condition(radioValue == RadioOneNormal || radioValue == RadioMultiNormal,
			g.Layout{
				g.Label(""),
				g.Checkbox("1 Sigma", &drawOneSigma),
				g.Checkbox("2 Sigma", &drawTwoSigma),
				g.Checkbox("3 Sigma", &drawThreeSigma),
			}, nil),

		// Once a number of throws have been accumulated, display the average score
		g.Condition(throwCount > 0,
			g.Layout{
				g.Label(""),
				g.Label("Throws: " + strconv.Itoa(int(throwCount))),
				g.Label("Total: " + strconv.Itoa(int(throwTotal))),
				g.Label("Average: " + strconv.FormatFloat(throwAverage, 'f', 1, 64)),
			},
			nil),
	}

}

func getAccuracyModel(radioValue int) simulation.AccuracyModel {
	switch radioValue {
	case RadioExact:
		return nil
	case RadioOneAvg:
		return simulation.NewUniformAccuracyModel(uniformCEPRadius)
	case RadioMultiAvg:
		return simulation.NewUniformAccuracyModel(uniformCEPRadius)
	case RadioOneNormal:
		return simulation.NewNormalAccuracyModel(normalCEPRadius, stubStandardDeviation)
	case RadioMultiNormal:
		return simulation.NewNormalAccuracyModel(normalCEPRadius, stubStandardDeviation)
	default:
		panic("Invalid radio button value")
		return simulation.NewPerfectAccuracyModel()
	}
}

func setUpWindow() *g.WindowWidget {
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
	return window
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
		switch radioValue {
		case RadioExact:
			dartboard.QueueTargetMarker(position)
			_, score, description := boardgeo.DescribeBoardPoint(position)
			messageDisplay = description
			scoreDisplay = strconv.Itoa(score) + " points"
		case RadioOneAvg:
			oneUniformThrow(dartboard, position, AccuracyModel)
		case RadioMultiAvg:
			multipleUniformThrows(dartboard, position, AccuracyModel)
		case RadioOneNormal:
			oneNormalThrow(dartboard, position, AccuracyModel)
		case RadioMultiNormal:
			multipleNormalThrows(dartboard, position, AccuracyModel)
		default:
			panic("Invalid radio button value")
		}
	}
}

func oneUniformThrow(dartboard Dartboard, position boardgeo.BoardPosition, model simulation.AccuracyModel) {
	//fmt.Printf("oneUniformThrow  %v,\n", position)

	//  Draw a marker to record where we clicked
	dartboard.QueueTargetMarker(position)

	//	Draw a circle showing the accuracy radius around the clicked point
	accuracyRadius := model.GetAccuracyRadius()
	dartboard.QueueAccuracyCircle(position, accuracyRadius)

	//	Get a modeled hit within the accuracy
	hit, err := model.GetThrow(position,
		dartboard.GetScoringRadiusPixels(),
		dartboard.GetSquareDimension(),
		dartboard.GetImageMinPoint())
	if err != nil {
		fmt.Printf("Error getting throw %v", err)
		return
	}
	//fmt.Printf("Hit: %#v \n", hit)

	//	Draw the hit within this circle
	dartboard.QueueHitMarker(hit, singleHitMarkerRadius)

	//	Calculate the hit score
	_, score, description := boardgeo.DescribeBoardPoint(hit)
	messageDisplay = description
	scoreDisplay = strconv.Itoa(score) + " points"
	throwCount++
	throwTotal += int64(score)
	throwAverage = float64(throwTotal) / float64(throwCount)
	g.Update()

}

func multipleUniformThrows(dartboard Dartboard, position boardgeo.BoardPosition, model simulation.AccuracyModel) {
	//fmt.Printf("multipleUniformThrows  %v,\n", position)

	//  Draw a marker to record where we clicked
	dartboard.QueueTargetMarker(position)

	//	Draw a circle showing the accuracy radius around the clicked point
	accuracyRadius := model.GetAccuracyRadius()
	dartboard.QueueAccuracyCircle(position, accuracyRadius)
	dartboard.AllocateHitsSpace(int(numThrowsField))

	throwCount = 0
	throwTotal = 0
	for i := 0; i < int(numThrowsField); i++ {
		//	Get a modeled hit within the accuracy
		hit, err := model.GetThrow(position,
			dartboard.GetScoringRadiusPixels(),
			dartboard.GetSquareDimension(),
			dartboard.GetImageMinPoint())
		if err != nil {
			fmt.Printf("Error getting throw %v", err)
			return
		}
		//fmt.Printf("Hit: %#v \n", hit)

		//	Draw the hit within this circle
		dartboard.QueueHitMarker(hit, multipleHitMarkerRadius)

		//	Calculate the hit score
		_, score, _ := boardgeo.DescribeBoardPoint(hit)
		throwCount++
		throwTotal += int64(score)
		throwAverage = float64(throwTotal) / float64(throwCount)
	}
	g.Update()

}

func oneNormalThrow(dartboard Dartboard, position boardgeo.BoardPosition, model simulation.AccuracyModel) {
	//fmt.Printf("oneNormalThrow  %v,\n", position)

	//  Draw a marker to record where we clicked
	dartboard.QueueTargetMarker(position)

	//	Draw circles showing requested sigma levels around the clicked point
	if drawOneSigma {
		oneSigmaRadius := model.GetSigmaRadius(1)
		dartboard.QueueStdDeviationCircle(position, 1, oneSigmaRadius)
	}
	if drawTwoSigma {
		twoSigmaRadius := model.GetSigmaRadius(2)
		dartboard.QueueStdDeviationCircle(position, 2, twoSigmaRadius)
	}
	if drawThreeSigma {
		threeSigmaRadius := model.GetSigmaRadius(3)
		dartboard.QueueStdDeviationCircle(position, 3, threeSigmaRadius)
	}

	//	Get a modeled hit within the accuracy
	hit, err := model.GetThrow(position,
		dartboard.GetScoringRadiusPixels(),
		dartboard.GetSquareDimension(),
		dartboard.GetImageMinPoint())
	if err != nil {
		fmt.Printf("Error getting throw %v", err)
		return
	}
	//fmt.Printf("Hit: %#v \n", hit)

	//	Draw the hit within this circle
	dartboard.QueueHitMarker(hit, singleHitMarkerRadius)

	//	Calculate the hit score
	_, score, description := boardgeo.DescribeBoardPoint(hit)
	messageDisplay = description
	scoreDisplay = strconv.Itoa(score) + " points"
	throwCount++
	throwTotal += int64(score)
	throwAverage = float64(throwTotal) / float64(throwCount)
	g.Update()

}

func multipleNormalThrows(dartboard Dartboard, position boardgeo.BoardPosition, model simulation.AccuracyModel) {
	//fmt.Printf("oneNormalThrow  %v,\n", position)

	//  Draw a marker to record where we clicked
	dartboard.QueueTargetMarker(position)

	//	Draw circles showing requested sigma levels around the clicked point
	if drawOneSigma {
		oneSigmaRadius := model.GetSigmaRadius(1)
		dartboard.QueueStdDeviationCircle(position, 1, oneSigmaRadius)
	}
	if drawTwoSigma {
		twoSigmaRadius := model.GetSigmaRadius(2)
		dartboard.QueueStdDeviationCircle(position, 2, twoSigmaRadius)
	}
	if drawThreeSigma {
		threeSigmaRadius := model.GetSigmaRadius(3)
		dartboard.QueueStdDeviationCircle(position, 3, threeSigmaRadius)
	}
	dartboard.AllocateHitsSpace(int(numThrowsField))

	throwCount = 0
	throwTotal = 0

	for i := 0; i < int(numThrowsField); i++ {
		//	Get a modeled hit within the accuracy
		hit, err := model.GetThrow(position,
			dartboard.GetScoringRadiusPixels(),
			dartboard.GetSquareDimension(),
			dartboard.GetImageMinPoint())
		if err != nil {
			fmt.Printf("Error getting throw %v", err)
			return
		}
		//fmt.Printf("Hit: %#v \n", hit)

		//	Draw the hit within this circle
		dartboard.QueueHitMarker(hit, multipleHitMarkerRadius)

		//	Calculate the hit score
		_, score, description := boardgeo.DescribeBoardPoint(hit)
		messageDisplay = description
		scoreDisplay = strconv.Itoa(score) + " points"
		throwCount++
		throwTotal += int64(score)
		throwAverage = float64(throwTotal) / float64(throwCount)
	}
	g.Update()

}
