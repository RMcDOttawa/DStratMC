package ui

import (
	boardgeo "DStratMC/board-geometry"
	"DStratMC/simulation"
	"context"
	_ "embed"
	"fmt"
	g "github.com/AllenDang/giu"
	"image"
	"math"
	"strconv"
	"time"
)

const (
	RadioExact        = iota // Just record exact hit where clicked
	RadioOneAvg              // Record one hit uniformly distributed within a circle
	RadioMultiAvg            // Record multiple hits uniformly distributed within a circle
	RadioOneNormal           // Record one hit normally distributed within a circle
	RadioMultiNormal         // Record multiple hits normally distributed within a circle
	RadioSearchNormal        // Search around the board, recording result of multi-normal at each search location
)

const LeftToolbarMinimumWidth = 200
const singleHitMarkerRadius = 5
const multipleHitMarkerRadius = 1

const throwsAtOneTarget = 5_000
const numThrowsTextWidth = 120

const uniformCEPRadius = 0.3

//const normalCEPRadius = 0.3

const normalCEPRadius = 0.25

// const normalCEPRadius = 0.1
const stubStandardDeviation = normalCEPRadius * 2

const testCoordinateConversion = true

type UserInterface interface {
	MainUiLoop()
}

type UserInterfaceInstance struct {
	radioValue                 int
	dartboardTexture           *g.Texture
	dartboard                  Dartboard
	accuracyModel              simulation.AccuracyModel
	scoreDisplay               string
	messageDisplay             string
	searchResultStrings        [10]string
	throwTotal                 int64
	throwCount                 int64
	throwAverage               float64
	numThrowsField             int32
	drawReferenceLinesCheckbox bool

	drawOneSigma   bool
	drawTwoSigma   bool
	drawThreeSigma bool

	searchComplete    bool
	searchingBlinkOn  bool
	cancelBlinkTimer  context.CancelFunc
	simResultsOneEach []simulation.OneResult
}

func NewUserInterface(loadedImage *image.RGBA) UserInterface {
	instance := &UserInterfaceInstance{
		radioValue:                 RadioOneNormal,
		drawOneSigma:               false,
		drawTwoSigma:               false,
		drawThreeSigma:             false,
		searchResultStrings:        [10]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
		dartboard:                  NewDartboard(),
		drawReferenceLinesCheckbox: true,
		numThrowsField:             throwsAtOneTarget,
	}
	g.EnqueueNewTextureFromRgba(loadedImage, func(t *g.Texture) {
		instance.dartboardTexture = t
	})
	//fmt.Println("NewTargetSupplier returns", instance)
	instance.dartboard.SetDrawRefLines(instance.drawReferenceLinesCheckbox)
	instance.dartboard.SetClickCallback(instance.dartboardClickCallback)
	return instance
}

func (u *UserInterfaceInstance) MainUiLoop() {
	window := u.setUpWindow()

	window.Layout(
		u.leftToolbarLayout(),
		g.Custom(u.dartboard.DrawFunction),
	)

}

func (u *UserInterfaceInstance) leftToolbarLayout() g.Widget {
	u.accuracyModel = u.getAccuracyModel(u.radioValue)
	return g.Layout{
		//	Checkbox controlling whether crosshairs are drawn
		g.Checkbox("Reference Lines", &u.drawReferenceLinesCheckbox).OnChange(func() { u.dartboard.SetDrawRefLines(u.drawReferenceLinesCheckbox) }),

		// Fields used to select type of interaction and display messages
		u.uiLayoutInteractionTypeRadioButtons(),
		u.uiLayoutResetButton(),
		u.uiLayoutOptionalResultsMessage(),

		//	The following fields may be presented depending on the type of interaction
		u.uiLayoutNumberOfThrowsField(),
		u.uiLayoutStdCircleCheckboxes(),
		u.uiLayoutSearchButton(),
		u.uiLayoutBlinkingSearchNotice(),
		u.uiLayoutSearchResults(),
		u.uiLayoutAverageScore(),
	}

}

