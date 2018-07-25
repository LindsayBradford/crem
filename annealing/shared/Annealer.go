// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package shared

import (
	. "github.com/LindsayBradford/crm/annealing/solution"
	. "github.com/LindsayBradford/crm/logging/handlers"
)

type Annealer interface {
	SetTemperature(temperature float64) error
	Temperature() float64

	SetCoolingFactor(coolingFactor float64) error
	CoolingFactor() float64

	SolutionTourer() SolutionTourer
	SetSolutionTourer(tourer SolutionTourer) error

	SetMaxIterations(iterations uint)
	MaxIterations() uint

	SetEventNotifier(notifier AnnealingEventNotifier) error
	AddObserver(observer AnnealingObserver) error
	Observers() []AnnealingObserver

	SetLogHandler(logger LogHandler) error
	LogHandler() LogHandler

	Initialise()

	CurrentIteration() uint

	Anneal()
}