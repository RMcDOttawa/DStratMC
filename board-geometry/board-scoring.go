package boardgeo

//	These are the data and functions used to determine the scoring value of a dart throw

import "math"

// Boundaries between the single, double, and treble areas on the board, measured
// in normalized radius units (0.0 to 1.0) from the centre of the board
var multiplierBoundaries = [...]float64{
	0.0, // Exact centre of board
	innerBullRadiusNormalized,
	outerBullRadiusNormalized,
	insideTrebleRadiusNormalized,
	outsideTrebleRadiusNormalized,
	insideDoubleRadiusNormalized,
	outsideDoubleRadiusNormalized,
	scoringAreaRadiusNormalized,
}

var multiplierList = [...]int{
	2, //	Red bull is a double
	1, //	Green bull is a single
	1, //	Single area between bull and treble
	3, //	Treble ring
	1, //	Single area between treble and double
	2, //	Double ring
	0, //	Outside board
}

// Point values of slices, clockwise from the top
var segmentPointValues = [...]int{
	20, 1, 18, 4, 13, 6, 10, 15, 2, 17,
	3, 19, 7, 16, 8, 11, 14, 9, 12, 5}

//	determineSinglePointValue determines the single point value of a spot on the board,
//	not taking double and triple rings into account.

//	Use the angle on the board to determine which point wedge has been hit (ignoring double,triple)
//	Angles come in as numbers between -180 and +180.  We'll convert them to 0 to 360.
//	0 is straight up, then increasing numbers are clockwise rotation around the board.
//	Wedges are 360/20 = 18 degrees wide, and are offset by 1/2 that, 9 degrees.

func determineSinglePointValue(radius, degrees float64) int {
	//	Special cases: the inner and outer bulls - rotation doesn't matter
	if radius < outerBullRadiusNormalized {
		return 25
	}

	//	Convert +/- 180 range to 0-360 range
	if degrees < 0 {
		degrees = degrees + 360.0
	}

	//	Shift by 9 degrees (half a segment) to de-center the slices. So the beginning of the "1" slice is
	//	18 degrees, the 18 slice is 36 degrees, etc.
	degrees += (360.0 / 20) / 2

	//	Convert that to an integer index, where index 0 is the 20 slice, [1] is the 1 slice, [2] is the 18, etc.
	//	(20 will be messed up, we'll fix in a moment)
	sliceIndex := int(math.Floor(degrees / (360.0 / 20)))

	//	If we're on the left side of the 20 slice, this will have generated an index of 20, which we
	//	need to convert to 0
	if sliceIndex == 20 {
		sliceIndex = 0
	}
	return segmentPointValues[sliceIndex]
}

// determineMultiplier determines the multiplier for a spot on the board, based on the radius
// i.e. whether it's a single, double, treble, or outside the board
func determineMultiplier(radius float64) int {
	foundMultiplierIndex := -1
	for i := 0; i < len(multiplierBoundaries); i++ {
		thisBoundary := multiplierBoundaries[i]
		if radius < thisBoundary {
			foundMultiplierIndex = i - 1
			break
		}
	}
	if foundMultiplierIndex == -1 {
		return 0
	}
	return multiplierList[foundMultiplierIndex]
}
