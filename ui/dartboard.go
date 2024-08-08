package ui

import (
	boardgeo "DStratMC/board-geometry"
	g "github.com/AllenDang/giu"
	"image"
	"image/color"
)

const drawReferenceLines = false
const accuracyCircleThickness = 2

type Dartboard interface {
	SetInfo(windowWidget *g.WindowWidget, texture *g.Texture,
		squareDimension float64,
		imageMin image.Point, imageMax image.Point)
	SetClickCallback(callback func(dartboard Dartboard, position boardgeo.BoardPosition))
	DrawFunction()
	dartboardClicked()
	RemoveThrowMarkers()
	DrawTargetMarker(position boardgeo.BoardPosition)
	DrawAccuracyCircle(position boardgeo.BoardPosition, radius float64)
	GetScoringRadiusPixels() float64
	GetImageMinPoint() image.Point
	GetSquareDimension() float64
	AddHitMarker(hit boardgeo.BoardPosition, markerRadius int)
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
	// We have drawn a circle showing the accuracy radius around a clicked point
	accuracyDrawn    bool
	accuracyPosition boardgeo.BoardPosition
	accuracyRadius   float64
	// Slice of zero or more hits resulting from modeled throw
	hitPositions    []boardgeo.BoardPosition
	hitMarkerRadius int
}

