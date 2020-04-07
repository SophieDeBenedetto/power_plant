package calculations

import (
	"math/rand"
	"time"
)

// Calculate calculates the value
func Calculate(max *float64, min *float64, stepSize *float64, value float64) float64 {
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	if value == 0 {
		value = r.Float64()*(*max-*min) + *min
	}
	var nom = (*max-*min)/2 + *min
	var maxStep, minStep float64

	if value < nom {
		maxStep = *stepSize
		minStep = -1 * *stepSize * (value - *min) / (nom - *min)
	} else {
		maxStep = *stepSize * (*max - value) / (*max - nom)
		minStep = -1 * *stepSize
	}
	value += r.Float64()*(maxStep-minStep) + minStep
	return value
}
