package model

// AccuracyModel is an abstract model that is used to determine a simulated thrower's accuracy -
// how close they will come to their intended target when they throw.
// Accuracy is determined by the type of accuracy model used - a variety of
// implementations will provide models of different levels of complexity.
type AccuracyModel interface {
	Initialize() error
	GetThrow(targetRadius float64, targetAngleDegrees float64) (float64, float64, error)
}
