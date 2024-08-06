package boardgeo

type BoardPosition struct {
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
