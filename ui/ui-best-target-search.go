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
	"runtime"
	"sync"
	"time"
)

const use_legacy_single_threaded_search = false

// Empirically, more threads than this isn't faster because of coordination costs.
// In a future version, we will try to make the threads more independent by accepting large
// blocks of targets and returning large blocks of results

const num_search_workers = 4

//	startSearchForBestThrow begins the search.  We spawn two sub-processes, to keep this, the mac-binary process,
//	running to keep the UI responsive.  One subprocess is the actual search, and the other cycles the flag
//	that displays the "searching, please wait" message on and off periodically

func (u *UserInterfaceInstance) startSearchForBestThrow(model simulation.AccuracyModel, numThrows int32) {
	u.searchResultStrings = [10]string{"", "", "", "", "", "", "", "", "", ""}
	u.dartboard.RemoveThrowMarkers()
	u.searchComplete = false
	timeBeforeSearch := time.Now()

	g.Update()
	//	Start a process to blink the "searching" label on and off
	var blinkContext context.Context
	blinkContext, u.cancelBlinkTimer = context.WithCancel(context.Background())
	go u.cycleBlinkFlag(blinkContext, func() {
		timeAfterSearch := time.Now()
		fmt.Printf("Search took %v\n", timeAfterSearch.Sub(timeBeforeSearch))
	})

	//	Start the actual search process
	var searchContext context.Context
	searchContext, u.cancelSearch = context.WithCancel(context.Background())
	go u.searchProcess(searchContext, model, numThrows)

}

// cycleBlinkFlag is the sub-process that cycles the blinking message on and off
// It is created with a context object to allow receiving a "cancel" message,
// and runs indefinitely until that cancel is received.
func (u *UserInterfaceInstance) cycleBlinkFlag(ctx context.Context, doneCallback func()) {
	for {
		select {
		case <-ctx.Done():
			u.searchingBlinkOn = false
			doneCallback()
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
	u.messageDisplay = ""
	u.scoreDisplay = ""
	g.Update()

	if use_legacy_single_threaded_search {
		//	Try each target
		u.loopThroughAllTargets(ctx, model, numThrows, targetSupplier, results)
	} else {
		u.multiThreadedSearch(ctx, model, numThrows, targetSupplier, results, num_search_workers)
	}

	if u.searchCancelled {
		u.dartboard.RemoveThrowMarkers()
		u.searchProgressPercent = 0
		u.messageDisplay = "Search cancelled"
	} else {
		//	Get results, sorted from best to worst
		if results.GetNumResults() > 0 {

			sortedResults := results.GetResultsSortedByHighScore()

			//  Filter results so each plain-language target is named only once
			u.simResultsOneEach = target_search.FilterToOneTargetEach(sortedResults)

			// Messages saying what were the best targets
			u.reportResults()

			//	Draw best target on the board
			bestTargetPosition := u.simResultsOneEach[0].Position
			u.searchResultsRadio = 0
			u.dartboard.SetStdDeviationCirclesCentre(bestTargetPosition)
			u.dartboard.QueueTargetMarker(bestTargetPosition)
		}
	}

	//	Stop the blink timer
	u.cancelSearchVisible = false
	u.cancelBlinkTimer()
	g.Update()
}

// loopThroughAllTargets uses the target supplier iterator to loop through every possible target, and throw
// a large number of darts at each, recording the average score for each
func (u *UserInterfaceInstance) loopThroughAllTargets(ctx context.Context,
	model simulation.AccuracyModel,
	numThrows int32,
	targetSupplier target_search.TargetSupplier,
	results target_search.SimResults) {
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
			u.searchProgressPercent = targetCount / float64(howManyTargetsExpected)
			target := targetSupplier.NextTarget()
			// Mark this target on the dartboard
			if u.searchShowEachTarget {
				u.dartboard.QueueTargetMarker(target)
			}
			g.Update()
			// Do a large number of throws at this target
			averageScore, err := u.multipleThrowsAtTarget(target, model, numThrows)
			if err != nil {
				fmt.Printf("Error throwing at target %v: %v", target, err)
				continue
			}
			//	record result for this target
			results.AddTargetResult(target_search.TargetResult{Position: target, Score: averageScore})
		}
	}
	// Clear progress bar
	u.searchProgressPercent = 0
	g.Update()
}

