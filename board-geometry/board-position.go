package boardgeo

import (
	"fmt"
	"math"
	"strconv"
)

// BoardPosition represents a point on the dartboard - either where we aimed, or where we hit
// We record it both in terms of the mouse position (x,y) inside the frame,
// and in a normalized polar coordinate system
type BoardPosition struct {
	XMouseInside int
	YMouseInside int
	Radius       float64 // 0 (centre) to 1.0 (outer edge of scoring area) or larger for outside
	Angle        float64 // Degrees, clockwise from 0 being straight up
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

// Board colours
const (
	Board_Colour_Black = iota
	Board_Colour_White
	Board_Colour_Red
	Board_Colour_Green
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

// GetColourForSegment tells what colour the dartboard segment provided is.  This is not used for drawing -
// we are using a pre-drawn image to display the dartboard. It is used to calculate a contrasting colour for
// markers of various kinds that are displayed on top of the dartboard.
func GetColourForSegment(segment int, score int) int {
	if segment == BoardArea_Out {
		return Board_Colour_Black
	} else if segment == BoardArea_InnerSingle || segment == BoardArea_OuterSingle {
		return colourForSingle(score)
	} else if segment == BoardArea_Double {
		return colourForDouble(score)
	} else if segment == BoardArea_Treble {
		return colourForTreble(score)
	} else if segment == BoardArea_InnerBull {
		return Board_Colour_Red
	} else if segment == BoardArea_OuterBull {
		return Board_Colour_Green
	} else {
		panic("Unexpected board segment type: " + strconv.Itoa(segment))
	}
}

var multiplierRingColours = []int{
	Board_Colour_Green, //	1
	Board_Colour_Red,   //	2
	Board_Colour_Red,   //	3
	Board_Colour_Green, //	4
	Board_Colour_Green, //	5
	Board_Colour_Green, //	6
	Board_Colour_Red,   //	7
	Board_Colour_Red,   //	8
	Board_Colour_Green, //	9
	Board_Colour_Red,   //	10
	Board_Colour_Green, //	11
	Board_Colour_Red,   //	12
	Board_Colour_Red,   //	13
	Board_Colour_Red,   //	14
	Board_Colour_Green, //	15
	Board_Colour_Green, //	16
	Board_Colour_Green, //	17
	Board_Colour_Red,   //	18
	Board_Colour_Green, //	19
	Board_Colour_Red,   //	20
}

func colourForDouble(score int) int {
	return multiplierRingColours[(score/2)-1]
}

func colourForTreble(score int) int {
	return multiplierRingColours[(score/3)-1]
}

var singleSegmentColours = []int{
	Board_Colour_White, //	1
	Board_Colour_Black, //	2
	Board_Colour_Black, //	3
	Board_Colour_White, //	4
	Board_Colour_White, //	5
	Board_Colour_White, //	6
	Board_Colour_Black, //	7
	Board_Colour_Black, //	8
	Board_Colour_White, //	9
	Board_Colour_Black, //	10
	Board_Colour_White, //	11
	Board_Colour_Black, //	12
	Board_Colour_Black, //	13
	Board_Colour_Black, //	14
	Board_Colour_White, //	15
	Board_Colour_White, //	16
	Board_Colour_White, //	17
	Board_Colour_Black, //	18
	Board_Colour_White, //	19
	Board_Colour_Black, //	20
}

func colourForSingle(score int) int {
	if score < 1 || score > 20 {
		panic("colourForSingle received Invalid single score: " + strconv.Itoa(score))
	}
	return singleSegmentColours[score-1]
}

// GetContrastingColour returns a contrasting colour for the given colour.  This is used to make sure that
// text and other markers are visible on the board, regardless of the colour of the segment they are on.
// The colour is returned as 3
func GetContrastingColour(colour int) (uint8, uint8, uint8) {
	switch colour {
	case Board_Colour_Black:
		return 220, 220, 220
	case Board_Colour_White:
		return 50, 50, 50
	case Board_Colour_Red:
		return 56, 240, 140
	case Board_Colour_Green:
		return 240, 112, 249
	}
	return 0, 0, 0
}
