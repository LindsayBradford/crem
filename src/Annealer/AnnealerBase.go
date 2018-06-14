// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package Annealer

import "fmt"

type annealerBase struct {
	temperature      float64
	coolingFactor    float64
	maxIterations    uint
	currentIteration uint
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

func (this *annealerBase) Anneal() {
	fmt.Printf("Start Temperature: %f\n", this.temperature)
	fmt.Printf("Max Iterations: %d\n", this.maxIterations)

	fmt.Println("Starting Annealing")

	for done := this.initialDoneValue(); !done; {
		this.currentIteration++

		fmt.Printf("  currentIteration : %d\n", this.currentIteration)

		this.cooldown()
		if this.shouldFinish() {
			done = true
		}
	}

	fmt.Println("Finished Annealing")
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
