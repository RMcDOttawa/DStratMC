package main

import "fmt"

func main() {
	fmt.Println("Main STUB")
	accuracyModel, err := CreateOrLoadAccuracyModel()
	if err != nil {
		fmt.Println("Error creating accuracy model:", err)
		return
	}
	simResults, err := RunSimulation(accuracyModel)
	if err != nil {
		fmt.Println("Error running simulation:", err)
	}
	err = ReportSimulationResults(simResults)
}
