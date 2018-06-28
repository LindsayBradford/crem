// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package annealing

import (
	. "fmt"
	. "time"
)

type ElapsedTimeTrackingAnnealer struct {
	annealerBase

	startTime Time
	finishTime Time
}

func (this *ElapsedTimeTrackingAnnealer) Anneal() {
	this.logger.Info("Elapsed-time tracking annealer")

	this.startTime = Now()
	this.annealerBase.Anneal()
	this.finishTime = Now()

	this.logger.Info(this.generateElapsedTimeString())
}

func (this *ElapsedTimeTrackingAnnealer) generateElapsedTimeString() string {
	return Sprintf("Total elapsed time of annealing = [%v]", this.ElapsedTime())
}

func (this *ElapsedTimeTrackingAnnealer) ElapsedTime() Duration {
	return this.finishTime.Sub(this.startTime)
}
