// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package shared

import (
	"errors"

	. "github.com/LindsayBradford/crm/annealing/objectives"
	. "github.com/LindsayBradford/crm/logging/handlers"
)

type SimpleAnnealer struct {
	temperature      float64
	coolingFactor    float64
	maxIterations    uint
	currentIteration uint
	eventNotifier    AnnealingEventNotifier
	objectiveManager ObjectiveManager
	logger           LogHandler
}

func (this *SimpleAnnealer) Initialise() {
	this.temperature = 1
	this.coolingFactor = 1
	this.maxIterations = 0
	this.currentIteration = 0
	this.eventNotifier = NULL_EVENT_NOTIFIER
	this.objectiveManager = NULL_OBJECTIVE_MANAGER
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

func (this *SimpleAnnealer) ObjectiveManager() ObjectiveManager {
	return this.objectiveManager
}

func (this *SimpleAnnealer) SetObjectiveManager(manager ObjectiveManager) error {
	if manager == nil {
		return errors.New("Invalid attempt to set Objective Manager to nil value")
	}
	this.objectiveManager = manager
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

func (this *SimpleAnnealer) AddObserver(newObserver AnnealingObserver) error {
	return this.eventNotifier.AddObserver(newObserver)
}

func (this *SimpleAnnealer) Observers() []AnnealingObserver {
	return this.eventNotifier.Observers()
}

func (this *SimpleAnnealer) notifyObservers(eventType AnnealingEventType) {
	this.eventNotifier.NotifyObserversOfAnnealingEvent(this, eventType)
}

func (this *SimpleAnnealer) Anneal() {
	this.objectiveManager.SetLogHandler(this.LogHandler())
	this.objectiveManager.Initialise()

	this.annealingStarted()

	for done := this.initialDoneValue(); !done; {
		this.iterationStarted()

		this.objectiveManager.TryRandomChange(this.temperature)

		this.iterationFinished()
		this.cooldown()
		done = this.checkIfDone()
	}

	this.objectiveManager.TearDown()
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
