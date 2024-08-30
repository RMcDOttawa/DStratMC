package ui

import (
	boardgeo "DStratMC/board-geometry"
	"DStratMC/simulation"
	target_search "DStratMC/target-search"
	"context"
	_ "embed"
	"fmt"
	g "github.com/AllenDang/giu"
	"image"
	"image/color"
	"math"
	"strconv"
)

// UserInterface models the overall UI as an object. Even thought there won't be multiple instances, modeling
// it this way allows control over the associated data, reducing the need for global variables and the errors they invite
type UserInterface interface {
	MainUiLoop()
}

// Drawing the std-dev circle manually involves several modes the UI can be in

type drawCircleState int

const (
	drawCircleStateOff drawCircleState = iota
	drawCircleStateStart
	drawCircleStateDrawing
)

// Calculating the std-dev by measuring real throws involves several modes the UI can be in

type measureStdDevState int

const (
	measureStdDevStateOff measureStdDevState = iota
	measureStdDevStateSelectTarget
	measureStdDevStateThrowing
)

// UserInterfaceInstance is the attribute data stored with the UI object
type UserInterfaceInstance struct {
	dartboardTexture *g.Texture
	dartboard        Dartboard
	accuracyModel    simulation.AccuracyModel
	mode             InterfaceMode

	scoreDisplay   string
	messageDisplay string
	throwTotal     int64
	throwCount     int64
	throwAverage   float64
	numThrowsField int32

	drawReferenceLinesCheckbox bool
	drawOneSigma               bool
	drawTwoSigma               bool
	drawThreeSigma             bool

	searchShowEachTarget  bool
	searchProgressPercent float64
	searchComplete        bool
	searchResultStrings   [10]string
	searchResultsRadio    int
	searchingBlinkOn      bool
	cancelSearchVisible   bool
	cancelBlinkTimer      context.CancelFunc
	cancelSearch          context.CancelFunc
	searchCancelled       bool
	simResultsOneEach     []target_search.OneResult
	stdDevInputField      float32

	// Drawing circle to represent standard deviation
	circleDrawingState drawCircleState
	dartboardImageMin  image.Point
	dartboardImageMax  image.Point

	//	Data structure collecting real throws to calculate standard deviation
	realThrows       simulation.RealThrowCollection
	measurementState measureStdDevState
	measuringTarget  boardgeo.BoardPosition
}

var panelBorderColour = color.RGBA{100, 100, 100, 255}

// NewUserInterface creates a new UserInterface object
func NewUserInterface(loadedImage *image.RGBA) UserInterface {
	instance := &UserInterfaceInstance{
		mode:                       Mode_OneNormal,
		messageDisplay:             "",
		scoreDisplay:               "",
		drawOneSigma:               false,
		drawTwoSigma:               false,
		drawThreeSigma:             false,
		searchShowEachTarget:       false,
		searchResultStrings:        [10]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
		dartboard:                  NewDartboard(),
		drawReferenceLinesCheckbox: true,
		numThrowsField:             throwsAtOneTarget,
		stdDevInputField:           0.15,
		circleDrawingState:         drawCircleStateOff,
		realThrows:                 simulation.NewRealThrowCollectionInstance(),
	}
	g.EnqueueNewTextureFromRgba(loadedImage, func(t *g.Texture) {
		instance.dartboardTexture = t
	})

	instance.dartboard.SetDrawRefLines(instance.drawReferenceLinesCheckbox)
	instance.dartboard.SetClickCallback(instance.dartboardClickCallback)
	return instance
}

// MainUiLoop is the mac-binary loop for the user interface. This called from the master window's Run method,
// 30 times a second repeatedly.  The GIU ui framework does not store state, so we are responsible for
// storing state information for anything that is to be displayed that seems to be constant on the screen

