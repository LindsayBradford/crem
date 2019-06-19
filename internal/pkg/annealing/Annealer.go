// Copyright (c) 2018 Australian Rivers Institute.

package annealing

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
	"github.com/LindsayBradford/crem/pkg/attributes"
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

	EventAttributes(eventType observer.EventType) attributes.Attributes

	Cloneable
	Anneal()
}

type Cloneable interface {
	DeepClone() Annealer
}