// uiLayoutInteractionTypeRadioButtons lays out the radio buttons to select the type of interaction and model
func (u *UserInterfaceInstance) uiLayoutInteractionTypeRadioButtons() g.Widget {
	return g.Layout{
		g.Label(""),
		g.RadioButton("One Exact", u.radioValue == RadioExact).OnChange(func() {
			u.radioValue = RadioExact
			u.accuracyModel = u.getAccuracyModel(u.radioValue)
			u.radioChanged()
		}),
		//g.RadioButton("One Throw Uniform", u.radioValue == RadioOneAvg).OnChange(func() {
		//	u.radioValue = RadioOneAvg
		//	u.accuracyModel = u.getAccuracyModel(u.radioValue)
		//	u.radioChanged()
		//}),
		//g.RadioButton("Multi Throw Uniform", u.radioValue == RadioMultiAvg).OnChange(func() {
		//	u.radioValue = RadioMultiAvg
		//	u.accuracyModel = u.getAccuracyModel(u.radioValue)
		//	u.radioChanged()
		//}),
		g.RadioButton("One Throw Normal", u.radioValue == RadioOneNormal).OnChange(func() {
			u.radioValue = RadioOneNormal
			u.accuracyModel = u.getAccuracyModel(u.radioValue)
			u.radioChanged()
		}),
		g.RadioButton("Multi Throw Normal", u.radioValue == RadioMultiNormal).OnChange(func() {
			u.radioValue = RadioMultiNormal
			u.accuracyModel = u.getAccuracyModel(u.radioValue)
			u.radioChanged()
		}),
		g.RadioButton("Search Normal", u.radioValue == RadioSearchNormal).OnChange(func() {
			u.radioValue = RadioSearchNormal
			u.accuracyModel = u.getAccuracyModel(u.radioValue)
			u.radioChanged()
		}),
	}
}

// uiLayoutResetButton lays out the Reset button in the left toolbar
func (u *UserInterfaceInstance) uiLayoutResetButton() g.Widget {
	return g.Layout{
		g.Label(""),
		g.Button("Reset").OnClick(u.radioChanged),
	}
}

// uiLayoutOptionalResultsMessage displays a generic message and throw score, if either is nonblanks
func (u *UserInterfaceInstance) uiLayoutOptionalResultsMessage() g.Widget {
	return g.Layout{
		g.Condition(u.messageDisplay != "" || u.scoreDisplay != "",
			g.Layout{
				g.Label(""),
				g.Label(u.messageDisplay),
				g.Label(u.scoreDisplay),
			}, nil),
	}
}

// uiLayoutStdCircleCheckboxes will, when we are doing normal distribution (and only then) offer 3 checkboxes for drawing reference
// circles at 1, 2, and 3 standard deviations
func (u *UserInterfaceInstance) uiLayoutStdCircleCheckboxes() g.Widget {
	return g.Layout{
		g.Condition(u.radioValue == RadioOneNormal || u.radioValue == RadioMultiNormal || u.radioValue == RadioSearchNormal,
			g.Layout{
				g.Label(""),
				g.Checkbox("1 Sigma", &u.drawOneSigma).OnChange(func() { u.dartboard.SetDrawOneSigma(u.drawOneSigma, u.accuracyModel.GetSigmaRadius(1)) }),
				g.Checkbox("2 Sigma", &u.drawTwoSigma).OnChange(func() { u.dartboard.SetDrawTwoSigma(u.drawTwoSigma, u.accuracyModel.GetSigmaRadius(2)) }),
				g.Checkbox("3 Sigma", &u.drawThreeSigma).OnChange(func() { u.dartboard.SetDrawThreeSigma(u.drawThreeSigma, u.accuracyModel.GetSigmaRadius(3)) }),
			}, nil),
	}
}

// uiLayoutNumberOfThrowsField displays a field to enter an integer number of throws
func (u *UserInterfaceInstance) uiLayoutNumberOfThrowsField() g.Widget {
	return g.Layout{
		// If we are doing multiple throws, allow the user to set the number of throws
		g.Condition(u.radioValue == RadioMultiAvg || u.radioValue == RadioMultiNormal || u.radioValue == RadioSearchNormal,
			g.Layout{
				g.Label(""),
				g.InputInt(&u.numThrowsField).Label("# Throws").
					Size(numThrowsTextWidth).
					StepSize(1).
					StepSizeFast(100),
			}, nil),
	}
}

// uiLayoutSearchButton will, If we are doing a search, offer a "SEARCH" button to begin
func (u *UserInterfaceInstance) uiLayoutSearchButton() g.Widget {
	return g.Layout{
		g.Condition(u.radioValue == RadioSearchNormal,
			g.Layout{
				g.Label(""),
				g.Button("SEARCH").OnClick(func() {
					u.startSearchForBestThrow(u.accuracyModel, u.numThrowsField)
				}),
			}, nil),
	}
}