// The UI is divided into two sections: a vertical toolbar on the left, and a square dartboard area on the right.
func (u *UserInterfaceInstance) MainUiLoop() {
	window := u.setUpWindow()

	window.Layout(
		u.leftToolbarLayout(),
		g.Custom(u.dartboard.DrawFunction),
	)

	//	Click Callbacks are processed only when the mouse is released.  Because, for purposes of
	//  tracing the standard deviation circle, we want to handle the cases of the mouse first
	//  being clicked, and dragging while down, we will test for those cases here after the
	//  usual UI processing.
	if u.circleDrawingState != drawCircleStateOff {
		// We are in some kind of circle drawing state
		u.handleDrawingCircle()
		return
	}

}

//		If we have started "draw circle" mode, we will trace a circle on the dartboard as long
//	 as the mouse button is down
func (u *UserInterfaceInstance) handleDrawingCircle() {
	mousePosition := boardgeo.CreateBoardPositionFromXY(g.GetMousePos(), u.dartboard.GetSquareDimension(),
		u.dartboardImageMin)

	//fmt.Println("Handling circle drawing, state: ", u.circleDrawingState)
	//fmt.Println("  Mouse position: ", mousePosition)
	//  We are in "draw circle" mode.
	//  There are 3 states of interest
	//  1. Start mode and mouse is down:  We record the center of the circle
	//		for drawing the next time through the loop after the mouse has moved a bit.
	//		Go to "Drawing" mode to start drawing on the next iteration.
	//  2. Drawing mode and mouse is still down: keep drawing, stay in drawing mode
	//  3. Drawing mode and mouse is released:  stop drawing, store result, set mode = off
	//		We handle this case in the real click callback routine

	//  1. Start mode and mouse is down:  We record the center point for the circle and set mode = drawing
	if u.circleDrawingState == drawCircleStateStart && g.IsMouseDown(g.MouseButtonLeft) {
		//fmt.Println("   First event for circle drawing. Set center to", mousePosition)
		u.circleDrawingState = drawCircleStateDrawing
		u.dartboard.StartTracingCircleAtCenter(mousePosition)
		return
	}

	//  2. Drawing mode and mouse is still down: keep drawing, stay in drawing mode
	if u.circleDrawingState == drawCircleStateDrawing && g.IsMouseDown(g.MouseButtonLeft) {
		if mousePosition.Radius > 1.0 {
			fmt.Println("Mouse position outside dartboard, ignoring")
			return
		}
		circleRadiusPixels := boardgeo.PixelDistanceBetweenBoardPositions(
			u.dartboard.GetTracingCircleCenter(),
			mousePosition, u.dartboard.GetSquareDimension())
		u.dartboard.SetTracingCircleRadius(circleRadiusPixels)
		//fmt.Printf("  Continue drawing circle center %v, mouse %v, radius %d\n",
		//	u.dartboard.GetTracingCircleCenter(), mousePosition, circleRadiusPixels)
		return
	}

	//  3. Drawing mode and mouse is released:  handled during click callback
	//if u.circleDrawingState == drawCircleStateDrawing && !g.IsMouseDown(g.MouseButtonLeft) {
	//	circleRadiusPixels := boardgeo.PixelDistanceBetweenBoardPositions(u.dartboard.GetTracingCircleCenter(), mousePosition)
	//	fmt.Printf("  Stop drawing circle, record result center %v, radius %d\n",
	//		u.dartboard.GetTracingCircleCenter(), circleRadiusPixels)
	//	u.circleDrawingState = drawCircleStateOff
	//	return
	//}

	// Getting her means we are in "start" mode but the user hasn't clicked the mouse yet.
	//	Verify that
	if u.circleDrawingState == drawCircleStateStart && !g.IsMouseDown(g.MouseButtonLeft) {
		//fmt.Println("  Mouse not clicked yet")
		return
	}

	fmt.Println("Invalid state for circle drawing ignored. State: ", u.circleDrawingState)

}

