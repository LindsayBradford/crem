// Copyright (c) 2018 Australian Rivers Institute.

package annealing

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/name"
)

type Annealer interface {
	Initialise()

	name.Identifiable

	parameters.Container
	explorer.Container
	logging.Container

	observer.EventNotifierContainer

	AddObserver(observer observer.Observer) error
	Observers() []observer.Observer

	Temperature() float64
	CoolingFactor() float64

	MaximumIterations() uint64
	CurrentIteration() uint64

	Cloneable
	Anneal()
}

type Cloneable interface {
	DeepClone() Annealer
}
