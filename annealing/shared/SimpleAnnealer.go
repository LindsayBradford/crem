// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package shared

import (
	"github.com/pkg/errors"

	. "github.com/LindsayBradford/crem/annealing/solution"
	. "github.com/LindsayBradford/crem/logging/handlers"
)

type SimpleAnnealer struct {
	id               string
	temperature      float64
	coolingFactor    float64
	maxIterations    uint64
	currentIteration uint64
	eventNotifier    AnnealingEventNotifier
	solutionExplorer Explorer
	logger           LogHandler
}

func (sa *SimpleAnnealer) Initialise() {
	sa.id = "Simple Annealer"
	sa.temperature = 1
	sa.coolingFactor = 1
	sa.maxIterations = 0
	sa.currentIteration = 0
	sa.eventNotifier = new(SynchronousAnnealingEventNotifier)
	sa.solutionExplorer = NULL_EXPLORER
	sa.logger = new(NullLogHandler)
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

func (sa *SimpleAnnealer) SolutionExplorer() Explorer {
	return sa.solutionExplorer
}

func (sa *SimpleAnnealer) SetSolutionExplorer(explorer Explorer) error {
	if explorer == nil {
		return errors.New("invalid attempt to set Solution Explorer to nil value")
	}
	sa.solutionExplorer = explorer
	return nil
}

func (sa *SimpleAnnealer) SetLogHandler(logger LogHandler) error {
	if logger == nil {
		return errors.New("invalid attempt to set log handler to nil value")
	}
	sa.logger = logger
	return nil
}

func (sa *SimpleAnnealer) LogHandler() LogHandler {
	return sa.logger
}

func (sa *SimpleAnnealer) SetEventNotifier(delegate AnnealingEventNotifier) error {
	if delegate == nil {
		return errors.New("invalid attempt to set event notifier to nil value")
	}
	sa.eventNotifier = delegate
	return nil
}

func (sa *SimpleAnnealer) AddObserver(observer AnnealingObserver) error {
	return sa.eventNotifier.AddObserver(observer)
}

func (sa *SimpleAnnealer) Observers() []AnnealingObserver {
	return sa.eventNotifier.Observers()
}

func (sa *SimpleAnnealer) notifyObservers(eventType AnnealingEventType) {
	sa.eventNotifier.NotifyObserversOfAnnealingEvent(sa.Clone(), eventType)
}

func (sa *SimpleAnnealer) Clone() Annealer {
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
	sa.notifyObservers(StartedAnnealing)
}

func (sa *SimpleAnnealer) iterationStarted() {
	sa.currentIteration++
	sa.notifyObservers(StartedIteration)
}

func (sa *SimpleAnnealer) iterationFinished() {
	sa.notifyObservers(FinishedIteration)
}

func (sa *SimpleAnnealer) annealingFinished() {
	sa.notifyObservers(FinishedAnnealing)
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
