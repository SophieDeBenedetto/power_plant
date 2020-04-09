package calculations

import (
	"math/rand"
	"time"
)

// Calculator calculates sensor frequency
type Calculator struct {
	Max      *float64
	Min      *float64
	StepSize *float64
}

// New returns a new calculator
func New(max *float64, min *float64, stepSize *float64) *Calculator {
	return &Calculator{
		Max:      max,
		Min:      min,
		StepSize: stepSize,
	}
}

// Calculate calculates the value
func (c *Calculator) Calculate(value float64) float64 {
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	if value == 0 {
		value = r.Float64()*(*c.Max-*c.Min) + *c.Min
	}
	var nom = (*c.Max-*c.Min)/2 + *c.Min
	var maxStep, minStep float64

	if value < nom {
		maxStep = *c.StepSize
		minStep = -1 * *c.StepSize * (value - *c.Min) / (nom - *c.Min)
	} else {
		maxStep = *c.StepSize * (*c.Max - value) / (*c.Max - nom)
		minStep = -1 * *c.StepSize
	}
	value += r.Float64()*(maxStep-minStep) + minStep
	return value
}
