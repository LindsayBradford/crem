// Copyright (c) 2018 Australian Rivers Institute.

package annealers

import (
	"github.com/LindsayBradford/crem/annealing"
	"github.com/LindsayBradford/crem/logging"
	"github.com/pkg/errors"

	"github.com/LindsayBradford/crem/annealing/explorer"
	"github.com/LindsayBradford/crem/logging/loggers"
)

type SimpleAnnealer struct {
	id               string
	temperature      float64
	coolingFactor    float64
	maxIterations    uint64
	currentIteration uint64
	eventNotifier    annealing.EventNotifier
	solutionExplorer explorer.Explorer
	logger           logging.Logger
}

func (sa *SimpleAnnealer) Initialise() {
	sa.id = "Simple Annealer"
	sa.temperature = 1
	sa.coolingFactor = 1
	sa.maxIterations = 0
	sa.currentIteration = 0
	sa.eventNotifier = new(annealing.SynchronousAnnealingEventNotifier)
	sa.solutionExplorer = explorer.NULL_EXPLORER
	sa.logger = new(loggers.NullLogger)
}

func (sa *SimpleAnnealer) SetId(title string) {
	sa.id = title
	sa.solutionExplorer.SetScenarioId(title)
}

func (sa *SimpleAnnealer) Id() string {
	return sa.id
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

func (sa *SimpleAnnealer) SetCoolingFactor(coolingFactor float64) error {
	if coolingFactor <= 0 || coolingFactor > 1 {
		return errors.New("invalid attempt to set annealer cooling factor to value <= 0 or > 1")
	}
	sa.coolingFactor = coolingFactor
	return nil
}

func (sa *SimpleAnnealer) CoolingFactor() float64 {
	return sa.coolingFactor
}

func (sa *SimpleAnnealer) SetMaxIterations(iterations uint64) {
	sa.maxIterations = iterations
}

func (sa *SimpleAnnealer) MaxIterations() uint64 {
	return sa.maxIterations
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
	return sa.maxIterations == 0
}

func (sa *SimpleAnnealer) checkIfDone() bool {
	return sa.currentIteration >= sa.maxIterations
}

func (sa *SimpleAnnealer) cooldown() {
	sa.temperature *= sa.coolingFactor
}
