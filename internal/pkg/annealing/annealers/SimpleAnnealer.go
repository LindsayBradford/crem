// Copyright (c) 2018 Australian Rivers Institute.

package annealers

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer/null"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"github.com/pkg/errors"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
)

type SimpleAnnealer struct {
	explorer.ContainedExplorer
	model.ContainedModel

	logging.ContainedLogger

	parameters Parameters

	ContainedObservable
}

func (sa *SimpleAnnealer) Initialise() {
	sa.id = "Simple Annealer"
	sa.temperature = 1
	sa.currentIteration = 0
	sa.SetEventNotifier(new(annealing.SynchronousAnnealingEventNotifier))
	sa.SetLogHandler(new(loggers.NullLogger))
	sa.SetSolutionExplorer(null.NullExplorer)
	sa.parameters.Initialise()
	sa.assignStateFromParameters()
}

func (sa *SimpleAnnealer) SetId(title string) {
	sa.ContainedObservable.SetId(title)
	sa.SolutionExplorer().SetScenarioId(title)
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

func (sa *SimpleAnnealer) notifyObservers(eventType annealing.EventType) {
	sa.EventNotifier().NotifyObserversOfAnnealingEvent(sa.CloneObservable(), eventType)
}

func (sa *SimpleAnnealer) CloneObservable() annealing.Observable {
	observable := sa.ContainedObservable
	explorerClone := sa.SolutionExplorer().CloneObservable()
	observable.SetObservableExplorer(explorerClone)
	return &observable
}

func (sa *SimpleAnnealer) Anneal() {
	sa.SolutionExplorer().SetLogHandler(sa.LogHandler())

	defer func() {
		if r := recover(); r != nil {
			baseError, ok := r.(error)
			if ok {
				wrappingError := errors.Wrap(baseError, "annealing function failed")
				panic(wrappingError)
			}
			panic(r)
		}
	}()

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

func (sa *SimpleAnnealer) annealingStarted() {
	sa.notifyObservers(annealing.StartedAnnealing)
}

func (sa *SimpleAnnealer) iterationStarted() {
	sa.currentIteration++
	sa.notifyObservers(annealing.StartedIteration)
}

func (sa *SimpleAnnealer) iterationFinished() {
	sa.notifyObservers(annealing.FinishedIteration)
}

func (sa *SimpleAnnealer) annealingFinished() {
	sa.notifyObservers(annealing.FinishedAnnealing)
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
