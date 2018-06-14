// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package Annealer

type Annealer interface {
	setTemperature(temperature float64)
	Temperature() float64

	setCoolingFactor(coolingFactor float64)
	CoolingFactor() float64

	setMaxIterations(iterations uint)
	MaxIterations() uint

	Initialise()

	CurrentIteration() uint

	Anneal()
}