// setupWindow sets up the window for the user interface
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
	u.dartboardImageMin = image.Pt(int(windowX)+leftToolbarWidth, int(windowY))
	u.dartboardImageMax = image.Pt(u.dartboardImageMin.X+int(squareDimension), u.dartboardImageMin.Y+int(squareDimension))
	//fmt.Printf("image min %d, max %d\n", imageMin, imageMax)

	u.dartboard.SetInfo(window, u.dartboardTexture, u.dartboardImageMin, u.dartboardImageMax, leftToolbarWidth)
	return window
}

// leftToolbarLayout lays out the left toolbar, which contains buttons and fields for user interaction
//
//	Some of the elements in the toolbar are only displayed when certain radio buttons are selected
func (u *UserInterfaceInstance) leftToolbarLayout() g.Widget {
	u.accuracyModel = u.getAccuracyModel(u.mode)
	return g.Layout{

		u.uiLayoutInteractionModePanel(),
		u.uiLayoutMessagesPanel(),
		u.uiLayoutNumberOfThrowsPanel(),
		u.uiLayoutNormalInfoPanel(),
		u.uiSearchControlsPanel(),
		u.uiRealThrowMeasurementControls(),

		u.uiLayoutSearchResults(),
		u.uiLayoutAverageScore(),
	}
}

// uiLayoutInteractionModePanel lays out the controls for the overall interaction mode of the UI:
// radio buttons to select the interaction mode, a checkbox to draw reference lines, and a reset button
func (u *UserInterfaceInstance) uiLayoutInteractionModePanel() g.Widget {
	fieldsLayout := g.Layout{
		g.Label("Interaction Mode:"),
		g.RadioButton("One Exact", u.mode == Mode_Exact).OnChange(func() {
			u.mode = Mode_Exact
			u.accuracyModel = u.getAccuracyModel(u.mode)
			u.radioChanged()
		}),
		g.RadioButton("Draw 95% Circle", u.mode == Mode_DrawCircle).OnChange(func() {
			u.mode = Mode_DrawCircle
			u.accuracyModel = u.getAccuracyModel(u.mode)
			u.radioChanged()
			u.messageDisplay = "Draw 95% Circle"
			// Record circle drawing state last because radioChanged resets it
			u.circleDrawingState = drawCircleStateStart
		}),
		g.RadioButton("Measure Real Throws", u.mode == Mode_EmpricalStdDev).OnChange(func() {
			u.mode = Mode_EmpricalStdDev
			u.measurementState = measureStdDevStateSelectTarget
			u.accuracyModel = u.getAccuracyModel(u.mode)
			u.radioChanged()
			u.messageDisplay = "Select Target"
		}),
		// The following two radio buttons were used in early development stages, but are deprecated
		// The code to implement them is still present, so you can un-comment them to resume their function
		//g.RadioButton("One Throw Uniform", u.mode == Mode_OneAvg).OnChange(func() {
		//	u.mode = Mode_OneAvg
		//	u.accuracyModel = u.getAccuracyModel(u.mode)
		//	u.radioChanged()
		//}),
		//g.RadioButton("Multi Throw Uniform", u.mode == Mode_MultiAvg).OnChange(func() {
		//	u.mode = Mode_MultiAvg
		//	u.accuracyModel = u.getAccuracyModel(u.mode)
		//	u.radioChanged()
		//}),
		g.RadioButton("One Throw Normal", u.mode == Mode_OneNormal).OnChange(func() {
			u.mode = Mode_OneNormal
			u.accuracyModel = u.getAccuracyModel(u.mode)
			u.radioChanged()
		}),
		g.RadioButton("Multi Throw Normal", u.mode == Mode_MultiNormal).OnChange(func() {
			u.mode = Mode_MultiNormal
			u.accuracyModel = u.getAccuracyModel(u.mode)
			u.radioChanged()
		}),
		g.RadioButton("Search Normal", u.mode == Mode_SearchNormal).OnChange(func() {
			u.mode = Mode_SearchNormal
			u.accuracyModel = u.getAccuracyModel(u.mode)
			u.radioChanged()
		}),
		g.Dummy(0, BlankLineHeight),
		g.Checkbox("Reference Lines", &u.drawReferenceLinesCheckbox).OnChange(func() { u.dartboard.SetDrawRefLines(u.drawReferenceLinesCheckbox) }),
		g.Dummy(0, BlankLineHeight),
		g.Button("Reset").OnClick(u.radioChanged),
	}
	const numRadioButtons = 6
	const numButtons = 1
	const numLabels = 3
	const numCheckboxes = 1
	return g.Style().
		// Fields inside a bordered panel
		SetColor(g.StyleColorBorder, panelBorderColour).
		To(
			g.Child().Border(true).
				Size(LeftToolbarChildWidth,
					numRadioButtons*uiRadioButtonHeight+
						numButtons*uiButtonHeight+
						numCheckboxes*uiCheckboxHeight+
						numLabels*uiLabelHeight+
						4).
				Layout(fieldsLayout),
		)
}

