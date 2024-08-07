package simulation

import (
	boardgeo "DStratMC/board-geometry"
	"fmt"
	"image"
	"math"
	"math/rand"
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
}

// PerfectAccuracyModel is a trivial implementation of the accuracy model where the result
// of a throw always hits the target
type PerfectAccuracyModel struct {
}

func (p PerfectAccuracyModel) GetAccuracyRadius() float64 {
	//TODO implement me
	panic("implement me")
}

func NewPerfectAccuracyModel() AccuracyModel {
	instance := &PerfectAccuracyModel{}
	fmt.Println("NewPerfectAccuracyModel returns", instance)
	return instance
}

func (p PerfectAccuracyModel) GetThrow(target boardgeo.BoardPosition,
	_ float64,
	_ float64,
	_ image.Point) (boardgeo.BoardPosition, error) {
	//fmt.Printf("PerfectAccuracyModel/GetThrow(%#v)\n", target)
	return target, nil
}

// CircularAccuracyModel assumes that the result of a throw is evenly distributed in a circle
//
//	centered on the target and with a given radius
type CircularAccuracyModel struct {
	CEPRadius float64
}

func NewCircularAccuracyModel(CEPRadius float64) AccuracyModel {
	instance := &CircularAccuracyModel{
		CEPRadius: CEPRadius,
	}
	//fmt.Println("NewCircularAccuracyModel returns", instance)
	return instance
}

func (p CircularAccuracyModel) GetAccuracyRadius() float64 {
	return p.CEPRadius
}

// Using cartesian
func (p CircularAccuracyModel) GetThrow(target boardgeo.BoardPosition,
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

// Using polar
//func (p CircularAccuracyModel) GetThrow(target boardgeo.BoardPosition,
//	scoringRadius float64,
//	squareDimension float64,
//	startPoint image.Point) (boardgeo.BoardPosition, error) {
//	fmt.Printf("CircularAccuracyModel/GetThrow target %v, CEP %g, sr %g, sd %g, stp %v\n",
//		target, p.CEPRadius, scoringRadius, squareDimension, startPoint)
//
//	//	Get a random point inside the circle centred at the given target and with radius CEP
//	randomThetaRadians := rand.Float64() * 2 * math.Pi
//	randomRadius := p.CEPRadius * math.Sqrt(rand.Float64())
//	fmt.Printf("  random theta %g, random radius %g\n", randomThetaRadians, randomRadius)
//
//	// Convert to absolute polar coordinates
//	radiusResult := math.Sqrt(target.Radius*target.Radius +
//		randomRadius*randomRadius +
//		2*target.Radius*randomRadius*math.Cos(randomThetaRadians))
//	thetaResult := target.Angle + math.Atan2(randomRadius*math.Sin(randomThetaRadians),
//		target.Radius+randomRadius*math.Cos(randomThetaRadians))
//	fmt.Printf("  absolute polar R=%g, theta=%g\n", radiusResult, thetaResult)
//
//	result := boardgeo.CreateBoardPositionFromXY()
//	fmt.Printf("   Result: %v\n", result)
//	return result, nil
//}