// uiLayoutBlinkingSearchNotice displays a "searching please wait" message that blinks on and
// off (blinking caused by displaying the message dependent on a flag being toggled by a background process)
func (u *UserInterfaceInstance) uiLayoutBlinkingSearchNotice() g.Widget {
	return g.Layout{
		g.Condition(u.radioValue == RadioSearchNormal,
			g.Layout{
				g.Label(""),
				g.Condition(u.searchingBlinkOn,
					g.CSSTag("waitlabel").To(
						g.Label("Searching, please wait"),
					),
					g.Label("")),
			}, nil),
	}
}

// uiLayoutAverageScore displays the average score from non-search clicks
func (u *UserInterfaceInstance) uiLayoutAverageScore() g.Widget {
	return g.Layout{
		g.Condition(u.throwCount > 0,
			g.Layout{
				g.Label(""),
				g.Label("Throws: " + strconv.Itoa(int(u.throwCount))),
				g.Label("Total: " + strconv.Itoa(int(u.throwTotal))),
				g.Label("Average: " + strconv.FormatFloat(u.throwAverage, 'f', 1, 64)),
			},
			nil),
	}
}

// uiLayoutSearchResults lays out the fields that report search results
func (u *UserInterfaceInstance) uiLayoutSearchResults() g.Widget {
	return g.Layout{
		g.Condition(u.radioValue == RadioSearchNormal && u.searchComplete,
			g.Layout{
				g.Label("Best 10 throws:"),
				g.Label(""),
				u.uiLayoutSearchResultLabels(10),
			}, nil)}
}

// uiLayoutSearchResultLabels lays out a number of label fields that will be used to display search results
func (u *UserInterfaceInstance) uiLayoutSearchResultLabels(numLabels int) g.Layout {
	widgetList := make([]g.Widget, 0, numLabels)
	for i := 0; i < numLabels; i++ {
		//thisItem := g.Label(u.searchResultStrings[i])
		thisItem := g.Button(u.searchResultStrings[i]).
			OnClick(func() {
				u.resultButtonClicked(i)
			})
		widgetList = append(widgetList, thisItem)
	}
	return widgetList
}

// resultButtonClicked is called when a button is clicked that corresponds to a search result
// It will draw a marker at the precise location of that reported search result
func (u *UserInterfaceInstance) resultButtonClicked(buttonIndex int) {
	if buttonIndex < len(u.simResultsOneEach) {
		//	Set the position where any requested standard deviation circles will be drawn
		targetPosition := u.simResultsOneEach[buttonIndex].Position
		u.dartboard.SetStdDeviationCirclesCentre(targetPosition)
		u.dartboard.QueueTargetMarker(targetPosition)
		g.Update()
	}
}

func (u *UserInterfaceInstance) startSearchForBestThrow(model simulation.AccuracyModel, numThrows int32) {
	u.searchResultStrings = [10]string{"", "", "", "", "", "", "", "", "", ""}
	u.dartboard.RemoveThrowMarkers()
	u.searchComplete = false
	g.Update()
	//fmt.Printf("Starting search for the best throw. model=%#v, numThrows=%d\n", model, numThrows)
	//	Start a process to blink the "searching" label on and off
	var ctx context.Context
	ctx, u.cancelBlinkTimer = context.WithCancel(context.Background())
	go u.cycleBlinkFlag(ctx)

	//	Start the actual search process
	go u.searchProcess(model, numThrows)
	//fmt.Println("Search started")
}