func NewDartboard(clickCallback func(dartboard Dartboard, position boardgeo.BoardPosition)) Dartboard {
	instance := &DartboardInstance{
		clickCallback: clickCallback,
		targetDrawn:   false,
		accuracyDrawn: false,
		hitPositions:  make([]boardgeo.BoardPosition, 0, ThrowsAtOneTarget),
	}
	//fmt.Println("NewDartboard returns", instance)
	return instance
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
	//fmt.Printf("After SetInfo, d = %#v\n", d)
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

const targetCrossLength = 20

// Eventually compute a colour guaranteed to contrast with the target location
var targetCrossColour = color.RGBA{R: 100, G: 100, B: 100, A: 255}

const targetCrossThickness = 2

func (d *DartboardInstance) DrawFunction() {
	if d.squareDimension == 0 {
		//fmt.Println("Squaredimension 0, returning")
		return
	}
	if d.imageMin.X < 0 || d.imageMin.Y < 0 || d.imageMax.X < 0 || d.imageMax.Y < 0 {
		//fmt.Println("imageMin or Max 0, returning")
		return
	}
	//fmt.Println("DartboardCustomFunc")
	//fmt.Printf("imageMin = %#v, imageMax = %#v\n", d.imageMin, d.imageMax)
	//fmt.Println("Square dimension:", d.squareDimension)
	canvas := g.GetCanvas()

	//	Basic test: draw a centred circle
	//radius := d.squareDimension / 2 * .5
	//stubCentre := image.Pt((d.imageMin.X+d.imageMax.X)/2, (d.imageMin.Y+d.imageMax.Y)/2)
	//stubRadius := float32(radius)
	//stubColour := color.RGBA{200, 0, 0, 255}
	//canvas.AddCircleFilled(stubCentre, stubRadius, stubColour)

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

	if drawReferenceLines {
		//	For testing, draw a semitransparent circle on the centre
		xCentre := (d.imageMin.X + d.imageMax.X) / 2
		yCentre := (d.imageMin.Y + d.imageMax.Y) / 2
		testCirclePosition := image.Pt(xCentre, yCentre)
		testCircleColour := color.RGBA{R: 0, G: 0, B: 255, A: 128}
		radius := float32(d.squareDimension / 8.0)
		canvas.AddCircle(testCirclePosition, radius, testCircleColour, 0, 2)

		//	And add centred vertical and horizontal lines to help calibrate angle measurement
		crossHairColour := color.RGBA{R: 150, G: 150, B: 150, A: 255}

		verticalFrom := image.Pt(xCentre, yCentre-int(d.squareDimension/2))
		verticalTo := image.Pt(xCentre, yCentre+int(d.squareDimension/2))
		canvas.AddLine(verticalFrom, verticalTo, crossHairColour, 1)

		horizontalFrom := image.Pt(xCentre-int(d.squareDimension/2), yCentre)
		horizontalTo := image.Pt(xCentre+int(d.squareDimension/2), yCentre)
		canvas.AddLine(horizontalFrom, horizontalTo, crossHairColour, 1)
	}

	//	If we have a target position to draw, do that
	if d.targetDrawn {
		d.DoTargetMarkerDraw(canvas)
	}

	//	If we have an accuracy circle to draw, do that
	if d.accuracyDrawn {
		d.DoAccuracyCircleDraw(canvas)
	}

	d.drawHitMarkers()
}

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

// Despite the name we don't actually draw the target marker here.  It's fine to let
// the caller think that happens. But it would immediately be un-drawn because the
// graphic UI is continually refreshed in a loop.  So, we record the desire for a marker,
// and do the actual redraw in the following function, called from the loop
func (d *DartboardInstance) DrawTargetMarker(position boardgeo.BoardPosition) {
	//fmt.Printf("DrawTargetMarker at %v\n", position)
	d.targetDrawn = true
	d.targetPosition = position
}

func (d *DartboardInstance) DoTargetMarkerDraw(canvas *g.Canvas) {
	//fmt.Printf("DoTargetMarkerDraw at %#v\n", d.targetPosition)
	//	Get the pixel coordinates of this point
	xCentre, yCentre := boardgeo.GetDrawingXY(d.targetPosition)
	xCentre += d.imageMin.X
	yCentre += d.imageMin.Y
	//	Draw an upright cross at this point

	verticalFrom := image.Pt(xCentre, yCentre-targetCrossLength/2)
	verticalTo := image.Pt(xCentre, yCentre+targetCrossLength/2)
	canvas.AddLine(verticalFrom, verticalTo, targetCrossColour, targetCrossThickness)

	horizontalFrom := image.Pt(xCentre-targetCrossLength/2, yCentre)
	horizontalTo := image.Pt(xCentre+targetCrossLength/2, yCentre)
	canvas.AddLine(horizontalFrom, horizontalTo, targetCrossColour, targetCrossThickness)

}

func (d *DartboardInstance) DrawAccuracyCircle(position boardgeo.BoardPosition, radius float64) {
	d.accuracyDrawn = true
	d.accuracyRadius = radius
	d.accuracyPosition = position
}

func (d *DartboardInstance) DoAccuracyCircleDraw(canvas *g.Canvas) {
	xCentre, yCentre := boardgeo.GetDrawingXY(d.accuracyPosition)
	accuracyCirclePosition := image.Pt(xCentre+d.imageMin.X, yCentre+d.imageMin.Y)
	accuracyCircleColour := targetCrossColour
	drawRadius := d.accuracyRadius * d.squareDimension * boardgeo.ScoringAreaFraction / 2
	//drawRadius := 1 * d.squareDimension * boardgeo.ScoringAreaFraction / 2
	canvas.AddCircle(accuracyCirclePosition, float32(drawRadius), accuracyCircleColour, 0, accuracyCircleThickness)
}

func (d *DartboardInstance) AddHitMarker(hit boardgeo.BoardPosition, markerRadius int) {
	d.hitPositions = append(d.hitPositions, hit)
	d.hitMarkerRadius = markerRadius
}

func (d *DartboardInstance) drawHitMarkers() {
	for _, hit := range d.hitPositions {
		xCentre, yCentre := boardgeo.GetDrawingXY(hit)
		xCentre += d.imageMin.X
		yCentre += d.imageMin.Y
		//	Draw a filled circle at this point
		hitPosition := image.Pt(xCentre, yCentre)
		hitColour := color.RGBA{R: 0, G: 0, B: 255, A: 64}
		hitRadius := float32(d.hitMarkerRadius)
		canvas := g.GetCanvas()
		canvas.AddCircleFilled(hitPosition, hitRadius, hitColour)
	}
}
