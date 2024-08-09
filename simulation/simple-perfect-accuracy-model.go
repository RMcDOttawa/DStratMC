package simulation

import (
	boardgeo "DStratMC/board-geometry"
	"fmt"
	"image"
)

// PerfectAccuracyModel is a trivial implementation of the accuracy model where the result
// of a throw always hits the target
type PerfectAccuracyModel struct {
}

func (p PerfectAccuracyModel) GetAccuracyRadius() float64 {
	panic("should not have beeen called")
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
