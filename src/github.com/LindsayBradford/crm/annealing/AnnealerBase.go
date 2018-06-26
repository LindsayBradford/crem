// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package annealing

import (
	"errors"
	"time"
)

type annealerBase struct {
	temperature      float64
	coolingFactor    float64
	maxIterations    uint
	currentIteration uint
	observers        []AnnealingObserver

	startTime time.Time
	finishTime time.Time
}

func (this *annealerBase) Initialise() {
	this.temperature = 1
	this.coolingFactor = 1
	this.maxIterations = 0
	this.currentIteration = 0
}

func (this *annealerBase) setTemperature(temperature float64) error {
	if temperature <= 0 {
		return errors.New("Invalid attempt to set annealer temperature to value <= 0")
	}
	this.temperature = temperature; return nil
}

func (this *annealerBase) Temperature() float64 {
	return this.temperature
}

func (this *annealerBase) setCoolingFactor(coolingFactor float64) error {
	if coolingFactor <= 0 || coolingFactor > 1 {
		return errors.New("Invalid attempt to set annealer cooling factor to value <= 0 or > 1")
	}
	this.coolingFactor = coolingFactor; return nil
}

func (this *annealerBase) CoolingFactor() float64 {
	return this.coolingFactor
}

func (this *annealerBase) setMaxIterations(iterations uint) {
	this.maxIterations = iterations
}

func (this *annealerBase) MaxIterations() uint {
	return this.maxIterations
}

func (this *annealerBase) CurrentIteration() uint {
	return this.currentIteration
}

func (this *annealerBase) AddObserver(newObserver AnnealingObserver) error {
	if newObserver == nil {
		return errors.New("Invalid attempt to add non-existant observer to annealer")
	}
	this.observers = append(this.observers, newObserver); return nil
}

func (this *annealerBase) notifyObserversWith(thisNote string) {
	event := AnnealingEvent{
		EventType: NOTE,
		Annealer:  this,
		Note:      thisNote}
	this.notifyObserversWithEvent(event)
}

func (this *annealerBase) notifyObservers(thisEventType AnnealingEventType) {
	event := AnnealingEvent{
		EventType: thisEventType,
		Annealer:  this}
	this.notifyObserversWithEvent(event)
}

func (this *annealerBase) notifyObserversWithEvent(event AnnealingEvent) {
	for _, currObserver := range this.observers {
		if currObserver != nil {
			currObserver.ObserveAnnealingEvent(event)
		}
	}
}

func (this *annealerBase) Anneal() {
	this.annealingStarted()

	for done := this.initialDoneValue(); !done; {
		this.iterationStarted()

		// do the actual objective function work here.

		this.cooldown()
		if this.shouldFinish() {
			done = true
		}
	}

	this.annealingFinished()
}

func (this *annealerBase) annealingStarted() {
	this.notifyObservers(STARTED_ANNEALING)
	this.startTime = time.Now()
}

func (this *annealerBase) iterationStarted() {
	this.currentIteration++
	this.notifyObservers(STARTED_ITERATION)
}

func (this *annealerBase) annealingFinished() {
	this.notifyObservers(FINISHED_ANNEALING)
	this.finishTime = time.Now()
}

func (this *annealerBase) ElapsedTime() time.Duration {
	return this.startTime.Sub(this.finishTime)
}

func (this *annealerBase) initialDoneValue() bool {
	return this.maxIterations == 0
}

func (this *annealerBase) shouldFinish() bool {
	return this.currentIteration >= this.maxIterations
}

func (this *annealerBase) cooldown() {
	this.temperature *= this.coolingFactor
}
