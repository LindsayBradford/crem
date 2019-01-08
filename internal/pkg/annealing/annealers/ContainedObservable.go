// Copyright (c) 2019 Australian Rivers Institute.

package annealers

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
)

type ContainedObservable struct {
	annealing.ContainedEventNotifier

	explorer.ContainedObservableExplorer

	id                string
	temperature       float64
	coolingFactor     float64
	maximumIterations uint64
	currentIteration  uint64
}

func (co *ContainedObservable) AddObserver(observer annealing.Observer) error {
	return co.EventNotifier().AddObserver(observer)
}

func (co *ContainedObservable) Observers() []annealing.Observer {
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