// uiLayoutOptionalResultsMessage displays a generic message and throw score
func (u *UserInterfaceInstance) uiLayoutMessagesPanel() g.Widget {
	fieldsLayout := g.Layout{
		g.Label(u.messageDisplay),
		g.Label(u.scoreDisplay),
	}
	const numLabelsInsideChild = 2
	return g.Condition(!(u.mode == Mode_SearchNormal || u.mode == Mode_MultiNormal),
		g.Layout{
			// Blank line before
			// Fields inside a bordered panel
			g.Style().
				SetColor(g.StyleColorBorder, panelBorderColour).
				To(
					g.Child().Border(true).
						Size(LeftToolbarChildWidth,
							numLabelsInsideChild*uiLabelHeight).
						Layout(fieldsLayout),
				),
		}, nil)
}

func (u *UserInterfaceInstance) uiLayoutNormalInfoPanel() g.Widget {
	fieldsLayout := g.Layout{
		g.Label("Normal Distribution"),
		g.Dummy(0, BlankLineHeight),
		g.InputFloat(&u.stdDevInputField).
			Label("StdDev 0-1").
			Size(stdDevTextWidth).
			OnChange(u.validateAndProcessStdDevField),
		g.Dummy(0, BlankLineHeight),
		g.Label("Show circles for:"),
		g.Checkbox("1 Sigma", &u.drawOneSigma).OnChange(func() { u.dartboard.SetDrawOneSigma(u.drawOneSigma, u.accuracyModel.GetSigmaRadius(1)) }),
		g.Checkbox("2 Sigma", &u.drawTwoSigma).OnChange(func() { u.dartboard.SetDrawTwoSigma(u.drawTwoSigma, u.accuracyModel.GetSigmaRadius(2)) }),
		g.Checkbox("3 Sigma", &u.drawThreeSigma).OnChange(func() { u.dartboard.SetDrawThreeSigma(u.drawThreeSigma, u.accuracyModel.GetSigmaRadius(3)) }),
	}
	const numLabels = 4
	const numCheckboxes = 3
	return g.Condition(u.mode != Mode_Exact && u.mode != Mode_EmpricalStdDev,
		g.Layout{
			g.Style().
				// Fields inside a bordered panel
				SetColor(g.StyleColorBorder, panelBorderColour).
				SetDisabled(!(u.mode == Mode_OneNormal || u.mode == Mode_MultiNormal || u.mode == Mode_SearchNormal)).
				To(
					g.Child().Border(true).
						Size(LeftToolbarChildWidth,
							numLabels*uiLabelHeight+
								numCheckboxes*uiCheckboxHeight+
								uiInputFieldHeight-20).
						Layout(fieldsLayout),
				),
		}, nil)
}

