package ui

//	Custom function that draws the actual dartboard and the various annotations such
//	as target markers that we may place on top of it.
//	The dartboard itself is an embedded image, not drawn.  It is carefully drawn in Adobe
//	Illustrator, which allows us to know the precise size of the various circles, and those
//	sizes are hard-coded here to assist with location calculations

import (
	boardgeo "DStratMC/board-geometry"
	"fmt"
	g "github.com/AllenDang/giu"
	"image"
	"image/color"
)

// Attributes of various annotations that may be drawn on top of the dartboard
const targetCrossAlpha = 230
const targetCrossLength = 20
const targetCrossThickness = 2

const hitMarkerAlpha = 200

const accuracyCircleThickness = 2

var accuracyCircleColour = color.RGBA{R: 100, G: 100, B: 255, A: 192}

// Dartboard models the dartboard as an object to keep control
// of the variables associated with it
type Dartboard interface {
	SetInfo(windowWidget *g.WindowWidget, texture *g.Texture,
		squareDimension float64,
		imageMin image.Point, imageMax image.Point)
	SetClickCallback(callback func(dartboard Dartboard, position boardgeo.BoardPosition))
	DrawFunction()
	dartboardClicked()
	RemoveThrowMarkers()
	QueueTargetMarker(position boardgeo.BoardPosition)
	QueueAccuracyCircle(position boardgeo.BoardPosition, radius float64)
	GetScoringRadiusPixels() float64
	GetImageMinPoint() image.Point
	GetSquareDimension() float64
	QueueHitMarker(hit boardgeo.BoardPosition, markerRadius int)
	AllocateHitsSpace(i int)
	SetDrawRefLines(checkbox bool)
	SetDrawOneSigma(draw bool, radius float64)
	SetDrawTwoSigma(draw bool, radius float64)
	SetDrawThreeSigma(draw bool, radius float64)
	SetStdDeviationCirclesCentre(position boardgeo.BoardPosition)
}

type DartboardInstance struct {
	window          *g.WindowWidget
	texture         *g.Texture
	squareDimension float64
	imageMin        image.Point
	imageMax        image.Point
	clickCallback   func(dartboard Dartboard, position boardgeo.BoardPosition)

	// We have drawn a marker showing where a throw was targeted
	targetDrawn    bool
	targetPosition boardgeo.BoardPosition

	// Circle showing the uniform accuracy radius around a clicked point
	drawAccuracyCircle     bool
	accuracyCircleRadius   float64
	accuracyCirclePosition boardgeo.BoardPosition

	//// Zero or more circles showing the standard deviation radii around a clicked point
	//stdDevCirclePositions   []boardgeo.BoardPosition
	//stdDevCircleMultipliers []float64
	//stdDevCircleRadii       []float64

	// Slice of zero or more hits resulting from modeled throw
	hitPositions    []boardgeo.BoardPosition
	hitMarkerRadius int

	// Draw the testing crosshair?
	drawReferenceLines bool

	//	Draw reference circles for 1, 2, and 3 standard deviations?
	stdDeviationCirclesCentre boardgeo.BoardPosition
	stdDevClicked             bool

	drawOneStdDeviation bool
	drawOneStdRadius    float64

	drawTwoStdDeviation bool
	drawTwoStdRadius    float64

	drawThreeStdDeviation bool
	drawThreeStdRadius    float64
}

// NewDartboard creates an instance of the dartboard object
func NewDartboard() Dartboard {
	instance := &DartboardInstance{
		clickCallback:      nil,
		targetDrawn:        false,
		drawAccuracyCircle: false,
		hitPositions:       make([]boardgeo.BoardPosition, 0, throwsAtOneTarget),
	}
	return instance
}

// SetInfo accepts and stores key size and dimension info for the dartboard
func (d *DartboardInstance) SetInfo(windowWidget *g.WindowWidget, texture *g.Texture, squareDimension float64, imageMin image.Point, imageMax image.Point) {
	d.window = windowWidget
	d.texture = texture
	d.squareDimension = squareDimension
	d.imageMin = imageMin
	d.imageMax = imageMax
}

// SetClickCallback sets the function that is called back when the mouse is clicked inside the dartboard scoring area
// The callback function will be called with the clicked mouse position
func (d *DartboardInstance) SetClickCallback(callback func(dartboard Dartboard, position boardgeo.BoardPosition)) {
	d.clickCallback = callback
}