// reportResults reports the results of the simulation by console messages and by setting the
// ui variables that will be displayed for the best 10 targets
func (u *UserInterfaceInstance) reportResults() {
	for i := 0; i < 10; i++ {
		_, score, description := boardgeo.DescribeBoardPoint(u.simResultsOneEach[i].Position)
		fmt.Printf("   %s (theoretical score %d, average %g)\n", description, score, u.simResultsOneEach[i].Score)
		u.searchResultStrings[i] = fmt.Sprintf("%s (%.2f)", description, u.simResultsOneEach[i].Score)
	}
	//	Setting the "search complete" flag allows the result labels to be displayed in the next UI loop pass
	u.searchComplete = true
}

func (u *UserInterfaceInstance) multiThreadedSearch(
	ctx context.Context,
	model simulation.AccuracyModel,
	throws int32,
	supplier target_search.TargetSupplier,
	results target_search.SimResults,
	numWorkers uint16) {
	fmt.Println("multiThreadedSearch starting")
	fmt.Println("  num workers: ", numWorkers)
	fmt.Println("  CPUs: ", runtime.NumCPU())
	gm := runtime.GOMAXPROCS(int(numWorkers))
	fmt.Println("  GOMAXPROCS: ", gm)
	//fmt.Println("  throws: ", throws)
	//fmt.Println("  supplier: ", supplier)
	//fmt.Println("  results: ", results)

	// Done counters
	var wg sync.WaitGroup

	// Channels to send targets and receive results
	channelCapacity := supplier.ForecastNumTargets()
	targetsChannel := make(chan boardgeo.BoardPosition, channelCapacity)
	resultsChannel := make(chan target_search.TargetResult, channelCapacity)

	//	Start worker threads
	for i := uint16(0); i < numWorkers; i++ {
		//fmt.Printf("  Starting worker number %d\n", i+1)
		wg.Add(1)
		go u.workerThread(i, ctx, model, throws, targetsChannel, resultsChannel, &wg)
	}

	// Fill the targets channel
	//fmt.Println("Filling targets channel")
	for supplier.HasNext() {
		target := supplier.NextTarget()
		targetsChannel <- target
	}
	close(targetsChannel)

	// Wait for all worker threads to finish
	go func() {
		//fmt.Println("Starting mini-task to wait for threads to finish")
		wg.Wait()
		//fmt.Println("All worker threads finished, closing results channel")
		close(resultsChannel)
	}()

	// Read the results as they come in
	//fmt.Println("Reading results")
	numResults := 0
	totalResultsDenominator := float64(channelCapacity)
reader:
	for {
		select {
		case <-ctx.Done():
			//fmt.Println("Search cancelled")
			u.searchCancelled = true
			break reader
		case result, ok := <-resultsChannel:
			if ok {
				//fmt.Printf("  Got result %v\n", result)
				results.AddTargetResult(result)
				numResults += 1
				u.searchProgressPercent = float64(numResults) / totalResultsDenominator
			} else {
				//fmt.Println("Results channel closed")
				break reader
			}
		}
	}
	u.cancelBlinkTimer()
	// Clear progress bar
	u.searchProgressPercent = 0
	g.Update()
	//fmt.Println("multiThreadedSearch ending")
}

func (u *UserInterfaceInstance) workerThread(
	threadNumber uint16,
	ctx context.Context,
	model simulation.AccuracyModel,
	throws int32,
	targetsChannel chan boardgeo.BoardPosition,
	resultsChannel chan target_search.TargetResult,
	wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Worker thread %d starting\n", threadNumber)

	//fmt.Printf("Worker %d requesting work from targets channel\n", threadNumber)
	for {
		select {
		case <-ctx.Done():
			//fmt.Printf("Worker %d thread ending", threadNumber)
			return
		case target, ok := <-targetsChannel:
			if ok {
				//fmt.Printf("  Worker %d received target: %v\n", threadNumber, target)
				averageScore, err := u.multipleThrowsAtTarget(target, model, throws)
				if err != nil {
					panic(err)
				} else {
					resultsChannel <- target_search.TargetResult{Position: target, Score: averageScore}
				}
			} else {
				//fmt.Printf("  Worker %d, supply channel closed\n", threadNumber)
				return
			}
		}
	}
}
