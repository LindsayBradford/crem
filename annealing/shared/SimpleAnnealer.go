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
	maxIterations    uint
	currentIteration uint
	eventNotifier    AnnealingEventNotifier
	solutionExplorer SolutionExplorer
	logger           LogHandler
}

func (this *SimpleAnnealer) Initialise() {
	this.temperature = 1
	this.coolingFactor = 1
	this.maxIterations = 0
	this.currentIteration = 0
	this.eventNotifier = new(SynchronousAnnealingEventNotifier)
	this.solutionExplorer = NULL_SOLUTION_EXPLORER
	this.logger = NULL_LOG_HANDLER
}

func (this *SimpleAnnealer) SetTemperature(temperature float64) error {
	if temperature <= 0 {
		return errors.New("Invalid attempt to set annealer temperature to value <= 0")
	}
	this.temperature = temperature
	return nil
}

func (this *SimpleAnnealer) Temperature() float64 {
	return this.temperature
}

func (this *SimpleAnnealer) SetCoolingFactor(coolingFactor float64) error {
	if coolingFactor <= 0 || coolingFactor > 1 {
		return errors.New("Invalid attempt to set annealer cooling factor to value <= 0 or > 1")
	}
	this.coolingFactor = coolingFactor
	return nil
}

func (this *SimpleAnnealer) CoolingFactor() float64 {
	return this.coolingFactor
}

func (this *SimpleAnnealer) SetMaxIterations(iterations uint) {
	this.maxIterations = iterations
}

func (this *SimpleAnnealer) MaxIterations() uint {
	return this.maxIterations
}

func (this *SimpleAnnealer) CurrentIteration() uint {
	return this.currentIteration
}

func (this *SimpleAnnealer) SolutionExplorer() SolutionExplorer {
	return this.solutionExplorer
}

func (this *SimpleAnnealer) SetSolutionExplorer(explorer SolutionExplorer) error {
	if explorer == nil {
		return errors.New("Invalid attempt to set Solution Explorer to nil value")
	}
	this.solutionExplorer = explorer
	return nil
}

func (this *SimpleAnnealer) SetLogHandler(logger LogHandler) error {
	if logger == nil {
		return errors.New("Invalid attempt to set log handler to nil value")
	}
	this.logger = logger
	return nil
}

func (this *SimpleAnnealer) LogHandler() LogHandler {
	return this.logger
}

func (this *SimpleAnnealer) SetEventNotifier(delegate AnnealingEventNotifier) error {
	if delegate == nil {
		return errors.New("Invalid attempt to set event notifier to nil value")
	}
	this.eventNotifier = delegate
	return nil
}

func (this *SimpleAnnealer) AddObserver(observer AnnealingObserver) error {
	return this.eventNotifier.AddObserver(observer)
}

func (this *SimpleAnnealer) Observers() []AnnealingObserver {
	return this.eventNotifier.Observers()
}

func (this *SimpleAnnealer) notifyObservers(eventType AnnealingEventType) {
	this.eventNotifier.NotifyObserversOfAnnealingEvent(this.cloneState(), eventType)
}

func (this *SimpleAnnealer) cloneState() *SimpleAnnealer {
	cloneOfThis := *this
	return &cloneOfThis
}

func (this *SimpleAnnealer) Anneal() {
	this.solutionExplorer.SetLogHandler(this.LogHandler())
	this.solutionExplorer.Initialise()

	this.annealingStarted()

	for done := this.initialDoneValue(); !done; {
		this.iterationStarted()

		this.solutionExplorer.TryRandomChange(this.temperature)

		this.iterationFinished()
		this.cooldown()
		done = this.checkIfDone()
	}

	this.solutionExplorer.TearDown()
	this.annealingFinished()
}

func (this *SimpleAnnealer) annealingStarted() {
	this.notifyObservers(STARTED_ANNEALING)
}

func (this *SimpleAnnealer) iterationStarted() {
	this.currentIteration++
	this.notifyObservers(STARTED_ITERATION)
}

func (this *SimpleAnnealer) iterationFinished() {
	this.notifyObservers(FINISHED_ITERATION)
}

func (this *SimpleAnnealer) annealingFinished() {
	this.notifyObservers(FINISHED_ANNEALING)
}

func (this *SimpleAnnealer) initialDoneValue() bool {
	return this.maxIterations == 0
}

func (this *SimpleAnnealer) checkIfDone() bool {
	return this.currentIteration >= this.maxIterations
}

func (this *SimpleAnnealer) cooldown() {
	this.temperature *= this.coolingFactor
}
