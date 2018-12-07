// Copyright (c) 2018 Australian Rivers Institute.

package annealing

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/pkg/logging"
)

type Annealer interface {
	Initialise()

	parameters.Container
	explorer.Container
	logging.Container

	Observable
	EventNotifierContainer

	Cloneable
	Anneal()
}

type Observable interface {
	AddObserver(observer Observer) error
	Observers() []Observer

	Identifiable
	Temperature() float64
	CoolingFactor() float64
	MaximumIterations() uint64
	CurrentIteration() uint64
}

type Identifiable interface {
	SetId(title string)
	Id() string
}

type Cloneable interface {
	DeepClone() Annealer
}
