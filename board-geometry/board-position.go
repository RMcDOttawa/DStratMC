package boardgeo

import (
	"fmt"
	"math"
	"strconv"
)

type BoardPosition struct {
	XMouseInside int
	YMouseInside int
	//XMouseZeroCentered int
	//YMouseZeroCentered int
	//XFractionBoard     float64
	//YFractionBoard     float64
	//XFractionScoring   float64
	//YFractionScoring   float64
	Radius float64 // 0 (centre) to 1.0 (outer edge of scoring area) or larger for outside
	Angle  float64 // Degrees, clockwise from 0 being straight up
}

// enum to represent the various scoring areas on the board
const (
	BoardArea_Out = iota
	BoardArea_InnerBull
	BoardArea_OuterBull
	BoardArea_InnerSingle
	BoardArea_OuterSingle
	BoardArea_Double
	BoardArea_Treble
)

// BoardAreaDescription translates board position areas to text (change this for different languages)
var BoardAreaDescription = map[int]string{
	BoardArea_Out:         "Out",
	BoardArea_InnerBull:   "Red Bull",
	BoardArea_OuterBull:   "Green Bull",
	BoardArea_InnerSingle: "Inner",
	BoardArea_OuterSingle: "Outer",
	BoardArea_Double:      "Double",
	BoardArea_Treble:      "Treble",
}

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