func (u *UserInterfaceInstance) uiSearchControlsPanel() g.Widget {
	fieldsLayout := g.Layout{
		g.Label("Search Controls"),
		g.Dummy(0, BlankLineHeight),
		g.Checkbox("Show Search", &u.searchShowEachTarget),
		g.Dummy(0, BlankLineHeight),
		g.Button("START SEARCH").OnClick(func() {
			u.startSearchForBestThrow(u.accuracyModel, u.numThrowsField)
		}),
		g.ProgressBar(float32(u.searchProgressPercent)).Size(LeftToolbarChildWidth-12, 0),
		g.Button("Cancel Search").OnClick(func() {
			fmt.Println("Cancelling Search")
			u.cancelSearch()
		}),
		g.Condition(u.searchingBlinkOn,
			g.CSSTag("waitlabel").To(
				g.Label("Searching, please wait"),
			),
			g.Dummy(0, BlankLineHeight)),
	}
	const numLabels = 4
	const numCheckboxes = 1
	const numButtons = 1
	return g.Condition(u.mode == Mode_SearchNormal,
		g.Layout{
			g.Style().
				// Fields inside a bordered panel
				SetColor(g.StyleColorBorder, panelBorderColour).
				SetDisabled(u.mode != Mode_SearchNormal).
				To(
					g.Child().Border(true).
						Size(LeftToolbarChildWidth,
							numLabels*uiLabelHeight+
								uiProgressBarHeight+
								numButtons*uiButtonHeight+
								numCheckboxes*uiCheckboxHeight+12).
						Layout(fieldsLayout),
				),
		}, nil)
}

func (u *UserInterfaceInstance) uiRealThrowMeasurementControls() g.Widget {
	fieldsLayout := g.Layout{
		g.Dummy(0, BlankLineHeight),
		g.Button("New Model").OnClick(func() {
			fmt.Println("New Model")
			u.realThrows = simulation.NewRealThrowCollectionInstance()
			u.measurementState = measureStdDevStateSelectTarget
			u.messageDisplay = "Click Target"
		}),
		g.Style().SetDisabled(u.measurementState != measureStdDevStateThrowing).To(
			g.Button("New Target").OnClick(func() {
				u.measurementState = measureStdDevStateSelectTarget
				u.messageDisplay = "Click Target"
			}),
		),
		g.Dummy(0, BlankLineHeight),
		g.Label(fmt.Sprintf("Data Points: %d", u.realThrows.GetNumThrows())),
		g.Label(fmt.Sprintf("Std Dev: %s", u.realThrows.GetStdDevString())),
		g.Dummy(0, BlankLineHeight),
		g.Button("Load").OnClick(func() { fmt.Println("Load STUB") }),
		g.Style().SetDisabled(u.realThrows.GetNumThrows() == 0).To(
			g.Button("Save").OnClick(func() { fmt.Println("Save STUB") }),
		),
		g.Dummy(0, BlankLineHeight),
		g.Style().SetDisabled(!u.realThrows.IsStdDevAvailable()).To(
			g.Button("Use StdDev").OnClick(func() {
				//
				//stdDev := u.realThrows.CalcStdDevOfThrows()
				//u.stdDevInputField = float32(stdDev)
				//u.accuracyModel.SetStandardDeviation(stdDev)
				//u.dartboard.SetDrawOneSigma(u.drawOneSigma, u.accuracyModel.GetSigmaRadius(1))
				//u.dartboard.SetDrawTwoSigma(u.drawTwoSigma, u.accuracyModel.GetSigmaRadius(2))
				//u.dartboard.SetDrawThreeSigma(u.drawThreeSigma, u.accuracyModel.GetSigmaRadius(3))
				//g.Update()
			}),
		),
	}
	return g.Condition(u.mode == Mode_EmpricalStdDev, fieldsLayout, nil)
}

