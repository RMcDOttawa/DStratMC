package boardgeo

import (
	"fmt"
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

// func CreateBoardPosition(window *g.WindowWidget) BoardPosition {
func CreateBoardPosition(mousePosition image.Point, squareDimension float64, imageMin image.Point, _ image.Point) BoardPosition {
	fmt.Printf("CreateBoardPosition(%g,%v) mp %v\n", squareDimension, imageMin, mousePosition)
	xMouseInside := mousePosition.X - imageMin.X
	yMouseInside := mousePosition.Y - imageMin.Y
	fmt.Printf("Absolute Mouse Position = (%d,%d), Relative Mouse = (%d,%d)\n",
		mousePosition.X, mousePosition.Y, xMouseInside, yMouseInside)

	xMouseZeroCentered := xMouseInside - int(math.Round(squareDimension/2))
	yMouseZeroCentered := -(yMouseInside - int(math.Round(squareDimension/2)))
	fmt.Printf("Mouse centred = (%d,%d)\n", xMouseZeroCentered, yMouseZeroCentered)

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

//func PositionToXY(position BoardPosition, squareDimension float64, imageMin image.Point, imageMax image.Point) (int, int) {
//	angleInRadians := position.Angle * math.Pi / 180.0
//	xAsFraction := position.Radius * math.Sin(angleInRadians)
//	yAsFraction := position.Radius * math.Cos(angleInRadians)
//
//	xZeroCentered := xAsFraction * (squareDimension / 2)
//	yZeroCentered := yAsFraction * (squareDimension / 2)
//
//	xPositive := xZeroCentered + (squareDimension / 2)
//	yPositive := (squareDimension / 2) - yZeroCentered
//
//	fmt.Printf("PositionToXY(%#v) gives (%g,%g)\n", position, xAsFraction, yAsFraction)
//	fmt.Printf("  x zero-centered %g, y zero-centered %g\n", xZeroCentered, yZeroCentered)
//	fmt.Printf("  x positive %g, y positive %g\n", xPositive, yPositive)
//
//	xStub := (imageMin.X + imageMax.X) / 2
//	yStub := (imageMin.Y + imageMax.Y) / 2
//	return xStub + 20, yStub + 20
//}

func GetDrawingXY(position BoardPosition) (int, int) {
	return position.XMouseInside, position.YMouseInside
}
