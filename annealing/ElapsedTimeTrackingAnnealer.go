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

func (annealer *ElapsedTimeTrackingAnnealer) Initialise() {
	annealer.SimpleAnnealer.Initialise()
	annealer.SimpleAnnealer.SetId("Elapsed-Time Tracking Annealer")
}

func (annealer *ElapsedTimeTrackingAnnealer) Anneal() {
	annealer.startTime = Now()
	annealer.SimpleAnnealer.Anneal()
	annealer.finishTime = Now()

	annealer.LogHandler().Info(annealer.generateElapsedTimeString())
}

func (annealer *ElapsedTimeTrackingAnnealer) generateElapsedTimeString() string {
	return Sprintf("%s: total elapsed time of annealing = [%v]", annealer.Id(), annealer.ElapsedTime())
}

func (annealer *ElapsedTimeTrackingAnnealer) ElapsedTime() Duration {
	return annealer.finishTime.Sub(annealer.startTime)
}

func (annealer *ElapsedTimeTrackingAnnealer) Clone() Annealer {
	clone := *annealer
	explorerClone := annealer.SolutionExplorer().Clone()
	clone.SetSolutionExplorer(explorerClone)
	return &clone
}