func (u *UserInterfaceInstance) validateAndProcessStdDevField() {
	if u.stdDevInputField < .00001 {
		u.stdDevInputField = 0
		u.messageDisplay = "StdDev must be 0 to 1"
		return
	}
	if u.stdDevInputField > 1 {
		u.stdDevInputField = 1
		u.messageDisplay = "StdDev must be 0 to 1"
		return
	}
	u.messageDisplay = ""
	u.accuracyModel.SetStandardDeviation(float64(u.stdDevInputField))
	u.dartboard.SetDrawOneSigma(u.drawOneSigma, u.accuracyModel.GetSigmaRadius(1))
	u.dartboard.SetDrawTwoSigma(u.drawTwoSigma, u.accuracyModel.GetSigmaRadius(2))
	u.dartboard.SetDrawThreeSigma(u.drawThreeSigma, u.accuracyModel.GetSigmaRadius(3))
}

// uiLayoutNumberOfThrowsField displays a field to enter an integer number of throws
func (u *UserInterfaceInstance) uiLayoutNumberOfThrowsPanel() g.Widget {
	fieldsLayout := g.Layout{
		// If we are doing multiple throws, allow the user to set the number of throws
		//g.Condition(u.mode == Mode_MultiAvg || u.mode == Mode_MultiNormal || u.mode == Mode_SearchNormal,
		// Fields inside a bordered panel
		g.InputInt(&u.numThrowsField).Label("Throws").
			Size(numThrowsTextWidth).
			StepSize(100).
			StepSizeFast(1000).
			OnChange(u.validateNumThrowsField),
	}
	return g.Condition(u.mode == Mode_MultiNormal || u.mode == Mode_SearchNormal,
		g.Layout{
			g.Style().
				// Fields inside a bordered panel
				SetColor(g.StyleColorBorder, panelBorderColour).
				SetDisabled(!(u.mode == Mode_MultiNormal || u.mode == Mode_SearchNormal)).
				To(
					g.Child().Border(true).
						Size(LeftToolbarChildWidth,
							uiInputFieldHeight).
						Layout(fieldsLayout),
				),
		}, nil)
}

func (u *UserInterfaceInstance) validateNumThrowsField() {
	if u.numThrowsField < 1 {
		u.numThrowsField = throwsAtOneTarget
		u.messageDisplay = "numTrows must b > 0"
		return
	}
	u.messageDisplay = ""
}

