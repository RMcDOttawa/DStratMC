package simulation

type BoardPosition struct {
	radius float64 // 0 (centre) to 1.0 (outer edge of scoring area) or larger for outside
	angle  float64 // Degrees, clockwise from 0 being straight up
}
