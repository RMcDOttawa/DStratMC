package ui

import (
	boardgeo "DStratMC/board-geometry"
	g "github.com/AllenDang/giu"
	"image"
	"image/color"
)

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
}

const accuracyCircleThickness = 2
const targetCrossAlpha = 230
const hitMarkerAlpha = 200
const targetCrossLength = 20

var accuracyCircleColour = color.RGBA{R: 100, G: 100, B: 100, A: 255}

const targetCrossThickness = 2

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
	// We have drawn a circle showing the accuracy radius around a clicked point
	accuracyDrawn    bool
	accuracyPosition boardgeo.BoardPosition
	accuracyRadius   float64
	// Slice of zero or more hits resulting from modeled throw
	hitPositions    []boardgeo.BoardPosition
	hitMarkerRadius int
	// Draw the testing crosshair?
	drawReferenceLines bool
}

func (d *DartboardInstance) AllocateHitsSpace(numHits int) {
	d.hitPositions = make([]boardgeo.BoardPosition, 0, numHits)
}

func NewDartboard(clickCallback func(dartboard Dartboard, position boardgeo.BoardPosition)) Dartboard {
	instance := &DartboardInstance{
		clickCallback: clickCallback,
		targetDrawn:   false,
		accuracyDrawn: false,
		hitPositions:  make([]boardgeo.BoardPosition, 0, ThrowsAtOneTarget),
	}
	return instance
}

func (d *DartboardInstance) SetDrawRefLines(checkbox bool) {
	d.drawReferenceLines = checkbox
}

func (d *DartboardInstance) GetSquareDimension() float64 {
	return d.squareDimension
}

func (d *DartboardInstance) GetImageMinPoint() image.Point {
	return d.imageMin
}

func (d *DartboardInstance) GetScoringRadiusPixels() float64 {
	radius := d.squareDimension * boardgeo.ScoringAreaFraction / 2
	return radius
}

func (d *DartboardInstance) SetInfo(windowWidget *g.WindowWidget, texture *g.Texture, squareDimension float64, imageMin image.Point, imageMax image.Point) {
	d.window = windowWidget
	d.texture = texture
	d.squareDimension = squareDimension
	d.imageMin = imageMin
	d.imageMax = imageMax
}

func (d *DartboardInstance) SetClickCallback(callback func(dartboard Dartboard, position boardgeo.BoardPosition)) {
	d.clickCallback = callback
}

func (d *DartboardInstance) RemoveThrowMarkers() {
	d.targetDrawn = false
	d.accuracyDrawn = false
	d.hitPositions = make([]boardgeo.BoardPosition, 0, ThrowsAtOneTarget)
	//fmt.Println("RemoveThrowMarkers STUB")
}

func (d *DartboardInstance) DrawFunction() {
	if d.squareDimension == 0 {
		//fmt.Println("Squaredimension 0, returning")
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

	//	If we have an accuracy circle to draw, do that
	if d.accuracyDrawn {
		d.DrawQueuedAccuracyCircle(canvas)
	}

	d.drawQueuedHitMarkers()
}

// drawReferenceLinesOnDartboard  draws a semitransparent circle and crosshair on the centre
// of the dartboard, to assist with testing coordinates translation
func (d *DartboardInstance) drawReferenceLinesOnDartboard(canvas *g.Canvas) {
	xCentre := (d.imageMin.X + d.imageMax.X) / 2
	yCentre := (d.imageMin.Y + d.imageMax.Y) / 2
	testCirclePosition := image.Pt(xCentre, yCentre)
	testCircleColour := color.RGBA{R: 0, G: 0, B: 255, A: 128}
	radius := float32(d.squareDimension / 8.0)
	canvas.AddCircle(testCirclePosition, radius, testCircleColour, 0, 1)

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
		//fmt.Println("  No callback function")
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
	//fmt.Printf("DrawQueuedTargetMarker at %#v\n", d.targetPosition)
	//	Get the pixel coordinates of this point
	xCentre, yCentre := boardgeo.GetDrawingXY(d.targetPosition)
	xCentre += d.imageMin.X
	yCentre += d.imageMin.Y

	//	Get a contrasting colour that will be visible on this board section
	red, green, blue := contrastingColourForPosition(d.targetPosition)
	colour := color.RGBA{R: red, G: green, B: blue, A: targetCrossAlpha}

	//	Draw an upright cross at this point

	verticalFrom := image.Pt(xCentre, yCentre-targetCrossLength/2)
	verticalTo := image.Pt(xCentre, yCentre+targetCrossLength/2)
	canvas.AddLine(verticalFrom, verticalTo, colour, targetCrossThickness)

	horizontalFrom := image.Pt(xCentre-targetCrossLength/2, yCentre)
	horizontalTo := image.Pt(xCentre+targetCrossLength/2, yCentre)
	canvas.AddLine(horizontalFrom, horizontalTo, colour, targetCrossThickness)
	boardgeo.DescribeBoardPoint(d.targetPosition)
}

// Get RGB values for a colour that contrasts with the colour under the given board position
func contrastingColourForPosition(position boardgeo.BoardPosition) (uint8, uint8, uint8) {
	segment, score, _ := boardgeo.DescribeBoardPoint(position)
	underlyingColour := boardgeo.GetColourForSegment(segment, score)
	red, green, blue := boardgeo.GetContrastingColour(underlyingColour)
	return red, green, blue
}

// QueueAccuracyCircle records the coordinates of a circle that will be drawn on the next UI loop pass
func (d *DartboardInstance) QueueAccuracyCircle(position boardgeo.BoardPosition, radius float64) {
	d.accuracyDrawn = true
	d.accuracyRadius = radius
	d.accuracyPosition = position
}

// DrawQueuedAccuracyCircle draws the accuracy circle that has been recorded
func (d *DartboardInstance) DrawQueuedAccuracyCircle(canvas *g.Canvas) {
	xCentre, yCentre := boardgeo.GetDrawingXY(d.accuracyPosition)
	accuracyCirclePosition := image.Pt(xCentre+d.imageMin.X, yCentre+d.imageMin.Y)
	drawRadius := d.accuracyRadius * d.squareDimension * boardgeo.ScoringAreaFraction / 2
	//drawRadius := 1 * d.squareDimension * boardgeo.ScoringAreaFraction / 2
	canvas.AddCircle(accuracyCirclePosition, float32(drawRadius), accuracyCircleColour, 0, accuracyCircleThickness)
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
