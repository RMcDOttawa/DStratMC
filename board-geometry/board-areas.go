package boardgeo

// BoardArea represents the different logical areas of the dartboard
type BoardArea int

// enum to represent the various scoring areas on the board
const (
	BoardArea_Out BoardArea = iota
	BoardArea_InnerBull
	BoardArea_OuterBull
	BoardArea_InnerSingle
	BoardArea_OuterSingle
	BoardArea_Double
	BoardArea_Treble
)

var BoardAreaDescription = map[BoardArea]string{
	BoardArea_Out:         "Out",
	BoardArea_InnerBull:   "Red Bull",
	BoardArea_OuterBull:   "Green Bull",
	BoardArea_InnerSingle: "Inner",
	BoardArea_OuterSingle: "Outer",
	BoardArea_Double:      "Double",
	BoardArea_Treble:      "Treble",
}
