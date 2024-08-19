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
	"math"
	"strconv"
)

// UserInterface models the overall UI as an object. Even thought there won't be multiple instances, modeling
// it this way allows control over the associated data, reducing the need for global variables and the errors they invite
type UserInterface interface {
	MainUiLoop()
}

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
}

// NewUserInterface creates a new UserInterface object
func NewUserInterface(loadedImage *image.RGBA) UserInterface {
	instance := &UserInterfaceInstance{
		mode:                       Mode_OneNormal,
		drawOneSigma:               false,
		drawTwoSigma:               false,
		drawThreeSigma:             false,
		searchShowEachTarget:       false,
		searchResultStrings:        [10]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
		dartboard:                  NewDartboard(),
		drawReferenceLinesCheckbox: true,
		numThrowsField:             throwsAtOneTarget,
		stdDevInputField:           0.15,
	}
	g.EnqueueNewTextureFromRgba(loadedImage, func(t *g.Texture) {
		instance.dartboardTexture = t
	})

	instance.dartboard.SetDrawRefLines(instance.drawReferenceLinesCheckbox)
	instance.dartboard.SetClickCallback(instance.dartboardClickCallback)
	return instance
}

// MainUiLoop is the main loop for the user interface. This called from the master window's Run method,
// 30 times a second repeatedly.  The GIU ui framework does not store state, so we are responsible for
// storing state information for anything that is to be displayed that seems to be constant on the screen

// The UI is divided into two sections: a vertical toolbar on the left, and a square dartboard area on the right.
func (u *UserInterfaceInstance) MainUiLoop() {
	window := u.setUpWindow()

	window.Layout(
		u.leftToolbarLayout(),
		g.Custom(u.dartboard.DrawFunction),
	)

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
	dartboardImageMin := image.Pt(int(windowX)+leftToolbarWidth, int(windowY))
	dartboardImageMax := image.Pt(dartboardImageMin.X+int(squareDimension), dartboardImageMin.Y+int(squareDimension))
	//fmt.Printf("image min %d, max %d\n", imageMin, imageMax)

	u.dartboard.SetInfo(window, u.dartboardTexture, squareDimension, dartboardImageMin, dartboardImageMax)
	return window
}

// leftToolbarLayout lays out the left toolbar, which contains buttons and fields for user interaction
//
//	Some of the elements in the toolbar are only displayed when certain radio buttons are selected
func (u *UserInterfaceInstance) leftToolbarLayout() g.Widget {
	u.accuracyModel = u.getAccuracyModel(u.mode)
	return g.Layout{
		//	Checkbox controlling whether crosshairs are drawn
		g.Checkbox("Reference Lines", &u.drawReferenceLinesCheckbox).OnChange(func() { u.dartboard.SetDrawRefLines(u.drawReferenceLinesCheckbox) }),

		// Fields used to select type of interaction and display messages
		u.uiLayoutInteractionTypeRadioButtons(),
		u.uiLayoutResetButton(),
		u.uiLayoutResultsMessage(),

		//	The following fields may be presented depending on the type of interaction
		u.uiLayoutNumberOfThrowsField(),
		u.uiLayoutStdDevField(),
		u.uiLayoutStdCircleCheckboxes(),
		u.uiShowSearchCheckbox(),
		u.uiLayoutSearchButton(),
		u.uiLayoutBlinkingSearchNotice(),
		u.uiLayoutSearchProgressBar(),
		u.uiLayoutCancelSearchButton(),
		u.uiLayoutSearchResults(),
		u.uiLayoutAverageScore(),
	}
}

