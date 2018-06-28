// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package annealing

import (
	. "github.com/LindsayBradford/crm/logging/handlers"
)

type Annealer interface {
	setTemperature(temperature float64) error
	Temperature() float64

	setCoolingFactor(coolingFactor float64) error
	CoolingFactor() float64

	SetObjectiveManager(manager ObjectiveManager) error

	setMaxIterations(iterations uint)
	MaxIterations() uint

	AddObserver(observer AnnealingObserver) error

	setLogHandler(logger LogHandler) error
	LogHandler() LogHandler

	Initialise()

	CurrentIteration() uint

	Anneal()
}
