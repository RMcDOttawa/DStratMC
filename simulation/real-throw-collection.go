package simulation

import (
	boardgeo "DStratMC/board-geometry"
	"encoding/json"
	"fmt"
	"math"
)

type RealThrowCollection interface {
	AddHit(target boardgeo.BoardPosition, hit boardgeo.BoardPosition)
	GetNumThrows() int
	GetStdDevString() string
	IsStdDevAvailable() bool
	CalcStdDevOfThrows() float64
	GetJsonData() string
	LoadStoredJsonData(content []byte)
}

//	A "real throw collection" is a collection of real throws that have been made by a player.
//  It is a two-level list.  The top level is a list of targets chosen.  Each chosen target
//	has a list of throws that were made at that target.

type hitsList []boardgeo.BoardPosition

type RealThrowCollectionInstance struct {
	targetsList  map[boardgeo.BoardPosition]hitsList
	dataChanged  bool
	cachedStdDev float64
}

func NewRealThrowCollectionInstance() RealThrowCollection {
	return &RealThrowCollectionInstance{
		targetsList: make(map[boardgeo.BoardPosition]hitsList),
		dataChanged: true,
	}
}

func (r *RealThrowCollectionInstance) AddHit(target boardgeo.BoardPosition, hit boardgeo.BoardPosition) {
	//fmt.Printf("Add hit %v at target %v\n", hit, target)
	// If the target is not in the map, add it
	if _, ok := r.targetsList[target]; !ok {
		r.targetsList[target] = hitsList{}
	}
	// Add the hit to the list of hits for that target
	r.targetsList[target] = append(r.targetsList[target], hit)
	r.dataChanged = true
}

// GetNumThrows returns the total number of throws that have been made at all targets
func (r *RealThrowCollectionInstance) GetNumThrows() int {
	countThrows := 0
	for _, hits := range r.targetsList {
		countThrows += len(hits)
	}
	return countThrows
}

// GetStdDevString returns a string representation of the standard deviation of the throws
// if there are enough data points to calculate it, or the string "N/A" if not
func (r *RealThrowCollectionInstance) GetStdDevString() string {
	if r.IsStdDevAvailable() {
		stdDev := r.CalcStdDevOfThrows()
		return fmt.Sprintf("%.3f", stdDev)
	}
	return "N/A"
}

// IsStdDevAvailable returns true if there are enough data points to calculate the standard deviation
func (r *RealThrowCollectionInstance) IsStdDevAvailable() bool {
	//TODO implement me
	return r.GetNumThrows() >= 3
}

func (r *RealThrowCollectionInstance) CalcStdDevOfThrows() float64 {
	if r.dataChanged {
		r.dataChanged = false
		//	We'll calculate a single overall standard deviation for all the throws
		//  at all the targets.

		//	Convert every throw to a floating point value representing the error
		//  in the polar radius.  (We don't care about the angle)
		//  By letting these values remain positive or negative, we will automatically
		//  get values that are normal around a mean of zero
		errorValues := make([]float64, 0, r.GetNumThrows())
		for target, hits := range r.targetsList {
			for _, hit := range hits {
				radiusError := hit.Radius - target.Radius
				// Multiply radius error by 2 since our code uses diameter error
				errorValues = append(errorValues, radiusError*2)
			}
		}

		//fmt.Println(" errors list", errorValues)
		// Calculate mean
		var sum float64 = 0
		for _, value := range errorValues {
			sum += value
		}
		meanValue := sum / float64(len(errorValues))
		//fmt.Println("Error Values", errorValues)
		//fmt.Println("Mean", meanValue)

		// Mean is fixed at zero
		//const meanValue float64 = 0

		// Calculate the sum of the squares of the differences between each throw and the mean
		var sumOfSquares float64 = 0
		for _, value := range errorValues {
			sumOfSquares += (value - meanValue) * (value - meanValue)
		}
		// Divide by the number of throws
		// Take the square root of the result
		// This is the standard deviation of the throws.
		sigma2 := math.Sqrt(sumOfSquares / float64(len(errorValues)))
		// But we have been measuring where darts hi in a small sample - we're going to assume
		// this is the 95% circle, i.e. 2 sigma.  So divide by 2 to get the standard deviation
		r.cachedStdDev = sigma2 / 2
		return r.cachedStdDev
	}
	return r.cachedStdDev
}

// GetJsonData returns the data in the collection as a JSON stream, suitable for a file
func (r *RealThrowCollectionInstance) GetJsonData() string {
	// We just store the data, not the cache
	jsonString, err := json.MarshalIndent(r.targetsList, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling JSON", err)
		return ""
	}
	return string(jsonString)
}

func (r *RealThrowCollectionInstance) LoadStoredJsonData(content []byte) {
	//fmt.Println("LoadStoredJsonData. Loaded file content: ", string(content))
	var decodedMap map[boardgeo.BoardPosition]hitsList
	err := json.Unmarshal(content, &decodedMap)
	if err != nil {
		fmt.Println("Error unmarshalling JSON", err)
		return
	}
	r.targetsList = decodedMap
	r.dataChanged = true
}
