package simulation

import (
	boardgeo "DStratMC/board-geometry"
	"fmt"
	"math/rand"
)

// AccuracyModel is an abstract model that is used to determine a simulated thrower's accuracy -
// how close they will come to their intended target when they throw.
// Accuracy is determined by the type of accuracy model used - a variety of
// implementations will provide models of different levels of complexity.
type AccuracyModel interface {
	GetThrow(target boardgeo.BoardPosition) (boardgeo.BoardPosition, error)
}

// PerfectAccuracyModel is a trivial implementation of the accuracy model where the result
// of a throw always hits the target
type PerfectAccuracyModel struct{}

func NewPerfectAccuracyModel() AccuracyModel {
	instance := &PerfectAccuracyModel{}
	fmt.Println("NewPerfectAccuracyModel returns", instance)
	return instance
}

func (p PerfectAccuracyModel) GetThrow(target boardgeo.BoardPosition) (boardgeo.BoardPosition, error) {
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
	fmt.Println("NewCircularAccuracyModel returns", instance)
	return instance
}

func (p CircularAccuracyModel) GetThrow(target boardgeo.BoardPosition) (boardgeo.BoardPosition, error) {
	//fmt.Printf("CircularAccuracyModel/GetThrow(%#v) STUB\n", target)

	//	Get random deviation from target radius as +/- CEP
	signedDeviation := rand.Float64()*(2*p.CEPRadius) - p.CEPRadius
	//fmt.Printf("With CEP %g, radius deviation is %g\n", p.CEPRadius, signedDeviation)

	//	We'll pick a new angle that deviates from the target angle by the CEP factor.  Since the deviation may
	//	drive it out of the range (0 - 360) we normalize
	angleDeviationFactor := rand.Float64()*(2*p.CEPRadius) - p.CEPRadius
	unNormalizedAngle := target.Angle + target.Angle*angleDeviationFactor
	newAngle := unNormalizedAngle
	if newAngle < 0 {
		newAngle += 360
	}
	if newAngle > 360 {
		newAngle -= 360
	}

	//fmt.Printf("Angle %g, deviation factor %g, deviation %g, raw %g, new %g\n",
	//	target.angle, angleDeviationFactor, unNormalizedAngle, newAngle)

	result := boardgeo.BoardPosition{
		Radius: target.Radius + signedDeviation,
		Angle:  newAngle,
	}

	//fmt.Printf("Result: %v\n", result)
	return result, nil
}
