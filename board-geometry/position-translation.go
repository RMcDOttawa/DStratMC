package boardgeo

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

// Scaling factor to normalize mouse positing inside board from 0 to 1 radius
const ScoringAreaFraction = float64(scoringAreaDiameter) / float64(displayedBoardDiameter)

// func CreateBoardPositionFromXY(window *g.WindowWidget) BoardPosition {
func CreateBoardPositionFromXY(mousePosition image.Point,
	squareDimension float64,
	imageMin image.Point) BoardPosition {
	//fmt.Printf("CreateBoardPositionFromXY(%g,%v) mp %v\n", squareDimension, imageMin, mousePosition)
	xMouseInside := mousePosition.X - imageMin.X
	yMouseInside := mousePosition.Y - imageMin.Y
	//fmt.Printf("Absolute Mouse Position = (%d,%d), Relative Mouse = (%d,%d)\n",
	//	mousePosition.X, mousePosition.Y, xMouseInside, yMouseInside)

	xMouseZeroCentered := xMouseInside - int(math.Round(squareDimension/2))
	yMouseZeroCentered := -(yMouseInside - int(math.Round(squareDimension/2)))
	//fmt.Printf("Mouse centerd = (%d,%d)\n", xMouseZeroCentered, yMouseZeroCentered)

	xFractionBoard := float64(xMouseZeroCentered) / (squareDimension / 2)
	yFractionBoard := float64(yMouseZeroCentered) / (squareDimension / 2)

	xFractionScoring := xFractionBoard / ScoringAreaFraction
	yFractionScoring := yFractionBoard / ScoringAreaFraction

	polarRadius := math.Sqrt(math.Pow(xFractionScoring, 2) + math.Pow(yFractionScoring, 2))
	polarTheta := math.Atan2(xFractionScoring, yFractionScoring)
	thetaAsDegrees := polarTheta * (180 / math.Pi)

	position := BoardPosition{
		XMouseInside: xMouseInside,
		YMouseInside: yMouseInside,
		Radius:       polarRadius,
		Angle:        thetaAsDegrees,
	}
	return position
}

func GetDrawingXY(position BoardPosition) (int, int) {
	return position.XMouseInside, position.YMouseInside
}
