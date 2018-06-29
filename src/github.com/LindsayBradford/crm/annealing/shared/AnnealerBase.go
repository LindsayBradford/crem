// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package shared

import (
	"errors"

	. "github.com/LindsayBradford/crm/annealing/objectives"
	. "github.com/LindsayBradford/crm/logging/handlers"
)

type AnnealerBase struct {
	temperature      float64
	coolingFactor    float64
	maxIterations    uint
	currentIteration uint
	observers        []AnnealingObserver
	objectiveManager ObjectiveManager
	logger           LogHandler
}

func (this *AnnealerBase) Initialise() {
	this.temperature = 1
	this.coolingFactor = 1
	this.maxIterations = 0
	this.currentIteration = 0
	this.objectiveManager = new(DumbObjectiveManager)
	this.logger = new(NullLogHandler)
}

func (this *AnnealerBase) SetTemperature(temperature float64) error {
	if temperature <= 0 {
		return errors.New("Invalid attempt to set annealer temperature to value <= 0")
	}
	this.temperature = temperature
	return nil
}

func (this *AnnealerBase) Temperature() float64 {
	return this.temperature
}

func (this *AnnealerBase) SetCoolingFactor(coolingFactor float64) error {
	if coolingFactor <= 0 || coolingFactor > 1 {
		return errors.New("Invalid attempt to set annealer cooling factor to value <= 0 or > 1")
	}
	this.coolingFactor = coolingFactor
	return nil
}

func (this *AnnealerBase) CoolingFactor() float64 {
	return this.coolingFactor
}

func (this *AnnealerBase) SetMaxIterations(iterations uint) {
	this.maxIterations = iterations
}

func (this *AnnealerBase) MaxIterations() uint {
	return this.maxIterations
}

func (this *AnnealerBase) CurrentIteration() uint {
	return this.currentIteration
}

func (this *AnnealerBase) ObjectiveManager() ObjectiveManager {
	return this.objectiveManager
}

func (this *AnnealerBase) SetObjectiveManager(manager ObjectiveManager) error {
	if manager == nil {
		return errors.New("Invalid attempt to set Objective Manager to nil value")
	}
	this.objectiveManager = manager
	return nil
}

func (this *AnnealerBase) SetLogHandler(logger LogHandler) error {
	if logger == nil {
		return errors.New("Invalid attempt to set log handler to nil value")
	}
	this.logger = logger
	return nil
}

func (this *AnnealerBase) LogHandler() LogHandler {
	return this.logger
}

func (this *AnnealerBase) AddObserver(newObserver AnnealingObserver) error {
	if newObserver == nil {
		return errors.New("Invalid attempt to add non-existant observer to annealer")
	}
	this.observers = append(this.observers, newObserver)
	return nil
}

func (this *AnnealerBase) notifyObserversWith(thisNote string) {
	event := AnnealingEvent{
		EventType: NOTE,
		Annealer:  this,
		Note:      thisNote}
	this.notifyObserversWithEvent(event)
}

func (this *AnnealerBase) notifyObservers(thisEventType AnnealingEventType) {
	event := AnnealingEvent{
		EventType: thisEventType,
		Annealer:  this}
	this.notifyObserversWithEvent(event)
}

func (this *AnnealerBase) notifyObserversWithEvent(event AnnealingEvent) {
	for _, currObserver := range this.observers {
		if currObserver != nil {
			currObserver.ObserveAnnealingEvent(event)
		}
	}
}

func (this *AnnealerBase) Anneal() {
	this.objectiveManager.Initialise()

	this.annealingStarted()

for done := this.initialDoneValue(); !done; {
		this.iterationStarted()

		this.objectiveManager.TryRandomChange(this.temperature)

		this.iterationFinished()
		this.cooldown()
		done = this.checkIfDone()
	}

	this.annealingFinished()
}

func (this *AnnealerBase) annealingStarted() {
	this.notifyObservers(STARTED_ANNEALING)
}

func (this *AnnealerBase) iterationStarted() {
	this.currentIteration++
	this.notifyObservers(STARTED_ITERATION)
}

func (this *AnnealerBase) iterationFinished() {
	this.notifyObservers(FINISHED_ITERATION)
}

func (this *AnnealerBase) annealingFinished() {
	this.notifyObservers(FINISHED_ANNEALING)
}

func (this *AnnealerBase) initialDoneValue() bool {
	return this.maxIterations == 0
}

func (this *AnnealerBase) checkIfDone() bool {
	return this.currentIteration >= this.maxIterations
}

func (this *AnnealerBase) cooldown() {
	this.temperature *= this.coolingFactor
}