// AllocateHitsSpace can be called to pre-allocate space in the slice that records hits,
// once that quantity is known.  (It's only an estimate but helps efficiency to preallocate)
func (d *DartboardInstance) AllocateHitsSpace(numHits int) {
	d.hitPositions = make([]boardgeo.BoardPosition, 0, numHits)
}

// SetStdDeviationCirclesCentre records the centre position where we will draw the optional
// standard deviation reference circles
func (d *DartboardInstance) SetStdDeviationCirclesCentre(position boardgeo.BoardPosition) {
	d.stdDeviationCirclesCentre = position
	d.stdDevClicked = true
}

// SetDrawOneSigma records the information to draw a circle representing 1 standard deviation from centre
func (d *DartboardInstance) SetDrawOneSigma(draw bool, radius float64) {
	fmt.Printf("SetDrawOneSigma(%t, %g)\n", draw, radius)
	d.drawOneStdDeviation = draw
	d.drawOneStdRadius = radius
}

// SetDrawTwoSigma records the information to draw a circle representing 2 standard deviations from centre
func (d *DartboardInstance) SetDrawTwoSigma(draw bool, radius float64) {
	fmt.Printf("SetDrawTwoSigma(%t, %g)\n", draw, radius)
	d.drawTwoStdDeviation = draw
	d.drawTwoStdRadius = radius
}

// SetDrawThreeSigma records the information to draw a circle representing 3 standard deviations from centre
func (d *DartboardInstance) SetDrawThreeSigma(draw bool, radius float64) {
	fmt.Printf("SetDrawThreeSigma(%t, %g)\n", draw, radius)
	d.drawThreeStdDeviation = draw
	d.drawThreeStdRadius = radius
}

// SetDrawRefLines requests whether crosshair reference lines should be drawn through the board centre
func (d *DartboardInstance) SetDrawRefLines(checkbox bool) {
	d.drawReferenceLines = checkbox
}

// GetSquareDimension retrieves the square dimension, in pixels, of the square subframe containing the board
func (d *DartboardInstance) GetSquareDimension() float64 {
	return d.squareDimension
}

// GetImageMinPoint returns the x,y point of the origin of the dartboard square in the containing window
func (d *DartboardInstance) GetImageMinPoint() image.Point {
	return d.imageMin
}

//		GetScoringRadiusPixels returns the radius, in pixels, of the largest circle on the dartboard that is
//	 in the scoring area (i.e. inside the outer radius of the Double ring)
func (d *DartboardInstance) GetScoringRadiusPixels() float64 {
	radius := d.squareDimension * boardgeo.ScoringAreaFraction / 2
	return radius
}

// RemoveThrowMarkers resets the visual appearance of the dartboard by removing any previously-drawn
// overlays such as throw markers, hits, and circles
func (d *DartboardInstance) RemoveThrowMarkers() {
	d.targetDrawn = false
	d.drawAccuracyCircle = false
	d.hitPositions = make([]boardgeo.BoardPosition, 0, throwsAtOneTarget)
	d.stdDevClicked = false
}

//	DrawFunction is the actual drawing function for the dartboard. It draws the underlying image,
//	places an invisible button on top of it to detect clicks, and draws any annotations such as the
//	reference lines, target markers, etc.

//	Note that because of how the GIU functions, drawing is not persistent.  This function is called 30 times
//	per second, any anything drawn is only there for 1/30 second.
//	So the board and annotations are remembered and re-drawn every time.

func (d *DartboardInstance) DrawFunction() {
	//	The way GIU works, this function can be called before we are ready to draw something meaningful.
	//	We detect this with squareDimension == 0 or positon coordinates == 0 and return
	if d.squareDimension == 0 {
		//fmt.Println("Square dimension 0, returning")
		return
	}
	if d.imageMin.X < 0 || d.imageMin.Y < 0 || d.imageMax.X < 0 || d.imageMax.Y < 0 {
		//fmt.Println("imageMin or Max 0, returning")
		return
	}

	canvas := g.GetCanvas()

	//	Position an invisible button on top of this image to detect clicks
	//	Remember and then restore drawing cursor so image comes out on top of this
	savedCsp := g.GetCursorScreenPos()
	g.SetCursorScreenPos(d.imageMin)
	g.InvisibleButton().Size(float32(d.squareDimension), float32(d.squareDimension)).
		OnClick(d.dartboardClicked).
		Build()
	g.SetCursorScreenPos(savedCsp)

	// Display dartboard image
	canvas.AddImage(d.texture, d.imageMin, d.imageMax)

	if d.drawReferenceLines {
		d.drawReferenceLinesOnDartboard(canvas)
	}

	//	If we have a target position to draw, do that
	if d.targetDrawn {
		d.DrawQueuedTargetMarker(canvas)
	}

	if d.drawAccuracyCircle {
		d.drawQueuedAccuracyCircle(canvas)
	}
	d.drawStdDeviationCircles(canvas)

	d.drawQueuedHitMarkers()
}

