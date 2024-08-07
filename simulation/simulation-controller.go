package simulation

import (
	"fmt"
)

type Simulation interface {
	//RunSimulation(model AccuracyModel) (SimResults, error)
}

type InstanceSimulation struct {
	accuracyModel AccuracyModel
}

func NewSimulation(accuracyModel AccuracyModel) Simulation {
	instance := &InstanceSimulation{
		accuracyModel: accuracyModel,
	}
	fmt.Println("NewSimulation returns", instance)
	return instance
}

//func (sim InstanceSimulation) RunSimulation(model AccuracyModel) (SimResults, error) {
//	fmt.Println("RunSimulation STUB", sim)
//	result := NewSimResults()
//	targetSupplier := NewTargetSupplier()
//	for targetSupplier.HasNext() {
//		newTarget := targetSupplier.NextTarget()
//		err := sim.throwsAtTarget(newTarget, numThrowsAtOneTarget, &result)
//		if err != nil {
//			fmt.Println("throwsAtTarget error", err)
//			continue
//		}
//	}
//	return result, nil
//}

//func (sim InstanceSimulation) throwsAtTarget(target boardgeo.BoardPosition, numThrows int, result *SimResults) error {
//	for i := 0; i < numThrows; i++ {
//		hit, err := sim.accuracyModel.GetThrow(target)
//		if err != nil {
//			fmt.Println("GetThrow error", err)
//			continue
//		}
//		_, score, _ := boardgeo.DescribeBoardPoint(hit)
//		(*result).AddThrow(target, score)
//	}
//	return nil
//}
