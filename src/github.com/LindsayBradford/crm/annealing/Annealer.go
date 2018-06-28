// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package annealing

type Annealer interface {
	setTemperature(temperature float64) error
	Temperature() float64

	setCoolingFactor(coolingFactor float64) error
	CoolingFactor() float64

	SetObjectiveManager(manager ObjectiveManager) error

	setMaxIterations(iterations uint)
	MaxIterations() uint

	AddObserver(observer AnnealingObserver) error

	Initialise()

	CurrentIteration() uint

	Anneal()
}