// drawReferenceLinesOnDartboard  draws a semitransparent circle and crosshair on the centre
// of the dartboard, to assist with testing coordinates translation
func (d *DartboardInstance) drawReferenceLinesOnDartboard(canvas *g.Canvas) {
	xCentre := (d.imageMin.X + d.imageMax.X) / 2
	yCentre := (d.imageMin.Y + d.imageMax.Y) / 2

	//	Early in development cycle, we included a reference circle
	//testCirclePosition := image.Pt(xCentre, yCentre)
	//testCircleColour := color.RGBA{R: 0, G: 0, B: 255, A: 128}
	//radius := float32(d.squareDimension / 8.0)
	//canvas.AddCircle(testCirclePosition, radius, testCircleColour, 0, 1)

	//	And add centred vertical and horizontal lines to help calibrate angle measurement
	crossHairColour := color.RGBA{R: 150, G: 150, B: 150, A: 255}

	verticalFrom := image.Pt(xCentre, yCentre-int(d.squareDimension/2))
	verticalTo := image.Pt(xCentre, yCentre+int(d.squareDimension/2))
	canvas.AddLine(verticalFrom, verticalTo, crossHairColour, 1)

	horizontalFrom := image.Pt(xCentre-int(d.squareDimension/2), yCentre)
	horizontalTo := image.Pt(xCentre+int(d.squareDimension/2), yCentre)
	canvas.AddLine(horizontalFrom, horizontalTo, crossHairColour, 1)
}

// dartboardClicked is the callback function for the invisible button that covers the dartboard image
// Here we determine where the mouse was and pass the click through to the provided callback function
func (d *DartboardInstance) dartboardClicked() {
	//fmt.Println("dartboard clicked")
	if d.clickCallback == nil {
		panic("  No callback function")
	} else {
		position := boardgeo.CreateBoardPositionFromXY(g.GetMousePos(), d.squareDimension,
			d.imageMin)
		d.clickCallback(d, position)
	}
}

// QueueTargetMarker records a target marker to be drawn on the next time through the ui loop
func (d *DartboardInstance) QueueTargetMarker(position boardgeo.BoardPosition) {
	//fmt.Printf("QueueTargetMarker at %v\n", position)
	d.targetDrawn = true
	d.targetPosition = position
}

// DrawQueuedTargetMarker draws the target marker that has been recorded
func (d *DartboardInstance) DrawQueuedTargetMarker(canvas *g.Canvas) {

	//	Get the pixel coordinates of this point
	xCentre, yCentre := boardgeo.GetDrawingXY(d.targetPosition)
	xCentre += d.imageMin.X
	yCentre += d.imageMin.Y

	//	Get a contrasting colour that will be visible on this board section
	//	(This isn't perfect - it's the contrasting colour only of the centre of the marker.
	//	If the marker crosses into other coloured backgrounds, we are not breaking it into
	//	segments with different contrasting colours.  That would be a neat future project.)
	red, green, blue := contrastingColourForPosition(d.targetPosition)
	colour := color.RGBA{R: red, G: green, B: blue, A: targetCrossAlpha}

	//	Draw an upright cross at this point

	verticalFrom := image.Pt(xCentre, yCentre-targetCrossLength/2)
	verticalTo := image.Pt(xCentre, yCentre+targetCrossLength/2)
	canvas.AddLine(verticalFrom, verticalTo, colour, targetCrossThickness)

	horizontalFrom := image.Pt(xCentre-targetCrossLength/2, yCentre)
	horizontalTo := image.Pt(xCentre+targetCrossLength/2, yCentre)
	canvas.AddLine(horizontalFrom, horizontalTo, colour, targetCrossThickness)

	//	Refactoring. I have no idea why the following line ended up in the code, and am 99.999% sure it
	//	was an error. Commenting out for now, and will delete in a while when I'm sure.
	//boardgeo.DescribeBoardPoint(d.targetPosition)
}

// Get RGB values for a colour that contrasts with the colour under the given board position
func contrastingColourForPosition(position boardgeo.BoardPosition) (uint8, uint8, uint8) {
	segment, score, _ := boardgeo.DescribeBoardPoint(position)
	underlyingColour := boardgeo.GetColourForSegment(segment, score)
	red, green, blue := boardgeo.GetContrastingColour(underlyingColour)
	return red, green, blue
}

