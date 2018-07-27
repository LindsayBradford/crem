// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package annealing

import (
	. "fmt"
	. "github.com/LindsayBradford/crm/annealing/shared"
	. "time"
)

type ElapsedTimeTrackingAnnealer struct {
	SimpleAnnealer

	startTime  Time
	finishTime Time
}

func (annealer *ElapsedTimeTrackingAnnealer) Anneal() {
	annealer.LogHandler().Info("Elapsed-time tracking annealer")

	annealer.startTime = Now()
	annealer.SimpleAnnealer.Anneal()
	annealer.finishTime = Now()

	annealer.LogHandler().Info(annealer.generateElapsedTimeString())
}

func (annealer *ElapsedTimeTrackingAnnealer) generateElapsedTimeString() string {
	return Sprintf("Total elapsed time of annealing = [%v]", annealer.ElapsedTime())
}

func (annealer *ElapsedTimeTrackingAnnealer) ElapsedTime() Duration {
	return annealer.finishTime.Sub(annealer.startTime)
}
