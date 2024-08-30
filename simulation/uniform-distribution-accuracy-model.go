package simulation

import (
	boardgeo "DStratMC/board-geometry"
	"image"
	"math"
	"math/rand"
)

// UniformAccuracyModel assumes that the result of a throw is evenly distributed in a circle
//	centered on the target and with a given radius

type UniformAccuracyModel struct {
	CEPRadius float64
}

// NewUniformAccuracyModel creates a new instance of the UniformAccuracyModel
func NewUniformAccuracyModel(CEPRadius float64) AccuracyModel {
	instance := &UniformAccuracyModel{
		CEPRadius: CEPRadius,
	}
	return instance
}

// GetAccuracyRadius returns the radius of the circle in which the throw will land
func (p UniformAccuracyModel) GetAccuracyRadius() float64 {
	return p.CEPRadius
}

// GetThrow returns the result of a throw, given a target position, a scoring radius, a square dimension, and a starting point
func (p UniformAccuracyModel) GetThrow(target boardgeo.BoardPosition,
	scoringRadius float64,
	squareDimension float64,
	startPoint image.Point) (boardgeo.BoardPosition, error) {
	//	Polar coordinate deviation
	randomTheta := rand.Float64() * 2 * math.Pi
	randomRadius := p.CEPRadius * math.Sqrt(rand.Float64()) * scoringRadius
	//	Convert to cartesian
	targetX, targetY := boardgeo.GetXY(target, squareDimension)
	newX := float64(targetX) + randomRadius*math.Cos(randomTheta)
	newY := float64(targetY) + randomRadius*math.Sin(randomTheta)
	//	Convert to board position
	point := image.Pt(int(math.Round(newX))+startPoint.X, int(math.Round(newY))+startPoint.Y)
	result := boardgeo.CreateBoardPositionFromXY(point, squareDimension, startPoint)
	return result, nil
}

// GetSigmaRadius is not meaningful for the uniform model
func (p UniformAccuracyModel) GetSigmaRadius(_ float64) float64 {
	panic("GetSigmaRadius not meaningful for uniform model")
	return 0
}

func (p UniformAccuracyModel) SetStandardDeviation(_ float64) {
	panic("SetStandardDeviation not meaningful for uniform model")
}
