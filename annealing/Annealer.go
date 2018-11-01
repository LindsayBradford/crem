// Copyright (c) 2018 Australian Rivers Institute.

// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package annealing

import (
	"github.com/LindsayBradford/crem/annealing/explorer"
	"github.com/LindsayBradford/crem/logging"
)

type Annealer interface {
	SetId(title string)
	Id() string

	SetTemperature(temperature float64) error
	Temperature() float64

	SetCoolingFactor(coolingFactor float64) error
	CoolingFactor() float64

	SolutionExplorer() explorer.Explorer
	SetSolutionExplorer(explorer explorer.Explorer) error

	SetMaxIterations(iterations uint64)
	MaxIterations() uint64

	SetEventNotifier(notifier EventNotifier) error
	AddObserver(observer Observer) error

	Observers() []Observer

	SetLogHandler(logger logging.Logger) error
	LogHandler() logging.Logger

	Initialise()
	Clone() Annealer

	CurrentIteration() uint64

	Anneal()
}
