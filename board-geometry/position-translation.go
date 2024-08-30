package boardgeo

//	Functions to translate between mouse position, polar, and cartesian board coordinates

import (
	"image"
	"math"
)

// Diameters of the various circles of importance (in millimeters)
const displayedBoardDiameter = 451.0
const innerBullDiameter = 12.7
const outerBullDiameter = 31.8
const insideTrebleDiameter = 194.0
const outsideTrebleDiameter = 213.0
const insideDoubleDiameter = 321.0
const outsideDoubleDiameter = 340.0
const scoringAreaDiameter = outsideDoubleDiameter

// Those same quantities as radii
// const displayedBoardRadius = displayedBoardDiameter / 2.0
const innerBullRadius = innerBullDiameter / 2.0
const outerBullRadius = outerBullDiameter / 2.0
const insideTrebleRadius = insideTrebleDiameter / 2.0
const outsideTrebleRadius = outsideTrebleDiameter / 2.0
const insideDoubleRadius = insideDoubleDiameter / 2.0
const outsideDoubleRadius = outsideDoubleDiameter / 2.0
const scoringAreaRadius = scoringAreaDiameter / 2.0

// Those same quantities as radii normalized to 0 - 1
const innerBullRadiusNormalized = innerBullRadius / scoringAreaRadius
const outerBullRadiusNormalized = outerBullRadius / scoringAreaRadius
const insideTrebleRadiusNormalized = insideTrebleRadius / scoringAreaRadius
const outsideTrebleRadiusNormalized = outsideTrebleRadius / scoringAreaRadius
const insideDoubleRadiusNormalized = insideDoubleRadius / scoringAreaRadius
const outsideDoubleRadiusNormalized = outsideDoubleRadius / scoringAreaRadius
const scoringAreaRadiusNormalized = scoringAreaRadius / scoringAreaRadius

// ScoringAreaFraction is a scaling factor to normalize mouse positing inside board from 0 to 1 radius
const ScoringAreaFraction = float64(scoringAreaDiameter) / float64(displayedBoardDiameter)

// CreateBoardPositionFromXY is the basic function to convert from a clicked mouse position
// to a BoardPosition object
func CreateBoardPositionFromXY(mousePosition image.Point,
	squareDimension float64,
	imageMin image.Point) BoardPosition {

	xMouseInside := mousePosition.X - imageMin.X
	yMouseInside := mousePosition.Y - imageMin.Y

	xMouseZeroCentered := xMouseInside - int(math.Round(squareDimension/2))
	yMouseZeroCentered := -(yMouseInside - int(math.Round(squareDimension/2)))

	xFractionBoard := float64(xMouseZeroCentered) / (squareDimension / 2)
	yFractionBoard := float64(yMouseZeroCentered) / (squareDimension / 2)

	xFractionScoring := xFractionBoard / ScoringAreaFraction
	yFractionScoring := yFractionBoard / ScoringAreaFraction

	polarRadius := math.Sqrt(math.Pow(xFractionScoring, 2) + math.Pow(yFractionScoring, 2))
	polarTheta := math.Atan2(xFractionScoring, yFractionScoring)
	thetaAsDegrees := polarTheta * (180 / math.Pi)

	position := BoardPosition{
		//XMouseInside: xMouseInside,
		//YMouseInside: yMouseInside,
		Radius: polarRadius,
		Angle:  thetaAsDegrees,
	}

	return position
}

func GetXY(bp BoardPosition, squareDimension float64) (int, int) {
	//	Get the x,y equivalents
	xFromPolar := bp.Radius * math.Sin(bp.Angle*math.Pi/180)
	yFromPolar := bp.Radius * math.Cos(bp.Angle*math.Pi/180)

	//	polarRadius, thetaDegrees, xFromPolar, yFromPolar)
	xScaledByScoringFraction := xFromPolar * ScoringAreaFraction
	yScaledByScoringFraction := yFromPolar * ScoringAreaFraction

	xScaledToWindow := xScaledByScoringFraction * (squareDimension / 2)
	yScaledToWindow := yScaledByScoringFraction * (squareDimension / 2)

	xInsideWindow := int(math.Round(xScaledToWindow + squareDimension/2))
	yInsideWindow := int(math.Round(squareDimension/2 - yScaledToWindow))

	return xInsideWindow, yInsideWindow

}

// CreateBoardPositionFromPolar creates a BoardPosition object from polar coordinates.
// This is used when generating throws from random normal distributions, since those calculations
// are done in polar coordinates.
func CreateBoardPositionFromPolar(polarRadius float64, thetaDegrees float64) BoardPosition {

	//	Get the x,y equivalents
	//xFromPolar := polarRadius * math.Sin(thetaDegrees*math.Pi/180)
	//yFromPolar := polarRadius * math.Cos(thetaDegrees*math.Pi/180)

	//	polarRadius, thetaDegrees, xFromPolar, yFromPolar)
	//xScaledByScoringFraction := xFromPolar * ScoringAreaFraction
	//yScaledByScoringFraction := yFromPolar * ScoringAreaFraction

	//xScaledToWindow := xScaledByScoringFraction * (squareDimension / 2)
	//yScaledToWindow := yScaledByScoringFraction * (squareDimension / 2)

	//xInsideWindow := int(math.Round(xScaledToWindow + squareDimension/2))
	//yInsideWindow := int(math.Round(squareDimension/2 - yScaledToWindow))

	position := BoardPosition{
		//XMouseInside: xInsideWindow,
		//YMouseInside: yInsideWindow,
		Radius: polarRadius,
		Angle:  thetaDegrees,
	}
	return position
}
