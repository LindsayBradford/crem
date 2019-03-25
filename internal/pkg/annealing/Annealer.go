// Copyright (c) 2018 Australian Rivers Institute.

package annealing

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/name"
)

type Annealer interface {
	Initialise()

	parameters.Container
	explorer.Container
	logging.Container

	Observable
	observer.EventNotifierContainer

	Cloneable
	Anneal()
}

type Observable interface {
	name.Identifiable
	AddObserver(observer observer.Observer) error
	Observers() []observer.Observer

	Temperature() float64
	CoolingFactor() float64
	MaximumIterations() uint64
	CurrentIteration() uint64

	ObservableExplorer() explorer.Observable
	SetObservableExplorer(explorer explorer.Observable) error
	Solution() solution.Solution
}

type Cloneable interface {
	DeepClone() Annealer
}
