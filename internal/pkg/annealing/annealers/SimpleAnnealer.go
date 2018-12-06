// Copyright (c) 2018 Australian Rivers Institute.

package annealers

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer/null"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/pkg/errors"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
)

type SimpleAnnealer struct {
	id               string
	temperature      float64
	currentIteration uint64
	eventNotifier    annealing.EventNotifier
	solutionExplorer explorer.Explorer
	model            model.Model
	logger           logging.Logger

	parameters Parameters
}

func (sa *SimpleAnnealer) Initialise() {
	sa.id = "Simple Annealer"
	sa.temperature = 1
	sa.currentIteration = 0
	sa.eventNotifier = new(annealing.SynchronousAnnealingEventNotifier)
	sa.solutionExplorer = null.NullExplorer
	sa.logger = new(loggers.NullLogger)
	sa.parameters.Initialise()
}

func (sa *SimpleAnnealer) SetId(title string) {
	sa.id = title
	sa.solutionExplorer.SetScenarioId(title)
}

func (sa *SimpleAnnealer) Id() string {
	return sa.id
}

func (sa *SimpleAnnealer) SetParameters(params parameters.Map) error {
	sa.parameters.Merge(params)

	temperature := sa.parameters.GetFloat64(StartingTemperature)
	sa.SetTemperature(temperature)

	return sa.parameters.ValidationErrors()
}

func (sa *SimpleAnnealer) ParameterErrors() error {
	return sa.parameters.ValidationErrors()
}

func (sa *SimpleAnnealer) SetTemperature(temperature float64) error {
	if temperature <= 0 {
		return errors.New("invalid attempt to set annealer temperature to value <= 0")
	}
	sa.temperature = temperature
	return nil
}

func (sa *SimpleAnnealer) Temperature() float64 {
	return sa.temperature
}

func (sa *SimpleAnnealer) CoolingFactor() float64 {
	return sa.parameters.GetFloat64(CoolingFactor)
}

func (sa *SimpleAnnealer) MaximumIterations() uint64 {
	return uint64(sa.parameters.GetInt64(MaximumIterations))
}

func (sa *SimpleAnnealer) CurrentIteration() uint64 {
	return sa.currentIteration
}

func (sa *SimpleAnnealer) SolutionExplorer() explorer.Explorer {
	return sa.solutionExplorer
}

func (sa *SimpleAnnealer) SetSolutionExplorer(explorer explorer.Explorer) error {
	if explorer == nil {
		return errors.New("invalid attempt to set Solution Explorer to nil value")
	}
	sa.solutionExplorer = explorer
	return nil
}

func (sa *SimpleAnnealer) SetLogHandler(logger logging.Logger) error {
	if logger == nil {
		return errors.New("invalid attempt to set log handler to nil value")
	}
	sa.logger = logger
	return nil
}

func (sa *SimpleAnnealer) LogHandler() logging.Logger {
	return sa.logger
}

func (sa *SimpleAnnealer) SetEventNotifier(delegate annealing.EventNotifier) error {
	if delegate == nil {
		return errors.New("invalid attempt to set event notifier to nil value")
	}
	sa.eventNotifier = delegate
	return nil
}

func (sa *SimpleAnnealer) AddObserver(observer annealing.Observer) error {
	return sa.eventNotifier.AddObserver(observer)
}

func (sa *SimpleAnnealer) Observers() []annealing.Observer {
	return sa.eventNotifier.Observers()
}

func (sa *SimpleAnnealer) notifyObservers(eventType annealing.EventType) {
	sa.eventNotifier.NotifyObserversOfAnnealingEvent(sa.Clone(), eventType)
}

func (sa *SimpleAnnealer) Clone() annealing.Annealer {
	clone := *sa
	explorerClone := sa.SolutionExplorer().Clone()
	clone.SetSolutionExplorer(explorerClone)
	return &clone
}

func (sa *SimpleAnnealer) Anneal() {
	sa.solutionExplorer.SetLogHandler(sa.LogHandler())

	defer func() {
		if r := recover(); r != nil {
			baseError, ok := r.(error)
			if ok {
				wrappingError := errors.Wrap(baseError, "annealing function failed")
				panic(wrappingError)
			}
		}
	}()

	sa.solutionExplorer.Initialise()
	defer sa.solutionExplorer.TearDown()

	sa.annealingStarted()

	for done := sa.initialDoneValue(); !done; {
		sa.iterationStarted()

		sa.solutionExplorer.TryRandomChange(sa.temperature)

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
