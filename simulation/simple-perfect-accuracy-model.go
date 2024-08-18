package simulation

import (
	boardgeo "DStratMC/board-geometry"
	"fmt"
	"image"
)

// PerfectAccuracyModel is a trivial implementation of the accuracy model where the result
// of a throw always hits the target exactly - no error or variation is introduced.
type PerfectAccuracyModel struct {
}

// NewPerfectAccuracyModel returns a new instance of the perfect accuracy model
func NewPerfectAccuracyModel() AccuracyModel {
	instance := &PerfectAccuracyModel{}
	fmt.Println("NewPerfectAccuracyModel returns", instance)
	return instance
}

// GetSigmaRadius should never be called with this instance of the accuracy model
func (p PerfectAccuracyModel) GetSigmaRadius(_ float64) float64 {
	panic("GetSigmaRadius not meaningful for perfect-accuracy model")
	return 0
}

// GetAccuracyRadius should never be called with this instance of the accuracy model
func (p PerfectAccuracyModel) GetAccuracyRadius() float64 {
	panic("should not have been called")
}

// GetThrow returns the target position as the result of the throw - perfect accuracy
func (p PerfectAccuracyModel) GetThrow(target boardgeo.BoardPosition,
	_ float64,
	_ float64,
	_ image.Point) (boardgeo.BoardPosition, error) {
	//fmt.Printf("PerfectAccuracyModel/GetThrow(%#v)\n", target)
	return target, nil
}
