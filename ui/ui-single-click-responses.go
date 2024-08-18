package ui

import (
	boardgeo "DStratMC/board-geometry"
	"DStratMC/simulation"
	"fmt"
	g "github.com/AllenDang/giu"
	"strconv"
)

// resultButtonClicked is called when a button is clicked that corresponds to a search result
// It will draw a marker at the precise location of that reported search result
func (u *UserInterfaceInstance) resultButtonClicked(buttonIndex int) {
	if buttonIndex < len(u.simResultsOneEach) {
		//	Set the position where any requested standard deviation circles will be drawn
		targetPosition := u.simResultsOneEach[buttonIndex].Position
		u.dartboard.SetStdDeviationCirclesCentre(targetPosition)
		u.dartboard.QueueTargetMarker(targetPosition)
		g.Update()
	}
}

// oneThrowsAtTarget will throw multiple dart at the target position, and return the result using the given accuracy model
func (u *UserInterfaceInstance) multipleThrowsAtTarget(target boardgeo.BoardPosition, model simulation.AccuracyModel, throws int32) (float64, error) {
	var total float64 = 0.0
	for i := 0; i < int(throws); i++ {
		hit, err := model.GetThrow(target,
			u.dartboard.GetScoringRadiusPixels(),
			u.dartboard.GetSquareDimension(),
			u.dartboard.GetImageMinPoint())
		if err != nil {
			return 0.0, err
		}
		_, score, _ := boardgeo.DescribeBoardPoint(hit)
		total += float64(score)
	}
	average := total / float64(throws)
	return average, nil
}

// oneThrowsAtTarget will throw a single dart at the target position, and return the result using the uniform distribution accuracy model
// The results are added to the running total for calculating statistics, and marked on the board with a hit marker
func (u *UserInterfaceInstance) oneUniformThrow(dartboard Dartboard, position boardgeo.BoardPosition, model simulation.AccuracyModel) {

	//  Draw a marker to record where we clicked
	dartboard.QueueTargetMarker(position)

	//	Draw a circle showing the accuracy radius around the clicked point
	accuracyRadius := model.GetAccuracyRadius()
	dartboard.QueueAccuracyCircle(position, accuracyRadius)

	//	Get a modeled hit within the accuracy
	hit, err := model.GetThrow(position,
		dartboard.GetScoringRadiusPixels(),
		dartboard.GetSquareDimension(),
		dartboard.GetImageMinPoint())
	if err != nil {
		fmt.Printf("Error getting throw %v", err)
		return
	}

	//	Draw the hit within this circle
	dartboard.QueueHitMarker(hit, singleHitMarkerRadius)

	//	Calculate the hit score
	_, score, description := boardgeo.DescribeBoardPoint(hit)
	u.messageDisplay = description
	u.scoreDisplay = strconv.Itoa(score) + " points"
	u.throwCount++
	u.throwTotal += int64(score)
	u.throwAverage = float64(u.throwTotal) / float64(u.throwCount)
	g.Update()

}

// multipleUniformThrows will throw multiple darts at the target position, and return the result using the uniform distribution accuracy model
// The results are added to the running total for calculating statistics, and marked on the board with a hit marker
func (u *UserInterfaceInstance) multipleUniformThrows(dartboard Dartboard, position boardgeo.BoardPosition, model simulation.AccuracyModel) {

	//  Draw a marker to record where we clicked
	dartboard.QueueTargetMarker(position)

	//	Draw a circle showing the accuracy radius around the clicked point
	accuracyRadius := model.GetAccuracyRadius()
	dartboard.QueueAccuracyCircle(position, accuracyRadius)
	dartboard.AllocateHitsSpace(int(u.numThrowsField))

	u.throwCount = 0
	u.throwTotal = 0
	for i := 0; i < int(u.numThrowsField); i++ {
		//	Get a modeled hit within the accuracy
		hit, err := model.GetThrow(position,
			dartboard.GetScoringRadiusPixels(),
			dartboard.GetSquareDimension(),
			dartboard.GetImageMinPoint())
		if err != nil {
			fmt.Printf("Error getting throw %v", err)
			return
		}

		//	Draw the hit within this circle
		dartboard.QueueHitMarker(hit, multipleHitMarkerRadius)

		//	Calculate the hit score
		_, score, _ := boardgeo.DescribeBoardPoint(hit)
		u.throwCount++
		u.throwTotal += int64(score)
		u.throwAverage = float64(u.throwTotal) / float64(u.throwCount)
	}
	g.Update()

}

// oneNormalThrow will throw a single dart at the target position, and return the result using the normal distribution accuracy model
// The results are added to the running total for calculating statistics, and marked on the board with a hit marker
func (u *UserInterfaceInstance) oneNormalThrow(dartboard Dartboard, position boardgeo.BoardPosition, model simulation.AccuracyModel) {

	//  Draw a marker to record where we clicked
	dartboard.QueueTargetMarker(position)

	//	Set the position where any requested standard deviation circles will be drawn
	dartboard.SetStdDeviationCirclesCentre(position)

	//	Get a modeled hit within the accuracy
	hit, err := model.GetThrow(position,
		dartboard.GetScoringRadiusPixels(),
		dartboard.GetSquareDimension(),
		dartboard.GetImageMinPoint())
	if err != nil {
		fmt.Printf("Error getting throw %v", err)
		return
	}

	//	Draw the hit within this circle
	dartboard.QueueHitMarker(hit, singleHitMarkerRadius)

	//	Calculate the hit score
	_, score, description := boardgeo.DescribeBoardPoint(hit)
	u.messageDisplay = description
	u.scoreDisplay = strconv.Itoa(score) + " points"
	u.throwCount++
	u.throwTotal += int64(score)
	u.throwAverage = float64(u.throwTotal) / float64(u.throwCount)
	g.Update()

}

// multipleNormalThrows will throw multiple darts at the target position, and return the result using the normal distribution accuracy model
// The results are added to the running total for calculating statistics, and marked on the board with a hit marker
func (u *UserInterfaceInstance) multipleNormalThrows(dartboard Dartboard, position boardgeo.BoardPosition, model simulation.AccuracyModel) {

	//  Draw a marker to record where we clicked
	dartboard.QueueTargetMarker(position)

	//	Set the position where any requested standard deviation circles will be drawn
	dartboard.SetStdDeviationCirclesCentre(position)

	dartboard.AllocateHitsSpace(int(u.numThrowsField))

	u.throwCount = 0
	u.throwTotal = 0

	for i := 0; i < int(u.numThrowsField); i++ {
		//	Get a modeled hit within the accuracy
		hit, err := model.GetThrow(position,
			dartboard.GetScoringRadiusPixels(),
			dartboard.GetSquareDimension(),
			dartboard.GetImageMinPoint())
		if err != nil {
			fmt.Printf("Error getting throw %v", err)
			return
		}

		//	Draw the hit within this circle
		dartboard.QueueHitMarker(hit, multipleHitMarkerRadius)

		//	Calculate the hit score
		_, score, _ := boardgeo.DescribeBoardPoint(hit)
		u.throwCount++
		u.throwTotal += int64(score)
		u.throwAverage = float64(u.throwTotal) / float64(u.throwCount)
	}
	g.Update()

}
