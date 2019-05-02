// Copyright (c) 2018 Australian Rivers Institute.

package annealers

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer/null"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"github.com/LindsayBradford/crem/pkg/name"
	"github.com/pkg/errors"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
)

type SimpleAnnealer struct {
	name.IdentifiableContainer

	explorer.ContainedExplorer
	model.ContainedModel

	loggers.ContainedLogger

	parameters Parameters

	observer.ContainedEventNotifier

	temperature   float64
	coolingFactor float64

	maximumIterations uint64
	currentIteration  uint64
}

func (sa *SimpleAnnealer) Initialise() {
	sa.SetSolutionExplorer(null.NullExplorer)
	sa.SetLogHandler(new(loggers.NullLogger))
	sa.SetEventNotifier(new(observer.SynchronousAnnealingEventNotifier))

	sa.SetId("Simple Annealer")

	sa.temperature = 1
	sa.currentIteration = 0

	sa.parameters.Initialise()
	sa.assignStateFromParameters()
}

func (sa *SimpleAnnealer) SetId(title string) {
	sa.IdentifiableContainer.SetId(title)
	sa.SolutionExplorer().SetId(title)
}

func (sa *SimpleAnnealer) SetParameters(params parameters.Map) error {
	sa.parameters.Merge(params)
	sa.assignStateFromParameters()
	return sa.parameters.ValidationErrors()
}

func (sa *SimpleAnnealer) assignStateFromParameters() {
	sa.SetTemperature(sa.parameters.GetFloat64(StartingTemperature))
	sa.coolingFactor = sa.parameters.GetFloat64(CoolingFactor)
	sa.maximumIterations = uint64(sa.parameters.GetInt64(MaximumIterations))
}

func (sa *SimpleAnnealer) ParameterErrors() error {
	return sa.parameters.ValidationErrors()
}

func (sa *SimpleAnnealer) DeepClone() annealing.Annealer {
	clone := *sa
	explorerClone := sa.SolutionExplorer().DeepClone()
	clone.SetSolutionExplorer(explorerClone)
	return &clone
}

func (sa *SimpleAnnealer) SetTemperature(temperature float64) error {
	if temperature <= 0 {
		return errors.New("invalid attempt to set annealer temperature to value <= 0")
	}
	sa.temperature = temperature
	return nil
}

func (sa *SimpleAnnealer) Model() model.Model {
	return sa.SolutionExplorer().Model()
}

func (sa *SimpleAnnealer) Anneal() {
	defer sa.handlePanicRecovery()

	sa.SolutionExplorer().Initialise()
	defer sa.SolutionExplorer().TearDown()

	sa.annealingStarted()

	for done := sa.initialDoneValue(); !done; {
		sa.iterationStarted()

		sa.SolutionExplorer().TryRandomChange(sa.temperature)

		sa.iterationFinished()
		sa.cooldown()
		done = sa.checkIfDone()
	}

	sa.annealingFinished()
}

func (sa *SimpleAnnealer) handlePanicRecovery() {
	if r := recover(); r != nil {
		if rAsError, isError := r.(error); isError {
			wrappingError := errors.Wrap(rAsError, "annealing function failed")
			panic(wrappingError)
		}
		panic(r)
	}
}

func (sa *SimpleAnnealer) annealingStarted() {
	event := observer.NewEvent(observer.StartedAnnealing).
		WithId(sa.Id()).
		WithAttribute("MaximumIterations", sa.maximumIterations).
		WithAttribute("CoolingFactor", sa.coolingFactor).
		JoiningAttributes(
			sa.SolutionExplorer().AttributesForEventType(observer.StartedAnnealing),
		)

	sa.EventNotifier().NotifyObserversOfEvent(*event)
}

func (sa *SimpleAnnealer) iterationStarted() {
	sa.currentIteration++

	event := observer.NewEvent(observer.StartedIteration).
		WithId(sa.Id()).
		WithAttribute("CurrentIteration", sa.currentIteration).
		WithAttribute("MaximumIterations", sa.maximumIterations).
		JoiningAttributes(
			sa.SolutionExplorer().AttributesForEventType(observer.StartedIteration),
		)

	sa.EventNotifier().NotifyObserversOfEvent(*event)
}

func (sa *SimpleAnnealer) iterationFinished() {
	event := observer.NewEvent(observer.FinishedIteration).
		WithId(sa.Id()).
		WithAttribute("CurrentIteration", sa.currentIteration).
		WithAttribute("MaximumIterations", sa.maximumIterations).
		JoiningAttributes(
			sa.SolutionExplorer().AttributesForEventType(observer.FinishedIteration),
		)

	sa.EventNotifier().NotifyObserversOfEvent(*event)
}

func (sa *SimpleAnnealer) annealingFinished() {
	event := observer.NewEvent(observer.FinishedAnnealing).
		WithId(sa.Id()).
		WithAttribute("CurrentIteration", sa.currentIteration).
		WithAttribute("MaximumIterations", sa.maximumIterations).
		JoiningAttributes(
			sa.SolutionExplorer().AttributesForEventType(observer.FinishedAnnealing),
		).
		WithAttribute("Solution", *sa.fetchFinalModelSolution())

	sa.EventNotifier().NotifyObserversOfEvent(*event)
}

func (sa *SimpleAnnealer) fetchFinalModelSolution() *solution.Solution {
	return new(SolutionBuilder).
		WithId(sa.Id()).
		ForModel(sa.SolutionExplorer().Model()).
		Build()
}

func (sa *SimpleAnnealer) initialDoneValue() bool {
	return sa.parameters.GetInt64(MaximumIterations) == 0
}

func (sa *SimpleAnnealer) checkIfDone() bool {
	return sa.currentIteration >= uint64(sa.parameters.GetInt64(MaximumIterations))
}

func (sa *SimpleAnnealer) cooldown() {
	sa.temperature *= sa.parameters.GetFloat64(CoolingFactor)
}

func (sa *SimpleAnnealer) AddObserver(observer observer.Observer) error {
	return sa.EventNotifier().AddObserver(observer)
}

func (sa *SimpleAnnealer) Observers() []observer.Observer {
	return sa.EventNotifier().Observers()
}

func (sa *SimpleAnnealer) Temperature() float64 {
	return sa.temperature
}

func (sa *SimpleAnnealer) CoolingFactor() float64 {
	return sa.coolingFactor
}

func (sa *SimpleAnnealer) MaximumIterations() uint64 {
	return sa.maximumIterations
}

func (sa *SimpleAnnealer) CurrentIteration() uint64 {
	return sa.currentIteration
}
