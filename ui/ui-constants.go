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
)

// Certain fixed UI sizes that I can't be bothered to figure out how to compute at runtime
const LeftToolbarMinimumWidth = 200
const singleHitMarkerRadius = 5
const multipleHitMarkerRadius = 1

// How many throws are used in a multi-throw averaging run
const throwsAtOneTarget = 5_000
const numThrowsTextWidth = 120

//	Eventually the following will become computed variables: the size of the target circle
//	for uniform modeling, or of the 2-standard-deviation circle for normal modeling.
//	(CEP is Circular Error Probable, a nod to terminology about the accuracy of ballistics)

const uniformCEPRadius = 0.3

//const normalCEPRadius = 0.3

const normalCEPRadius = 0.25

// const normalCEPRadius = 0.1
const stubStandardDeviation = normalCEPRadius * 2

const testCoordinateConversion = true
