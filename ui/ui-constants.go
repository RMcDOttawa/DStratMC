package ui

//	Various constants used by the User Interface

//		A radio button allows the user to select the over all mode: data aquisition, simple
//	 click response, or optimal target search
type InterfaceMode int

const (
	Mode_Exact        InterfaceMode = iota // Just record exact hit where clicked
	Mode_OneAvg                            // Record one hit uniformly distributed within a circle
	Mode_MultiAvg                          // Record multiple hits uniformly distributed within a circle
	Mode_OneNormal                         // Record one hit normally distributed within a circle
	Mode_MultiNormal                       // Record multiple hits normally distributed within a circle
	Mode_SearchNormal                      // Search around the board, recording result of multi-normal at each search location
	Mode_DrawCircle                        // Draw the 2-sigma (95%) standard deviation circle
)

// Certain fixed UI sizes that I can't be bothered to figure out how to compute at runtime
const LeftToolbarMinimumWidth = 200
const LeftToolbarChildWidth = LeftToolbarMinimumWidth - 15
const singleHitMarkerRadius = 5
const multipleHitMarkerRadius = 1
const numThrowsTextWidth = 120
const stdDevTextWidth = 64.0

const uiFramePadVertical = 2
const uiRadioButtonHeight = 22
const uiCheckboxHeight = uiRadioButtonHeight
const uiLabelHeight = 26
const uiButtonHeight = uiRadioButtonHeight
const uiInputFieldHeight = 36
const uiProgressBarHeight = 20

// How many throws are used in a multi-throw averaging run
const throwsAtOneTarget = 5000

const numSearchResultsToDisplay = 10

//	Eventually the following will become computed variables:

// the size of the target circle for uniform modeling,
const uniformCEPRadius = 0.3

const MasterWindowWidth = 1000 + LeftToolbarMinimumWidth
const MasterWindowHeight = 1000

const testCoordinateConversion = true
