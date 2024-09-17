package target_search

import (
	boardgeo "DStratMC/board-geometry"
	"slices"
	"sort"
)

// SimResults stores the results of a simulation run - each target position tried, and its average score

type TargetResult struct {
	Position boardgeo.BoardPosition
	Score    float64
}

type SimResults interface {
	AddTargetResult(result TargetResult)
	GetResultsSortedByHighScore() []OneResult
	GetResultsSlice() []OneResult
	GetNumResults() uint32
}

// SimResultsInstance is data for the instance of the SimResults object
type SimResultsInstance struct {
	resultsMap map[boardgeo.BoardPosition]float64
}

// NewSimResults creates a new SimResults object
func NewSimResults() SimResults {
	results := &SimResultsInstance{
		resultsMap: make(map[boardgeo.BoardPosition]float64, 4000),
	}
	//fmt.Println("NewSimResults  returns", results)
	return results
}

// OneResult is a single result - a position tried and the average score at that position
type OneResult struct {
	Position boardgeo.BoardPosition
	Score    float64
}

// GetResultsSlice returns the results as a slice of OneResult objects, in no particular order
func (s SimResultsInstance) GetResultsSlice() []OneResult {
	slice := make([]OneResult, 0, len(s.resultsMap))
	for pos, score := range s.resultsMap {
		slice = append(slice, OneResult{Position: pos, Score: score})
	}
	return slice
}

func (s SimResultsInstance) GetNumResults() uint32 {
	return uint32(len(s.resultsMap))
}

// GetResultsSortedByHighScore returns the positions sorted from the highest average score to the lowest
func (s SimResultsInstance) GetResultsSortedByHighScore() []OneResult {
	// Convert map to slice
	slice := s.GetResultsSlice()
	// Sort the slice by descending order of score
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].Score > slice[j].Score
	})
	return slice
}

// AddTargetResult adds a target position and its average score to the results list
func (s SimResultsInstance) AddTargetResult(result TargetResult) {
	s.resultsMap[result.Position] = result.Score
}

// FilterToOneTargetEach returns a slice of OneResult objects, with only one result for each target position
// "target position" in the sense of board segment (e.g., "treble 20", "double 2"), not precise coordinates
func FilterToOneTargetEach(results []OneResult) []OneResult {
	targetsUsed := make([]string, 0, len(results))
	outputSlice := make([]OneResult, 0, len(results))
	for i := 0; i < len(results); i++ {
		_, _, descriptionString := boardgeo.DescribeBoardPoint(results[i].Position)
		if slices.Contains(targetsUsed, descriptionString) {
			// Already have reported on this one, no more
		} else {
			outputSlice = append(outputSlice, results[i])
			targetsUsed = append(targetsUsed, descriptionString)
		}
	}
	return outputSlice
}
