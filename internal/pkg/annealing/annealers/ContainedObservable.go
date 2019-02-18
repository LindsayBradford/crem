// Copyright (c) 2019 Australian Rivers Institute.

package annealers

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
)

type ContainedObservable struct {
	observer.ContainedEventNotifier

	explorer.ContainedObservableExplorer

	id                string
	temperature       float64
	coolingFactor     float64
	maximumIterations uint64
	currentIteration  uint64
}

func (co *ContainedObservable) AddObserver(observer observer.Observer) error {
	return co.EventNotifier().AddObserver(observer)
}

func (co *ContainedObservable) Observers() []observer.Observer {
	return co.EventNotifier().Observers()
}

func (co *ContainedObservable) SetId(title string) {
	co.id = title
}

func (co *ContainedObservable) Id() string {
	return co.id
}

func (co *ContainedObservable) Temperature() float64 {
	return co.temperature
}

func (co *ContainedObservable) CoolingFactor() float64 {
	return co.coolingFactor
}

func (co *ContainedObservable) MaximumIterations() uint64 {
	return co.maximumIterations
}

func (co *ContainedObservable) CurrentIteration() uint64 {
	return co.currentIteration
}
