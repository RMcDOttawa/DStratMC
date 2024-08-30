package boardgeo

import (
	"fmt"
	"math"
	"strconv"
)

// BoardPosition represents a point on the dartboard - either where we aimed, or where we hit
// We record it both in terms of the mouse position (x,y) inside the frame,
// and in a normalized polar coordinate system
//
//	The x,y position has the (0,0) origin at the top left corner of the frame, with x increasing to the right
//	and y increasing downwards.  The x,y position is in pixels.
//	The polar coordinate system has the origin at the centre of the board, with the radius normalized to 0-1
//	and the angle in degrees, clockwise from the top of the board.
type BoardPosition struct {
	Radius float64 // 0 (centre) to 1.0 (outer edge of scoring area) or larger for outside
	Angle  float64 // Degrees, clockwise from 0 being straight up
}

// DescribeBoardPoint describes a point on the board by the area code, the score, and a text description
func DescribeBoardPoint(point BoardPosition) (BoardArea, int, string) {

	//	Immediately handle the case where the point is outside the board's scoring area
	if point.Radius > 1 {
		return BoardArea_Out, 0, BoardAreaDescription[BoardArea_Out]
	}

	singlePointValue := determineSinglePointValue(math.Abs(point.Radius), point.Angle)
	multiplier := determineMultiplier(math.Abs(point.Radius))
	if singlePointValue == 0 || multiplier == 0 {
		return BoardArea_Out, 0, BoardAreaDescription[BoardArea_Out]
	}
	score := singlePointValue * multiplier

	//	Handle the special cases of the bullseyes
	if score == 50 {
		return BoardArea_InnerBull, 50, BoardAreaDescription[BoardArea_InnerBull]
	}
	if score == 25 {
		return BoardArea_OuterBull, 25, BoardAreaDescription[BoardArea_OuterBull]
	}

	//	Handle the special cases of the double and treble rings
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

	//	Handle the single score segments
	segment := BoardArea_OuterSingle
	if math.Abs(point.Radius) < insideTrebleRadiusNormalized {
		segment = BoardArea_InnerSingle
	}
	asString := BoardAreaDescription[segment] + " " + strconv.Itoa(singlePointValue)
	return segment, score, asString
}

func PixelDistanceBetweenBoardPositions(a BoardPosition, b BoardPosition, squareDimension float64) int {
	aX, aY := GetXY(a, squareDimension)
	bX, bY := GetXY(b, squareDimension)
	xDiff := float64(aX - bX)
	yDiff := float64(aY - bY)
	distance := math.Sqrt(xDiff*xDiff + yDiff*yDiff)
	return int(math.Round(distance))

}