func (u *UserInterfaceInstance) cycleBlinkFlag(ctx context.Context) {
	//fmt.Println("Starting blink timer routine")
	for {
		select {
		case <-ctx.Done():
			//fmt.Println("Blink timer stopped")
			u.searchingBlinkOn = false
			return
		default:
			//fmt.Println("Blink timer running")
			u.searchingBlinkOn = !u.searchingBlinkOn
			g.Update()
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func (u *UserInterfaceInstance) searchProcess(model simulation.AccuracyModel, numThrows int32) {
	//fmt.Println("search process entered")
	//	Get target iterator and results aggregator
	targetSupplier := simulation.NewTargetSupplier(u.dartboard.GetSquareDimension(), u.dartboard.GetImageMinPoint())
	results := simulation.NewSimResults()
	// Loop through all targets
	for targetSupplier.HasNext() {
		target := targetSupplier.NextTarget()
		// Do throws at this target
		averageScore, err := u.multipleThrowsAtTarget(target, model, numThrows)
		if err != nil {
			fmt.Printf("Error getting throw %v", err)
			continue
		}
		//	record result for this target
		results.AddTargetResult(target, averageScore)
	}
	//	Get results, sorted from best to worst
	sortedResults := results.GetResultsSortedByHighScore()
	//fmt.Println("First few sorted results:", sortedResults[:5])
	//  Filter results so each plain-language target is named only once
	u.simResultsOneEach = simulation.FilterToOneTargetEach(sortedResults)
	//fmt.Println("First few one each results:", oneEach[:5])
	// Message saying what was the best target
	fmt.Println("First 5 best choices, from best down:")
	for i := 0; i < 10; i++ {
		_, score, description := boardgeo.DescribeBoardPoint(u.simResultsOneEach[i].Position)
		fmt.Printf("   %s (theoretical score %d, average %g)\n", description, score, u.simResultsOneEach[i].Score)
		u.searchResultStrings[i] = fmt.Sprintf("%s (%.2f)", description, u.simResultsOneEach[i].Score)
		u.searchComplete = true
	}
	//	Draw best target on the board
	bestTargetPosition := u.simResultsOneEach[0].Position
	u.dartboard.SetStdDeviationCirclesCentre(bestTargetPosition)
	u.dartboard.QueueTargetMarker(bestTargetPosition)
	//	Stop the blink timer
	u.cancelBlinkTimer()
	g.Update()
	//fmt.Println("Search process ends")
}

func (u *UserInterfaceInstance) multipleThrowsAtTarget(target boardgeo.BoardPosition, model simulation.AccuracyModel, throws int32) (float64, error) {
	//fmt.Println("multipleThrowsAtTarget", target, model, throws)
	var total float64 = 0.0
	for i := 0; i < int(throws); i++ {
		hit, err := model.GetThrow(target,
			u.dartboard.GetScoringRadiusPixels(),
			u.dartboard.GetSquareDimension(),
			u.dartboard.GetImageMinPoint())
		if err != nil {
			return 0.0, err
		}
		_, score, _ := boardgeo.DescribeBoardPoint(hit)
		total += float64(score)
	}
	average := total / float64(throws)
	//fmt.Println("   average", average)
	return average, nil
}

func (u *UserInterfaceInstance) getAccuracyModel(radioValue int) simulation.AccuracyModel {
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
	case RadioSearchNormal:
		return simulation.NewNormalAccuracyModel(normalCEPRadius, stubStandardDeviation)
	default:
		panic("Invalid radio button value")
		return simulation.NewPerfectAccuracyModel()
	}
}

func (u *UserInterfaceInstance) setUpWindow() *g.WindowWidget {
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

	u.dartboard.SetInfo(window, u.dartboardTexture, squareDimension, dartboardImageMin, dartboardImageMax)
	return window
}

func (u *UserInterfaceInstance) radioChanged() {
	u.scoreDisplay = ""
	u.messageDisplay = ""
	u.throwTotal = 0
	u.throwCount = 0
	u.throwAverage = 0
	u.searchResultStrings = [10]string{"", "", "", "", "", "", "", "", "", ""}
	u.dartboard.RemoveThrowMarkers()
	u.searchComplete = false
	u.searchingBlinkOn = false
}

func (u *UserInterfaceInstance) dartboardClickCallback(dartboard Dartboard, position boardgeo.BoardPosition) {
	//fmt.Printf("Dartboard clicked at position %#v\n", position)
	// This is a good place to verify that coordinate conversion is working
	if testCoordinateConversion {
		testConvertPolar := boardgeo.CreateBoardPositionFromPolar(position.Radius, position.Angle, dartboard.GetSquareDimension())
		if position.Radius != testConvertPolar.Radius || position.Angle != testConvertPolar.Angle {
			panic("Coordinate conversion failed: polar coordinates do not match")
		}
		xDelta := math.Abs(float64(position.XMouseInside) - float64(testConvertPolar.XMouseInside))
		yDelta := math.Abs(float64(position.YMouseInside) - float64(testConvertPolar.YMouseInside))
		if xDelta > 1 || yDelta > 1 {
			details := fmt.Sprintf("x %d,%d  y %d,%d",
				position.XMouseInside, testConvertPolar.XMouseInside,
				position.YMouseInside, testConvertPolar.YMouseInside)
			panic("Coordinate conversion failed: cartesian coordinates do not match: " + details)
		}
	}
	//fmt.Printf("  Polar converted back to %#v\n", testConvertPolar)
	if position.Radius <= 1.0 {
		u.messageDisplay = ""
		u.scoreDisplay = ""
		dartboard.RemoveThrowMarkers()
		switch u.radioValue {
		case RadioExact:
			dartboard.QueueTargetMarker(position)
			_, score, description := boardgeo.DescribeBoardPoint(position)
			u.messageDisplay = description
			u.scoreDisplay = strconv.Itoa(score) + " points"
		case RadioOneAvg:
			u.oneUniformThrow(dartboard, position, u.accuracyModel)
		case RadioMultiAvg:
			u.multipleUniformThrows(dartboard, position, u.accuracyModel)
		case RadioOneNormal:
			u.oneNormalThrow(dartboard, position, u.accuracyModel)
		case RadioMultiNormal:
			u.multipleNormalThrows(dartboard, position, u.accuracyModel)
		case RadioSearchNormal:
			u.messageDisplay = "Click SEARCH to begin"
		default:
			panic("Invalid radio button value")
		}
	}
}

func (u *UserInterfaceInstance) oneUniformThrow(dartboard Dartboard, position boardgeo.BoardPosition, model simulation.AccuracyModel) {
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
	u.messageDisplay = description
	u.scoreDisplay = strconv.Itoa(score) + " points"
	u.throwCount++
	u.throwTotal += int64(score)
	u.throwAverage = float64(u.throwTotal) / float64(u.throwCount)
	g.Update()

}

func (u *UserInterfaceInstance) multipleUniformThrows(dartboard Dartboard, position boardgeo.BoardPosition, model simulation.AccuracyModel) {
	//fmt.Printf("multipleUniformThrows  %v,\n", position)

	//  Draw a marker to record where we clicked
	dartboard.QueueTargetMarker(position)

	//	Draw a circle showing the accuracy radius around the clicked point
	accuracyRadius := model.GetAccuracyRadius()
	dartboard.QueueAccuracyCircle(position, accuracyRadius)
	dartboard.AllocateHitsSpace(int(u.numThrowsField))

	u.throwCount = 0
	u.throwTotal = 0
	for i := 0; i < int(u.numThrowsField); i++ {
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
		u.throwCount++
		u.throwTotal += int64(score)
		u.throwAverage = float64(u.throwTotal) / float64(u.throwCount)
	}
	g.Update()

}

func (u *UserInterfaceInstance) oneNormalThrow(dartboard Dartboard, position boardgeo.BoardPosition, model simulation.AccuracyModel) {
	//fmt.Printf("oneNormalThrow  %v,\n", position)

	//  Draw a marker to record where we clicked
	dartboard.QueueTargetMarker(position)

	//	Set the position where any requested standard deviation circles will be drawn
	dartboard.SetStdDeviationCirclesCentre(position)

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
	u.messageDisplay = description
	u.scoreDisplay = strconv.Itoa(score) + " points"
	u.throwCount++
	u.throwTotal += int64(score)
	u.throwAverage = float64(u.throwTotal) / float64(u.throwCount)
	g.Update()

}

func (u *UserInterfaceInstance) multipleNormalThrows(dartboard Dartboard, position boardgeo.BoardPosition, model simulation.AccuracyModel) {
	//fmt.Printf("oneNormalThrow  %v,\n", position)

	//  Draw a marker to record where we clicked
	dartboard.QueueTargetMarker(position)

	//	Set the position where any requested standard deviation circles will be drawn
	dartboard.SetStdDeviationCirclesCentre(position)

	dartboard.AllocateHitsSpace(int(u.numThrowsField))

	u.throwCount = 0
	u.throwTotal = 0

	for i := 0; i < int(u.numThrowsField); i++ {
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
		u.throwCount++
		u.throwTotal += int64(score)
		u.throwAverage = float64(u.throwTotal) / float64(u.throwCount)
	}
	g.Update()

}
