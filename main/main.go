package main

import (
	boardgeo "DStratMC/board-geometry"
	"DStratMC/simulation"
	"DStratMC/ui"
	"fmt"
	g "github.com/AllenDang/giu"
	"image"
)

var LoadedImage *image.RGBA

func main() {
	//modelPath := ""
	//if len(os.Args) > 1 {
	//	modelPath = os.Args[1]
	//}
	//accuracyModel, err := CreateOrLoadAccuracyModel(modelPath)
	//if err != nil {
	//	fmt.Println("Error creating accuracy model:", err)
	//	return
	//}

	//simulator := simulation.NewSimulation(accuracyModel)
	wnd := g.NewMasterWindow("Dartboard", 1000, 1000, 0)
	var err error
	LoadedImage, err = g.LoadImage("./Dartboard Illustration.png")
	if err != nil {
		fmt.Println("Unable to load dartboard image:", err)
		return
	}
	g.EnqueueNewTextureFromRgba(LoadedImage, func(t *g.Texture) {
		ui.DartboardInfo.Texture = t
	})
	wnd.Run(ui.MainUiLoop)

}

// CreateOrLoadAccuracyModel loads an existing accuracy model if specified, load it.
// If not, it creates a new empty model.
func CreateOrLoadAccuracyModel(modelPath string) (simulation.AccuracyModel, error) {
	if modelPath != "" {
		panic("Loading model from file not implemented yet")
	}
	fmt.Printf("CreateOrLoadAccuracyModel(%s) STUB\n", modelPath)
	return simulation.NewCircularAccuracyModel(0.1), nil
	//return simulation.NewCircularAccuracyModel(0.2), nil
	//return simulation.NewCircularAccuracyModel(0.23), nil
	//return simulation.NewCircularAccuracyModel(0.3), nil
	//return simulation.NewCircularAccuracyModel(0.4), nil
	//return simulation.NewPerfectAccuracyModel(), nil
}

func ReportSimulationResults(results simulation.SimResults) error {
	positions := results.GetPositionsSortedByHighScore()
	const justDoTheFirst = 10
	for _, p := range positions[:justDoTheFirst] {
		averageScore := results.GetAverageScore(p)
		_, _, description := boardgeo.DescribeBoardPoint(p)
		fmt.Printf("Position %v (%s) average score %g\n", p, description, averageScore)
	}

	return nil
}
