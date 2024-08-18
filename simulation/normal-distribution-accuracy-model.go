package simulation

// NormalAccuracyModel assumes that the result of a throw follows a normal distribution, centered on the target
// and with a given standard deviation centered on the target and with a given radius

import (
	boardgeo "DStratMC/board-geometry"
	"gonum.org/v1/gonum/stat/distuv"
	"image"
	"math"
)

type NormalAccuracyModel struct {
	CEPRadius          float64 // Temporary. Eventually won't need this - just use the standard deviation
	standardDeviation  float64
	normalDistribution distuv.Normal
}

// NewNormalAccuracyModel creates a new instance of the NormalAccuracyModel
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

// GetAccuracyRadius should never be called with this instance of the accuracy model
func (p NormalAccuracyModel) GetAccuracyRadius() float64 {
	panic("GetAccuracyRadius not meaningful for normal model")
	return 0.5
}

// GetThrow generates a throw based on a normal distribution
//
//	We are given the coordinates the player actually aimed at, and use the normal distribution to determine
//	where the dart actually lands.
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

// GetSigmaRadius returns the radius of a circle that contains a given number of standard deviations
// For readers not versed in statistics, the "sigma" is the standard deviation of a normal distribution.
// For a normal distribution, 68% of the data falls within 1 sigma of the mean, 95% within 2 sigmas, and 99.7% within 3 sigmas.
// So a "2 sigma" circle would represent the area that you would expect most darts to land.
// A 3 sigma circle should catch almost all darts - darts outside this circle would be classified "wild throws" or "outliers".
func (p NormalAccuracyModel) GetSigmaRadius(numSigmas float64) float64 {
	return p.CEPRadius * numSigmas * p.standardDeviation
}
