package boardgeo

import (
	"fmt"
	g "github.com/AllenDang/giu"
	"math"
	"strconv"
)

// Diameters of the various circles of importance (in millimeters)
const displayedBoardDiameter = 451
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

var multiplierBoundaries = [...]float64{
	0.0, // Exact centre of board
	innerBullRadiusNormalized,
	outerBullRadiusNormalized,
	insideTrebleRadiusNormalized,
	outsideTrebleRadiusNormalized,
	insideDoubleRadiusNormalized,
	outsideDoubleRadiusNormalized,
	scoringAreaRadiusNormalized,
}

var multiplierList = [...]int{
	2, //	Red bull is a double
	1, //	Green bull is a single
	1, //	Single area between bull and treble
	3, //	Treble ring
	1, //	Single area between treble and double
	2, //	Double ring
	0, //	Outside board
}

// Point values of silces, clockwise from the top
var segmentPointValues = [...]int{
	20, 1, 18, 4, 13, 6, 10, 15, 2, 17,
	3, 19, 7, 16, 8, 11, 14, 9, 12, 5}

// Scaling factor to normalize mouse positing inside board from 0 to 1 radius
const scoringAreaFraction = float64(scoringAreaDiameter) / float64(displayedBoardDiameter)

// DescribeBoardPoint describes a point on the board by the area and the score
func DescribeBoardPoint(point BoardPosition) (int, int, string) {
	//fmt.Printf("DescribeBoardPoint(%#v)\n", point)
	if point.Radius > 1 {
		return BoardArea_Out, 0, BoardAreaDescription[BoardArea_Out]
	}

	singlePointValue := determineSinglePointValue(math.Abs(point.Radius), point.Angle)
	multiplier := determineMultiplier(math.Abs(point.Radius))
	if singlePointValue == 0 || multiplier == 0 {
		return BoardArea_Out, 0, BoardAreaDescription[BoardArea_Out]
	}
	score := singlePointValue * multiplier

	if score == 50 {
		return BoardArea_InnerBull, 50, BoardAreaDescription[BoardArea_InnerBull]
	}
	if score == 25 {
		return BoardArea_OuterBull, 25, BoardAreaDescription[BoardArea_OuterBull]
	}
	if multiplier == 3 {
		asString := BoardAreaDescription[BoardArea_Treble] + " " + strconv.Itoa(singlePointValue)
		return BoardArea_Treble, score, asString
	}
	if multiplier == 2 {
		asString := BoardAreaDescription[BoardArea_Double] + " " + strconv.Itoa(singlePointValue)
		return BoardArea_Double, score, asString
	}
	if score == 0 {
		fmt.Println("Unexpected score 0")
	}

	segment := BoardArea_OuterSingle
	if math.Abs(point.Radius) < insideTrebleRadiusNormalized {
		segment = BoardArea_InnerSingle
	}
	asString := BoardAreaDescription[segment] + " " + strconv.Itoa(singlePointValue)
	return segment, score, asString
}

//	Use the angle on the board to determine which point wedge has been hit (ignoring double,triple)
//	Angles come in as numbers between -180 and +180.  We'll convert them to 0 to 360.
//	0 is straight up, then increasing numbers are clockwise rotation around the board.
//	Wedges are 360/20 = 18 degrees wide, and are offset by 1/2 that, 9 degrees.

func determineSinglePointValue(radius, degrees float64) int {
	//	Special cases: the inner and outer bulls - rotation doesn't matter
	if radius < outerBullRadiusNormalized {
		return 25
	}

	//	Convert +/- 180 range to 0-360 range
	if degrees < 0 {
		degrees = degrees + 360.0
	}

	//	Shift by 9 degrees (half a segment) to de-center the slices. So the beginning of the "1" slice is
	//	18 degrees, the 18 slice is 36 degrees, etc.
	degrees += (360.0 / 20) / 2

	//	Convert that to an integer index, where index 0 is the 20 slice, [1] is the 1 slice, [2] is the 18, etc.
	//	(20 will be messed up, we'll fix in a moment)
	sliceIndex := int(math.Floor(degrees / (360.0 / 20)))

	//	If we're on the left side of the 20 slice, this will have generated an index of 20, which we
	//	need to convert to 0
	if sliceIndex == 20 {
		sliceIndex = 0
	}
	return segmentPointValues[sliceIndex]
}

func determineMultiplier(radius float64) int {
	foundMultiplierIndex := -1
	for i := 0; i < len(multiplierBoundaries); i++ {
		thisBoundary := multiplierBoundaries[i]
		if radius < thisBoundary {
			foundMultiplierIndex = i - 1
			break
		}
	}
	if foundMultiplierIndex == -1 {
		return 0
	}
	return multiplierList[foundMultiplierIndex]
}

func CalcMousePolarPosition(window *g.WindowWidget) BoardPosition {
	wx32, wy32 := window.CurrentPosition()
	windowX := float64(wx32)
	windowY := float64(wy32)

	w32, h32 := window.CurrentSize()
	width := float64(w32)
	height := float64(h32)

	squareDimension := math.Min(width, height)
	xPadding := int(math.Round((width - squareDimension) / 2))
	yPadding := int(math.Round((height - squareDimension) / 2))
	//fmt.Printf("Window size (%d x %d), Padding (%d x %d)\n", int(width), int(height), xPadding, yPadding)
	mp := g.GetMousePos()
	xMouseInside := mp.X - int(windowX) - xPadding
	yMouseInside := mp.Y - int(windowY) - yPadding
	//fmt.Printf("Window Position = (%g,%g), Absolute Mouse Position = (%d,%d), Relative Mouse = (%d,%d)\n",
	//	windowX, windowY, mp.X, mp.Y, xMouseInside, yMouseInside)

	xMouseZeroCentered := xMouseInside - int(math.Round(squareDimension/2))
	yMouseZeroCentered := -(yMouseInside - int(math.Round(squareDimension/2)))

	xFractionBoard := float64(xMouseZeroCentered) / (squareDimension / 2)
	yFractionBoard := float64(yMouseZeroCentered) / (squareDimension / 2)

	xFractionScoring := xFractionBoard / scoringAreaFraction
	yFractionScoring := yFractionBoard / scoringAreaFraction

	polarRadius := math.Sqrt(math.Pow(xFractionScoring, 2) + math.Pow(yFractionScoring, 2))
	polarTheta := math.Atan2(xFractionScoring, yFractionScoring)
	thetaAsDegrees := polarTheta * (180 / math.Pi)

	position := BoardPosition{
		Radius: polarRadius,
		Angle:  thetaAsDegrees,
	}
	return position
}
