package ui

//	UI functions the implement the search for an optimal throw target.
//	We cycle through numerous targets around the board, and throw a large number of darts at
//	each target, then report which targets produced the best average score

import (
	boardgeo "DStratMC/board-geometry"
	"DStratMC/simulation"
	target_search "DStratMC/target-search"
	"context"
	"fmt"
	g "github.com/AllenDang/giu"
	"time"
)

//	startSearchForBestThrow begins the search.  We spawn two sub-processes, to keep this, the main process,
//	running to keep the UI responsive.  One subprocess is the actual search, and the other cycles the flag
//	that displays the "searching, please wait" message on and off periodically

func (u *UserInterfaceInstance) startSearchForBestThrow(model simulation.AccuracyModel, numThrows int32) {
	u.searchResultStrings = [10]string{"", "", "", "", "", "", "", "", "", ""}
	u.dartboard.RemoveThrowMarkers()
	u.searchComplete = false
	g.Update()

	//	Start a process to blink the "searching" label on and off
	var blinkContext context.Context
	blinkContext, u.cancelBlinkTimer = context.WithCancel(context.Background())
	go u.cycleBlinkFlag(blinkContext)

	//	Start the actual search process
	var searchContext context.Context
	searchContext, u.cancelSearch = context.WithCancel(context.Background())
	go u.searchProcess(searchContext, model, numThrows)
}

// cycleBlinkFlag is the sub-process that cycles the blinking message on and off
// It is created with a context object to allow receiving a "cancel" message,
// and runs indefinitely until that cancel is received.
func (u *UserInterfaceInstance) cycleBlinkFlag(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			u.searchingBlinkOn = false
			return
		default:
			u.searchingBlinkOn = !u.searchingBlinkOn
			g.Update()
			time.Sleep(500 * time.Millisecond)
		}
	}
}

// searchProcess is the subprocess that runs the actual target search.
func (u *UserInterfaceInstance) searchProcess(ctx context.Context, model simulation.AccuracyModel, numThrows int32) {
	//	Get target iterator and results aggregator
	targetSupplier := target_search.NewTargetSupplier(u.dartboard.GetSquareDimension(), u.dartboard.GetImageMinPoint())
	results := target_search.NewSimResults()
	u.cancelSearchVisible = true
	u.searchComplete = false
	g.Update()

	//	Try each target
	u.loopThroughAllTargets(ctx, model, numThrows, targetSupplier, results)

	if u.searchCancelled {
		u.dartboard.RemoveThrowMarkers()
		u.searchProgressPercent = 0
		u.messageDisplay = "Search cancelled"
	} else {
		//	Get results, sorted from best to worst
		sortedResults := results.GetResultsSortedByHighScore()

		//  Filter results so each plain-language target is named only once
		u.simResultsOneEach = target_search.FilterToOneTargetEach(sortedResults)

		// Messages saying what were the best targets
		u.reportResults()

		//	Draw best target on the board
		bestTargetPosition := u.simResultsOneEach[0].Position
		u.dartboard.SetStdDeviationCirclesCentre(bestTargetPosition)
		u.dartboard.QueueTargetMarker(bestTargetPosition)
	}

	//	Stop the blink timer
	u.cancelSearchVisible = false
	u.cancelBlinkTimer()
	g.Update()
}

// loopThroughAllTargets uses the target supplier iterator to loop through every possible target, and throw
// a large number of darts at each, recording the average score for each
func (u *UserInterfaceInstance) loopThroughAllTargets(ctx context.Context, model simulation.AccuracyModel, numThrows int32, targetSupplier target_search.TargetSupplier, results target_search.SimResults) {
	u.searchProgressPercent = 0
	// Loop through all targets
	targetCount := float64(0)
	howManyTargetsExpected := targetSupplier.ForecastNumTargets()
	for targetSupplier.HasNext() {
		select {
		case <-ctx.Done():
			fmt.Println("Search cancelled")
			u.searchCancelled = true
			return
		default:
			//	Provide visual feedback of what's going on
			targetCount += 1
			u.searchProgressPercent = targetCount / howManyTargetsExpected
			target := targetSupplier.NextTarget()
			// Mark this target on the dartboard
			if u.searchShowEachTarget {
				u.dartboard.QueueTargetMarker(target)
				g.Update()
			}
			// Do a large number of throws at this target
			averageScore, err := u.multipleThrowsAtTarget(target, model, numThrows)
			if err != nil {
				fmt.Printf("Error throwing at target %v: %v", target, err)
				continue
			}
			//	record result for this target
			results.AddTargetResult(target, averageScore)
		}
	}
}

// reportResults reports the results of the simulation by console messages and by setting the
// ui variables that will be displayed for the best 10 targets
func (u *UserInterfaceInstance) reportResults() {
	fmt.Println("First 10 best choices, from best down:")
	for i := 0; i < 10; i++ {
		_, score, description := boardgeo.DescribeBoardPoint(u.simResultsOneEach[i].Position)
		fmt.Printf("   %s (theoretical score %d, average %g)\n", description, score, u.simResultsOneEach[i].Score)
		u.searchResultStrings[i] = fmt.Sprintf("%s (%.2f)", description, u.simResultsOneEach[i].Score)
	}
	//	Setting the "search complete" flag allows the result labels to be displayed in the next UI loop pass
	u.searchComplete = true
}
