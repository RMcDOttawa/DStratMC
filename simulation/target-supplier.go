package simulation

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
	HasNext() bool
	NextTarget() boardgeo.BoardPosition // Radius, theta
}

type CircularTargetSupplierInstance struct {
	nextRadius      float64
	nextAngle       float64
	radiusIncrement float64
	angleIncrement  float64
	squareDimension float64
	imageMinPoint   image.Point
	windowX         int
	windowY         int
}

func NewTargetSupplier(squareDimension float64, imageMinPoint image.Point, windowX int, windowY int) TargetSupplier {
	instance := &CircularTargetSupplierInstance{
		nextRadius:      0.0,
		nextAngle:       0.0,
		radiusIncrement: 0.1,
		angleIncrement:  1.0,
		squareDimension: squareDimension,
		imageMinPoint:   imageMinPoint,
		windowX:         windowX,
		windowY:         windowY,
	}
	//fmt.Println("NewTargetSupplier returns", instance)
	return instance
}

func (t *CircularTargetSupplierInstance) HasNext() bool {
	return t.nextRadius <= 1.0
}

func (t *CircularTargetSupplierInstance) NextTarget() boardgeo.BoardPosition {
	//result := boardgeo.BoardPosition{
	//	Radius: t.nextRadius,
	//	Angle:  t.nextAngle,
	//}
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