// QueueAccuracyCircle records the coordinates of a circle that will be drawn on the next UI loop pass
// marking the defined uniform accuracy circle
func (d *DartboardInstance) QueueAccuracyCircle(position boardgeo.BoardPosition, radius float64) {
	d.accuracyCircleRadius = radius
	d.accuracyCirclePosition = position
	d.drawAccuracyCircle = true
}

// drawStdDeviationCircles draws the standard deviation circles that have been recorded
func (d *DartboardInstance) drawQueuedAccuracyCircle(canvas *g.Canvas) {
	xCentre, yCentre := boardgeo.GetDrawingXY(d.accuracyCirclePosition)
	accuracyCirclePosition := image.Pt(xCentre+d.imageMin.X, yCentre+d.imageMin.Y)
	drawRadius := d.accuracyCircleRadius * d.squareDimension * boardgeo.ScoringAreaFraction / 2
	canvas.AddCircle(accuracyCirclePosition, float32(drawRadius), accuracyCircleColour, 0, accuracyCircleThickness)
}

// drawStdDeviationCircles draws the standard deviation circles that have been recorded
func (d *DartboardInstance) drawStdDeviationCircles(canvas *g.Canvas) {
	if d.stdDevClicked {
		//fmt.Printf("std dev clicked, d = %#v\n", d)
		if d.drawOneStdDeviation {
			d.drawStdDeviationCircle(canvas, 1, "68", d.drawOneStdRadius)
		}
		if d.drawTwoStdDeviation {
			d.drawStdDeviationCircle(canvas, 2, "95", d.drawTwoStdRadius)
		}
		if d.drawThreeStdDeviation {
			d.drawStdDeviationCircle(canvas, 3, "99.7", d.drawThreeStdRadius)
		}
	}
}

// drawStdDeviationCircle draws a circle representing a standard deviation from the centre
//
//	We draw a circle at the stored centre position, with the given radius, and label it with the
//	multiplier and percentage of the normal distribution that it represents
func (d *DartboardInstance) drawStdDeviationCircle(canvas *g.Canvas, multiplier float64, percentageString string, radius float64) {
	// Draw the circle for this standard deviation reference
	xCentre, yCentre := boardgeo.GetDrawingXY(d.stdDeviationCirclesCentre)
	circlePosition := image.Pt(xCentre+d.imageMin.X, yCentre+d.imageMin.Y)
	drawRadius := radius * d.squareDimension * boardgeo.ScoringAreaFraction / 2
	canvas.AddCircle(circlePosition, float32(drawRadius), accuracyCircleColour, 0, accuracyCircleThickness)

	//	Label the top of the circle with the multiplier
	circleLabel := fmt.Sprintf("%2g std (%s%%)", multiplier, percentageString)
	labelWidth, labelHeight := g.CalcTextSize(circleLabel)

	labelPosition := circlePosition
	labelPosition.X -= int(labelWidth / 2)
	labelPosition.Y -= int(drawRadius + float64(labelHeight))
	canvas.AddText(labelPosition, accuracyCircleColour, circleLabel)
}

// QueueHitMarker records the position of a throw hit in a list (there may be many). The queued markers will be drawn
// on the next UI loop pass
func (d *DartboardInstance) QueueHitMarker(hit boardgeo.BoardPosition, markerRadius int) {
	d.hitPositions = append(d.hitPositions, hit)
	d.hitMarkerRadius = markerRadius
}

// drawQueuedHitMarkers draws all the hit markers that have been queued
func (d *DartboardInstance) drawQueuedHitMarkers() {
	canvas := g.GetCanvas()
	// Loop through all the hit markers that are stored for display
	for _, hit := range d.hitPositions {
		// Get screen coordinates for this hit
		xCentre, yCentre := boardgeo.GetDrawingXY(hit)
		xCentre += d.imageMin.X
		yCentre += d.imageMin.Y
		//	Draw a tiny filled circle at this point
		segment, score, _ := boardgeo.DescribeBoardPoint(hit)
		underlyingColour := boardgeo.GetColourForSegment(segment, score)
		red, green, blue := boardgeo.GetContrastingColour(underlyingColour)
		hitPosition := image.Pt(xCentre, yCentre)
		hitColour := color.RGBA{R: red, G: green, B: blue, A: hitMarkerAlpha}
		hitRadius := float32(d.hitMarkerRadius)
		canvas.AddCircleFilled(hitPosition, hitRadius, hitColour)
	}
}
