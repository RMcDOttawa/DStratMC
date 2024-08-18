package boardgeo

import "strconv"

// BoardColour is an abstract representation of the colour of a segment of the dartboard
// These are just abstract descriptions of the colour - not RGB specifications for drawing

type BoardColour int

const (
	Board_Colour_Black BoardColour = iota
	Board_Colour_White
	Board_Colour_Red
	Board_Colour_Green
)

// GetColourForSegment tells what colour the dartboard segment provided is.  This is not used for drawing -
// we are using a pre-drawn image to display the dartboard. It is used to calculate a contrasting colour for
// markers of various kinds that are displayed on top of the dartboard.
func GetColourForSegment(segment BoardArea, score int) BoardColour {
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
		panic("Unexpected board segment type: " + strconv.Itoa(int(segment)))
	}
}

// The double and treble rings alternate colours around the board.
// These are the colours for the double and treble rings for each score
// e.g., the first item says that the double and treble rings for the "1" segment are green
// (Wouldn't it be nice if Go had constant arrays?)
var multiplierRingColours = []BoardColour{
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

// colourForDouble returns the colour for the double ring for the given score
// (the score has already been doubled, so we need to divide it by 2 to get the underlying single value)
func colourForDouble(score int) BoardColour {
	if score%2 != 0 {
		panic("colourForDouble received Invalid double score: " + strconv.Itoa(score))
	}
	return multiplierRingColours[(score/2)-1]
}

// colourForTreble returns the colour for the treble ring for the given score
// (the score has already been tripled, so we need to divide it by 3 to get the underlying single value)
func colourForTreble(score int) BoardColour {
	if score%3 != 0 {
		panic("colourForTreble received Invalid triple score: " + strconv.Itoa(score))
	}
	return multiplierRingColours[(score/3)-1]
}

var singleSegmentColours = []BoardColour{
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

func colourForSingle(score int) BoardColour {
	if score < 1 || score > 20 {
		panic("colourForSingle received Invalid single score: " + strconv.Itoa(score))
	}
	return singleSegmentColours[score-1]
}

// GetContrastingColour returns a contrasting colour for the given colour.  This is used to make sure that
// text and other markers are visible on the board, regardless of the colour of the segment they are on.
// The colour is returned as 3
func GetContrastingColour(colour BoardColour) (uint8, uint8, uint8) {
	switch colour {
	case Board_Colour_Black:
		return 220, 220, 220
	case Board_Colour_White:
		return 50, 50, 50
	case Board_Colour_Red:
		return 10, 128, 36
	case Board_Colour_Green:
		return 247, 33, 89
	}
	return 0, 0, 0
}
