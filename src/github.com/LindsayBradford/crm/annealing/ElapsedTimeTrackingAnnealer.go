// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package annealing

import (
	. "fmt"
	. "time"
	. "github.com/LindsayBradford/crm/annealing/shared"
)

type ElapsedTimeTrackingAnnealer struct {
	SimpleAnnealer

	startTime  Time
	finishTime Time
}

func (this *ElapsedTimeTrackingAnnealer) Anneal() {
	this.LogHandler().Info("Elapsed-time tracking annealer")

	this.startTime = Now()
	this.SimpleAnnealer.Anneal()
	this.finishTime = Now()

	this.LogHandler().Info(this.generateElapsedTimeString())
}

func (this *ElapsedTimeTrackingAnnealer) generateElapsedTimeString() string {
	return Sprintf("Total elapsed time of annealing = [%v]", this.ElapsedTime())
}

func (this *ElapsedTimeTrackingAnnealer) ElapsedTime() Duration {
	return this.finishTime.Sub(this.startTime)
}