// uiLayoutAverageScore displays the average score from non-search clicks
func (u *UserInterfaceInstance) uiLayoutAverageScore() g.Widget {
	return g.Layout{
		g.Condition(u.throwCount > 0,
			g.Layout{
				g.Dummy(0, BlankLineHeight),
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
		g.Condition(u.mode == Mode_SearchNormal && u.searchComplete,
			g.Layout{
				g.Label(fmt.Sprintf("Best %d targets:", numSearchResultsToDisplay)),
				u.uiLayoutSearchResultLabels(numSearchResultsToDisplay),
			}, nil)}
}

// uiLayoutSearchResultLabels lays out a number of radio buttons that will be used to display search results
func (u *UserInterfaceInstance) uiLayoutSearchResultLabels(numLabels int) g.Layout {
	widgetList := make([]g.Widget, 0, numLabels)
	for i := 0; i < numLabels; i++ {
		thisItem := g.RadioButton(u.searchResultStrings[i], u.searchResultsRadio == i).OnChange(func() {
			u.searchResultsRadio = i
			u.resultButtonClicked(i)
		})
		widgetList = append(widgetList, thisItem)
	}
	return widgetList
}

// getAccuracyModel returns the accuracy model that corresponds to the selected mode button
func (u *UserInterfaceInstance) getAccuracyModel(mode InterfaceMode) simulation.AccuracyModel {
	switch mode {
	case Mode_Exact:
		return nil
	case Mode_OneAvg:
		return simulation.NewUniformAccuracyModel(uniformCEPRadius)
	case Mode_MultiAvg:
		return simulation.NewUniformAccuracyModel(uniformCEPRadius)
	case Mode_OneNormal:
		return simulation.NewNormalAccuracyModel(float64(u.stdDevInputField))
	case Mode_MultiNormal:
		return simulation.NewNormalAccuracyModel(float64(u.stdDevInputField))
	case Mode_SearchNormal:
		return simulation.NewNormalAccuracyModel(float64(u.stdDevInputField))
	case Mode_DrawCircle:
		// Doesn't matter what model we return, as it isn't used in this mode
		return simulation.NewNormalAccuracyModel(float64(u.stdDevInputField))
	case Mode_EmpricalStdDev:
		// Doesn't matter what model we return, as it isn't used in this mode
		return simulation.NewNormalAccuracyModel(float64(u.stdDevInputField))
	default:
		panic("Invalid radio button value")
		return simulation.NewPerfectAccuracyModel()
	}
}

// radioChanged responds to a change to the mode radio button by resetting various display fields and counters
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
	u.circleDrawingState = drawCircleStateOff
}

// dartboardClickCallback is called when the user clicks on the dartboard. It is the mac-binary entry point for
// the UI to respond to user input
func (u *UserInterfaceInstance) dartboardClickCallback(dartboard Dartboard, position boardgeo.BoardPosition) {
	// This is a good place to verify that coordinate conversion is working
	if testCoordinateConversion {
		testConvertPolar := boardgeo.CreateBoardPositionFromPolar(position.Radius, position.Angle)
		if position.Radius != testConvertPolar.Radius || position.Angle != testConvertPolar.Angle {
			panic("Coordinate conversion failed: polar coordinates do not match")
		}
		posX, posY := boardgeo.GetXY(position, u.dartboard.GetSquareDimension())
		convertX, convertY := boardgeo.GetXY(testConvertPolar, u.dartboard.GetSquareDimension())
		xDelta := math.Abs(float64(posX) - float64(convertX))
		yDelta := math.Abs(float64(posY) - float64(convertY))
		if xDelta > 1 || yDelta > 1 {
			details := fmt.Sprintf("x %d,%d  y %d,%d",
				posX, convertX,
				posY, convertY)
			panic("Coordinate conversion failed: cartesian coordinates do not match: " + details)
		}
	}

	if position.Radius <= 1.0 || u.mode == Mode_DrawCircle {
		u.messageDisplay = ""
		u.scoreDisplay = ""
		dartboard.RemoveThrowMarkers()

		//	If it wasn't "start drawing circle" mode, we just take this as a throw marker
		switch u.mode {
		case Mode_Exact:
			dartboard.QueueTargetMarker(position)
			_, score, description := boardgeo.DescribeBoardPoint(position)
			u.messageDisplay = description
			u.scoreDisplay = strconv.Itoa(score) + " points"
		case Mode_OneAvg:
			u.oneUniformThrow(dartboard, position, u.accuracyModel)
		case Mode_MultiAvg:
			u.multipleUniformThrows(dartboard, position, u.accuracyModel)
		case Mode_OneNormal:
			u.oneNormalThrow(dartboard, position, u.accuracyModel)
		case Mode_MultiNormal:
			u.multipleNormalThrows(dartboard, position, u.accuracyModel)
		case Mode_SearchNormal:
			u.messageDisplay = "Click SEARCH to begin"
		case Mode_DrawCircle:
			fmt.Println("Click callback for draw-circle ")
			if u.circleDrawingState == drawCircleStateDrawing {
				mousePosition := boardgeo.CreateBoardPositionFromXY(g.GetMousePos(), u.dartboard.GetSquareDimension(),
					u.dartboardImageMin)
				circleRadiusPixels := boardgeo.PixelDistanceBetweenBoardPositions(u.dartboard.GetTracingCircleCenter(),
					mousePosition, u.dartboard.GetSquareDimension())
				if circleRadiusPixels > 0 {
					fmt.Printf("  Stop drawing circle, record result center %v, radius %d\n",
						u.dartboard.GetTracingCircleCenter(), circleRadiusPixels)
					pixelDiameter := float64(2 * circleRadiusPixels)
					fmt.Printf("  Circle diameter is %g pixels\n", pixelDiameter)
					fmt.Printf("  Square dimension is %g pixels\n", u.dartboard.GetSquareDimension())
					fmt.Printf("  Scoring area fraction is %g\n", boardgeo.ScoringAreaFraction)

					normalizedDiameter := pixelDiameter / (u.dartboard.GetSquareDimension() * boardgeo.ScoringAreaFraction)
					fmt.Println("Normalized diameter", normalizedDiameter)
					stdDeviation := normalizedDiameter / 2
					u.stdDevInputField = float32(stdDeviation)
					u.accuracyModel.SetStandardDeviation(stdDeviation)
					u.dartboard.SetDrawOneSigma(u.drawOneSigma, u.accuracyModel.GetSigmaRadius(1))
					u.dartboard.SetDrawTwoSigma(u.drawTwoSigma, u.accuracyModel.GetSigmaRadius(2))
					u.dartboard.SetDrawThreeSigma(u.drawThreeSigma, u.accuracyModel.GetSigmaRadius(3))
					u.dartboard.StopTracingCircle()
				}
			}
			// Since the radio button is still set for drawing, we reset the state to
			// Start in case they click the button again - that would be repeating the draw
			u.circleDrawingState = drawCircleStateStart
			u.messageDisplay = ""
			g.Update()
		case Mode_EmpricalStdDev:
			u.handleEmpiricalModeClick(position)
		default:
			panic("Invalid radio button value")
		}
	}
}

// StartDrawStdDevMode is called when the "Draw Standard Deviation" button is clicked
// We set a flag causing the next mouse click to be used to draw a circle representing the
// 2-standard deviation circle for the normal distribution
func (u *UserInterfaceInstance) StartDrawStdDevMode() {
	u.messageDisplay = "Draw 95% circle"
	u.scoreDisplay = ""
	u.circleDrawingState = drawCircleStateStart
	g.Update()
}

//		handleEmpiricalModeClick is called when the user clicks on the dartboard in "Measure Real Throws" mode
//	 What happens depends on the current state of the measurement process
func (u *UserInterfaceInstance) handleEmpiricalModeClick(position boardgeo.BoardPosition) {
	//fmt.Println("Handling empirical mode click at ", position)
	//fmt.Println("  Measurement state is ", u.measurementState)
	//	In "waiting for user to select target" state, when they click we record the target position
	//  and enter "gathering hits for this target" state
	switch u.measurementState {
	case measureStdDevStateSelectTarget:
		// We have just started a first, or new, target. The clicked position is where we will be throwing,
		// so we record this target position and move to "collecting hits" mode
		//fmt.Println("  Selecting target at ", position)
		u.measuringTarget = position
		u.measurementState = measureStdDevStateThrowing
		u.messageDisplay = "Throw, click hits"
		g.Update()
	case measureStdDevStateThrowing:
		//fmt.Println("  Throwing at target, hit at ", position)
		u.realThrows.AddHit(u.measuringTarget, position)
		if u.realThrows.IsStdDevAvailable() {
			stdDev := u.realThrows.CalcStdDevOfThrows()
			u.stdDevInputField = float32(stdDev)
			u.accuracyModel.SetStandardDeviation(stdDev)
			u.dartboard.SetDrawOneSigma(u.drawOneSigma, u.accuracyModel.GetSigmaRadius(1))
			u.dartboard.SetDrawTwoSigma(u.drawTwoSigma, u.accuracyModel.GetSigmaRadius(2))
			u.dartboard.SetDrawThreeSigma(u.drawThreeSigma, u.accuracyModel.GetSigmaRadius(3))
		}
		u.messageDisplay = "Throw, click hits"
	default:
		panic("  Invalid state for empirical mode click")
	}
}
