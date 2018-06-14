// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package Annealer

type annealerBase struct {
	temperature      float64
	coolingFactor    float64
	maxIterations    uint
	currentIteration uint
	observers        []AnnealingObserver
}

func (this *annealerBase) Initialise() {
	this.temperature = 1
	this.coolingFactor = 1
	this.maxIterations = 0
	this.currentIteration = 0
}

func (this *annealerBase) setTemperature(temperature float64) {
	this.temperature = temperature
}

func (this *annealerBase) Temperature() float64 {
	return this.temperature
}

func (this *annealerBase) setCoolingFactor(coolingFactor float64) {
	// PRE: 0 < coolingFactor <= 1
	this.coolingFactor = coolingFactor
}

func (this *annealerBase) CoolingFactor() float64 {
	return this.coolingFactor
}

func (this *annealerBase) setMaxIterations(iterations uint) {
	// PRE: iterations >= 1
	this.maxIterations = iterations
}

func (this *annealerBase) MaxIterations() uint {
	return this.maxIterations
}

func (this *annealerBase) CurrentIteration() uint {
	return this.currentIteration
}

func (this *annealerBase) AddObserver(newObserver AnnealingObserver) {
	this.observers = append(this.observers, newObserver)
}

func (this *annealerBase) notifyObservers(event AnnealingEvent) {
	for _, currObserver := range this.observers {
		if currObserver != nil {
			currObserver.ObserveAnnealingEvent(event, this)
		}
	}
}

func (this *annealerBase) Anneal() {
	this.notifyObservers(STARTED_ANNEALING)

	for done := this.initialDoneValue(); !done; {
		this.currentIteration++
		this.notifyObservers(STARTED_ITERATION)

		// do the actual objective function work here.

		this.cooldown()
		if this.shouldFinish() {
			done = true
		}
	}

	this.notifyObservers(FINISHED_ANNEALING)
}

func (this *annealerBase) initialDoneValue() bool {
	return !(this.maxIterations > 0)
}

func (this *annealerBase) shouldFinish() bool {
	return this.currentIteration >= this.maxIterations
}

func (this *annealerBase) cooldown() {
	this.temperature = this.temperature * this.coolingFactor
}
