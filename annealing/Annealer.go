// Copyright (c) 2018 Australian Rivers Institute.

package annealing

import (
	"github.com/LindsayBradford/crem/annealing/explorer"
	"github.com/LindsayBradford/crem/annealing/parameters"
	"github.com/LindsayBradford/crem/logging"
)

type Annealer interface {
	SetId(title string)
	Id() string

	SetParameters(params parameters.Map) error
	Temperature() float64
	CoolingFactor() float64
	MaximumIterations() uint64

	SolutionExplorer() explorer.Explorer
	SetSolutionExplorer(explorer explorer.Explorer) error

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
