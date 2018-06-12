// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package Annealer

type Annealer interface {
	Anneal()
	setTemperature(temperature float64)
	Temperature() float64
	setIterationsLeft(iterations uint)
	IterationsLeft() uint
}

type abstractAnnealer struct {
	temperature    float64
	iterationsLeft uint
}

func (annealer *abstractAnnealer) setTemperature(temperature float64) {
	annealer.temperature = temperature
}

func (annealer *abstractAnnealer) Temperature() float64 {
	return annealer.temperature
}

func (annealer *abstractAnnealer) setIterationsLeft(iterations uint) {
	annealer.iterationsLeft = iterations
}

func (annealer *abstractAnnealer) IterationsLeft() uint {
	return annealer.iterationsLeft
}
