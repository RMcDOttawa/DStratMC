package simulation

import (
	boardgeo "DStratMC/board-geometry"
	"gonum.org/v1/gonum/stat/distuv"
	"image"
	"math"
)

// NormalAccuracyModel assumes that the result of a throw follows a normal distribution, centered on the target
// and with a given standard deviation
//
//	centered on the target and with a given radius
type NormalAccuracyModel struct {
	CEPRadius          float64
	standardDeviation  float64
	normalDistribution distuv.Normal
}

func NewNormalAccuracyModel(CEPRadius float64, stdDev float64) AccuracyModel {
	instance := &NormalAccuracyModel{
		CEPRadius:         CEPRadius,
		standardDeviation: stdDev,
		normalDistribution: distuv.Normal{
			Mu:    0.0,
			Sigma: stdDev,
		},
	}
	return instance
}

func (p NormalAccuracyModel) GetAccuracyRadius() float64 {
	panic("GetAccuracyRadius not meaningful for normal model")
	return 0.5
}

func (p NormalAccuracyModel) GetThrow(target boardgeo.BoardPosition,
	scoringRadius float64,
	squareDimension float64,
	startPoint image.Point) (boardgeo.BoardPosition, error) {

	// Generate normally distributed random offsets
	randomXDeviation := p.normalDistribution.Rand()
	randomYDeviation := p.normalDistribution.Rand()
	deltaX := randomXDeviation * scoringRadius * p.CEPRadius
	deltaY := randomYDeviation * scoringRadius * p.CEPRadius

	// Calculate the final coordinates in Cartesian form
	xFinal := target.XMouseInside + int(math.Round(deltaX))
	yFinal := target.YMouseInside + int(math.Round(deltaY))
	point := image.Pt(xFinal+startPoint.X, yFinal+startPoint.Y)

	result := boardgeo.CreateBoardPositionFromXY(point, squareDimension, startPoint)
	return result, nil
}

func (p NormalAccuracyModel) GetSigmaRadius(numSigmas float64) float64 {
	return p.CEPRadius * numSigmas * p.standardDeviation
}