// uiLayoutInteractionTypeRadioButtons lays out the radio buttons to select the type of interaction and model
func (u *UserInterfaceInstance) uiLayoutInteractionTypeRadioButtons() g.Widget {
	return g.Layout{
		g.Label(""),
		g.RadioButton("One Exact", u.mode == Mode_Exact).OnChange(func() {
			u.mode = Mode_Exact
			u.accuracyModel = u.getAccuracyModel(u.mode)
			u.radioChanged()
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
	}
}

// uiLayoutResetButton lays out the Reset button in the left toolbar
func (u *UserInterfaceInstance) uiLayoutResetButton() g.Widget {
	return g.Layout{
		g.Label(""),
		g.Button("Reset").OnClick(u.radioChanged),
	}
}

// uiLayoutOptionalResultsMessage displays a generic message and throw score
func (u *UserInterfaceInstance) uiLayoutResultsMessage() g.Widget {
	return g.Layout{
		g.Label(""),
		g.Label(u.messageDisplay),
		g.Label(u.scoreDisplay),
	}
}

// uiLayoutStdCircleCheckboxes will, when we are doing normal distribution (and only then) offer 3 checkboxes for drawing reference
// circles at 1, 2, and 3 standard deviations
func (u *UserInterfaceInstance) uiLayoutStdCircleCheckboxes() g.Widget {
	return g.Layout{
		g.Condition(u.mode == Mode_OneNormal || u.mode == Mode_MultiNormal || u.mode == Mode_SearchNormal,
			g.Layout{
				g.Label(""),
				g.Label("Show circles for:"),
				g.Checkbox("1 Sigma", &u.drawOneSigma).OnChange(func() { u.dartboard.SetDrawOneSigma(u.drawOneSigma, u.accuracyModel.GetSigmaRadius(1)) }),
				g.Checkbox("2 Sigma", &u.drawTwoSigma).OnChange(func() { u.dartboard.SetDrawTwoSigma(u.drawTwoSigma, u.accuracyModel.GetSigmaRadius(2)) }),
				g.Checkbox("3 Sigma", &u.drawThreeSigma).OnChange(func() { u.dartboard.SetDrawThreeSigma(u.drawThreeSigma, u.accuracyModel.GetSigmaRadius(3)) }),
			}, nil),
	}
}

// uiLayoutStdDevField displays a field to enter an floating point number for standard deviation
func (u *UserInterfaceInstance) uiLayoutStdDevField() g.Widget {
	return g.Layout{
		// If we are doing multiple throws, allow the user to set the number of throws
		g.Condition(u.mode == Mode_OneNormal || u.mode == Mode_MultiNormal || u.mode == Mode_SearchNormal,
			g.Layout{
				g.Label(""),
				g.InputFloat(&u.stdDevInputField).
					Label("StdDev 0-1").
					Size(stdDevTextWidth).
					OnChange(u.validateAndProcessStdDevField),
			}, nil),
	}
}

func (u *UserInterfaceInstance) validateAndProcessStdDevField() {
	if u.stdDevInputField < 0 {
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

// uiShowSearchCheckbox displays a checkbox that determines whether we show a target marker for the search while in progress
func (u *UserInterfaceInstance) uiShowSearchCheckbox() g.Widget {
	return g.Layout{
		g.Condition(u.mode == Mode_SearchNormal,
			g.Layout{
				g.Label(""),
				g.Checkbox("Show Search", &u.searchShowEachTarget),
			}, nil),
	}
}

// uiLayoutNumberOfThrowsField displays a field to enter an integer number of throws
func (u *UserInterfaceInstance) uiLayoutNumberOfThrowsField() g.Widget {
	return g.Layout{
		// If we are doing multiple throws, allow the user to set the number of throws
		g.Condition(u.mode == Mode_MultiAvg || u.mode == Mode_MultiNormal || u.mode == Mode_SearchNormal,
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
		g.Condition(u.mode == Mode_SearchNormal,
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
		g.Condition(u.mode == Mode_SearchNormal,
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

// uiLayoutSearchProgressBar will, If we are doing a search, display a progress bar for search progress
func (u *UserInterfaceInstance) uiLayoutSearchProgressBar() g.Widget {
	return g.Layout{
		g.Condition(u.mode == Mode_SearchNormal,
			g.Layout{
				g.Label(""),
				g.ProgressBar(float32(u.searchProgressPercent)).Size(LeftToolbarMinimumWidth-10, 0),
			}, nil),
	}
}

// uiLayoutCancelSearchButton provides a button to cancel search, once search is running
func (u *UserInterfaceInstance) uiLayoutCancelSearchButton() g.Widget {
	return g.Layout{
		g.Condition(u.cancelSearchVisible,
			g.Layout{
				g.Label(""),
				g.Button("Cancel Search").OnClick(func() {
					fmt.Println("Cancelling Search")
					u.cancelSearch()
				}),
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
		g.Condition(u.mode == Mode_SearchNormal && u.searchComplete,
			g.Layout{
				g.Label("Best 10 throws:"),
				g.Label(""),
				u.uiLayoutSearchResultLabels(10),
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
}

// dartboardClickCallback is called when the user clicks on the dartboard. It is the main entry point for
// the UI to respond to user input
func (u *UserInterfaceInstance) dartboardClickCallback(dartboard Dartboard, position boardgeo.BoardPosition) {
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

	if position.Radius <= 1.0 {
		u.messageDisplay = ""
		u.scoreDisplay = ""
		dartboard.RemoveThrowMarkers()
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
		default:
			panic("Invalid radio button value")
		}
	}
}
