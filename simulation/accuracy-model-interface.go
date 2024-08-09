package simulation

import (
	boardgeo "DStratMC/board-geometry"
	"image"
)

// AccuracyModel is an abstract model that is used to determine a simulated thrower's accuracy -
// how close they will come to their intended target when they throw.
// Accuracy is determined by the type of accuracy model used - a variety of
// implementations will provide models of different levels of complexity.
type AccuracyModel interface {
	GetThrow(target boardgeo.BoardPosition,
		scoringRadius float64,
		squareDimension float64,
		startPoint image.Point) (boardgeo.BoardPosition, error)
	GetAccuracyRadius() float64
	GetSigmaRadius(numSigmas float64) float64
}
