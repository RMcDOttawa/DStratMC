package target_search

import (
	boardgeo "DStratMC/board-geometry"
	"image"
	"math"
)

// 	TargetSupplier is an iterator that supplies throw targets for the simulation, one at a time
// 	Implementations may use simple mechanical coverage, or might do something more human-like, such
//	as focusing on specific visual targets on the board (like bull or treble).  The "next" function
//	should return the next target, and guarantee that each target is unique. The order in which
//	targets are returned is not defined. Simple implementations will do something sequential, but
//	random orders are possible and OK.

type TargetSupplier interface {
	HasNext() bool // True if there are more targets to return
	NextTarget() boardgeo.BoardPosition
	ForecastNumTargets() float64
}

// 	CircularTargetSupplierInstance is a simple implementation of TargetSupplier that returns targets
//	in a circular pattern, starting from the centre and spiralling outwards.

type CircularTargetSupplierInstance struct {
	nextRadius      float64
	nextAngle       float64
	radiusIncrement float64
	angleIncrement  float64
	squareDimension float64
	imageMinPoint   image.Point
}

const radiusIncrement = 0.02
const angleIncrement = 0.5

// NewTargetSupplier creates a new instance of CircularTargetSupplierInstance, with the given squareDimension
func NewTargetSupplier(squareDimension float64, imageMinPoint image.Point) TargetSupplier {
	instance := &CircularTargetSupplierInstance{
		nextRadius:      0.0,
		nextAngle:       0.0,
		radiusIncrement: radiusIncrement,
		angleIncrement:  angleIncrement,
		squareDimension: squareDimension,
		imageMinPoint:   imageMinPoint,
	}
	return instance
}

func (t *CircularTargetSupplierInstance) ForecastNumTargets() float64 {
	numRadiusSteps := 1.0 / t.radiusIncrement
	numAngleSteps := 360 / t.angleIncrement
	return numRadiusSteps * numAngleSteps
}

// HasNext returns true if there are more targets to return
func (t *CircularTargetSupplierInstance) HasNext() bool {
	return t.nextRadius <= 1.0
}

// NextTarget returns the next target in the sequence.
// To find the next target:
// 1.  Increment the angle, so we are working our way around the circle at the current distance from centre
// 2.  If we have gone all the way around, reset the angle to zero and increment the radius
func (t *CircularTargetSupplierInstance) NextTarget() boardgeo.BoardPosition {
	result := boardgeo.CreateBoardPositionFromPolar(t.nextRadius, t.nextAngle,
		t.squareDimension)
	//	Prepare for next results.  Rotate the angle.  At 360 degrees, reset angle to 0 and incrmenet radius
	//	Special case: radius zero is the centre - rotating angle is meaningless, so skip directly to next radius
	if t.nextRadius == 0.0 {
		t.nextRadius = AddWithoutNoise(t.nextRadius, t.radiusIncrement)
	} else {
		t.nextAngle = AddWithoutNoise(t.nextAngle, t.angleIncrement)
		if t.nextAngle >= 360.0 {
			t.nextAngle = 0.0
			t.nextRadius = AddWithoutNoise(t.nextRadius, t.radiusIncrement)
		}
	}
	return result
}

const noiseDecimalPlacesFactor = 10_000_000.0

//  AddWithoutNoise adds the given two floating point numbers, but then rounds the result to
//	a bunch of decimal places, to avoid the weird floating point result that sometimes, e.g.
//	0.1 + 0.1 = 0.199999999999999999 instead of 0.2

func AddWithoutNoise(a float64, b float64) float64 {
	sum := a + b
	rounded := math.Round(sum*noiseDecimalPlacesFactor) / noiseDecimalPlacesFactor
	return rounded
}
