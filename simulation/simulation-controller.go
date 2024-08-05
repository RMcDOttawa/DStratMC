package simulation

import "fmt"

type Simulation interface {
	RunSimulation(model AccuracyModel) (SimResults, error)
}

type InstanceSimulation struct {
}

func NewSimulation() Simulation {
	instance := &InstanceSimulation{}
	fmt.Println("NewSimulation returns", instance)
	return instance
}

func (sim InstanceSimulation) RunSimulation(model AccuracyModel) (SimResults, error) {
	fmt.Println("RunSimulation STUB", sim)
	result := Ne
	return nil, nil
}
