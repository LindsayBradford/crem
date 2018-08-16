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

	SolutionExplorer() Explorer
	SetSolutionExplorer(explorer Explorer) error

	SetMaxIterations(iterations uint64)
	MaxIterations() uint64

	SetEventNotifier(notifier AnnealingEventNotifier) error
	AddObserver(observer AnnealingObserver) error

	Observers() []AnnealingObserver

	SetLogHandler(logger LogHandler) error
	LogHandler() LogHandler

	Initialise()

	CurrentIteration() uint64

	Anneal()
}
