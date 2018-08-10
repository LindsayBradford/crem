// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package shared

import (
	"errors"

	. "github.com/LindsayBradford/crm/annealing/solution"
	. "github.com/LindsayBradford/crm/logging/handlers"
)

type SimpleAnnealer struct {
	temperature      float64
	coolingFactor    float64
	maxIterations    uint64
	currentIteration uint64
	eventNotifier    AnnealingEventNotifier
	solutionExplorer SolutionExplorer
	logger           LogHandler
}

func (sa *SimpleAnnealer) Initialise() {
	sa.temperature = 1
	sa.coolingFactor = 1
	sa.maxIterations = 0
	sa.currentIteration = 0
	sa.eventNotifier = new(SynchronousAnnealingEventNotifier)
	sa.solutionExplorer = NULL_SOLUTION_EXPLORER
	sa.logger = new(NullLogHandler)
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

func (sa *SimpleAnnealer) SolutionExplorer() SolutionExplorer {
	return sa.solutionExplorer
}

func (sa *SimpleAnnealer) SetSolutionExplorer(explorer SolutionExplorer) error {
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
	sa.eventNotifier.NotifyObserversOfAnnealingEvent(sa.cloneState(), eventType)
}

func (sa *SimpleAnnealer) cloneState() *SimpleAnnealer {
	cloneOfThis := *sa
	return &cloneOfThis
}

func (sa *SimpleAnnealer) Anneal() {
	sa.solutionExplorer.SetLogHandler(sa.LogHandler())
	sa.solutionExplorer.Initialise()

	sa.annealingStarted()

	for done := sa.initialDoneValue(); !done; {
		sa.iterationStarted()

		sa.solutionExplorer.TryRandomChange(sa.temperature)

		sa.iterationFinished()
		sa.cooldown()
		done = sa.checkIfDone()
	}

	sa.solutionExplorer.TearDown()
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
