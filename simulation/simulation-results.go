package simulation

import (
	boardgeo "DStratMC/board-geometry"
	"slices"
	"sort"
)

type SimResults interface {
	AddTargetResult(position boardgeo.BoardPosition, score float64)
	GetResultsSortedByHighScore() []OneResult
	GetResultsSlice() []OneResult
}

type SimResultsInstance struct {
	resultsMap map[boardgeo.BoardPosition]float64
}

type OneResult struct {
	Position boardgeo.BoardPosition
	Score    float64
}

//func (s SimResultsInstance) GetAverageScore(p boardgeo.BoardPosition) float64 {
//	sum := s.SumMap[p]
//	count := s.CountMap[p]
//	if count == 0 {
//		return 0
//	}
//	return float64(sum) / float64(count)
//}
//
//func (s SimResultsInstance) GetPositions() []boardgeo.BoardPosition {
//	positions := make([]boardgeo.BoardPosition, 0, len(s.SumMap))
//	for pos := range s.SumMap {
//		positions = append(positions, pos)
//	}
//	return positions
//}

// Return the positions sorted from the highest average score to the lowest
func (s SimResultsInstance) GetResultsSortedByHighScore() []OneResult {
	// Convert map to slice
	slice := s.GetResultsSlice()
	// Sort the slice by descending order of score
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].Score > slice[j].Score
	})
	return slice
}

func (s SimResultsInstance) GetResultsSlice() []OneResult {
	slice := make([]OneResult, 0, len(s.resultsMap))
	for pos, score := range s.resultsMap {
		slice = append(slice, OneResult{Position: pos, Score: score})
	}
	return slice
}

func NewSimResults() SimResults {
	results := &SimResultsInstance{
		resultsMap: make(map[boardgeo.BoardPosition]float64, 4000),
	}
	//fmt.Println("NewSimResults  returns", results)
	return results
}

func (s SimResultsInstance) AddTargetResult(position boardgeo.BoardPosition, score float64) {
	//fmt.Println("AddTargetResult STUB", position, score)
	s.resultsMap[position] = score
}

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
