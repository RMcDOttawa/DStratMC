package simulation

import (
	boardgeo "DStratMC/board-geometry"
	"image"
	"math"
	"math/rand"
)

// UniformAccuracyModel assumes that the result of a throw is evenly distributed in a circle
//
//	centered on the target and with a given radius
type UniformAccuracyModel struct {
	CEPRadius float64
}

func NewUniformAccuracyModel(CEPRadius float64) AccuracyModel {
	instance := &UniformAccuracyModel{
		CEPRadius: CEPRadius,
	}
	//fmt.Println("NewUniformAccuracyModel returns", instance)
	return instance
}

func (p UniformAccuracyModel) GetAccuracyRadius() float64 {
	return p.CEPRadius
}

func (p UniformAccuracyModel) GetThrow(target boardgeo.BoardPosition,
	scoringRadius float64,
	squareDimension float64,
	startPoint image.Point) (boardgeo.BoardPosition, error) {
	//	Polar coordinate deviation
	randomTheta := rand.Float64() * 2 * math.Pi
	randomRadius := p.CEPRadius * math.Sqrt(rand.Float64()) * scoringRadius
	//	Convert to cartesian
	newX := float64(target.XMouseInside) + randomRadius*math.Cos(randomTheta)
	newY := float64(target.YMouseInside) + randomRadius*math.Sin(randomTheta)
	//	Convert to board position
	point := image.Pt(int(math.Round(newX))+startPoint.X, int(math.Round(newY))+startPoint.Y)
	result := boardgeo.CreateBoardPositionFromXY(point, squareDimension, startPoint)
	return result, nil
}

func (p UniformAccuracyModel) GetSigmaRadius(_ float64) float64 {
	panic("GetSigmaRadius not meaningful for uniform model")
	return 0
}
