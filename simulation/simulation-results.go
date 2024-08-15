package simulation

import (
	boardgeo "DStratMC/board-geometry"
	"fmt"
)

type SimResults interface {
	AddTargetResult(position boardgeo.BoardPosition, score float64)
	//GetPositions() []boardgeo.BoardPosition
	//GetAverageScore(p boardgeo.BoardPosition) float64
	//GetPositionsSortedByHighScore() []boardgeo.BoardPosition
}

type SimResultsInstance struct {
	//SumMap   map[boardgeo.BoardPosition]int
	//CountMap map[boardgeo.BoardPosition]int
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
//
//// Return the positions sorted from the highest average score to the lowest
//func (s SimResultsInstance) GetPositionsSortedByHighScore() []boardgeo.BoardPosition {
//	positions := s.GetPositions()
//	//	Sort the positions by average score
//	//	We'll use a simple bubble sort, because it's easy and we're not expecting a lot of positions
//	for i := 0; i < len(positions); i++ {
//		for j := i + 1; j < len(positions); j++ {
//			if s.GetAverageScore(positions[i]) < s.GetAverageScore(positions[j]) {
//				positions[i], positions[j] = positions[j], positions[i]
//			}
//		}
//	}
//	return positions
//}

func NewSimResults() SimResults {
	results := &SimResultsInstance{
		//SumMap:   make(map[boardgeo.BoardPosition]int, 1000),
		//CountMap: make(map[boardgeo.BoardPosition]int, 1000),
	}
	fmt.Println("NewSimResults STUB returns", results)
	return results
}

func (s SimResultsInstance) AddTargetResult(position boardgeo.BoardPosition, score float64) {
	fmt.Println("AddTargetResult STUB", position, score)
	//s.SumMap[position] += score
	//s.CountMap[position]++
}
