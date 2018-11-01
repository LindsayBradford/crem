// Copyright (c) 2018 Australian Rivers Institute.

// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package annealers

import (
	"fmt"
	"time"

	"github.com/LindsayBradford/crem/annealing"
)

type ElapsedTimeTrackingAnnealer struct {
	SimpleAnnealer

	startTime  time.Time
	finishTime time.Time
}

func (annealer *ElapsedTimeTrackingAnnealer) Initialise() {
	annealer.SimpleAnnealer.Initialise()
	annealer.SimpleAnnealer.SetId("Elapsed-Time Tracking Annealer")
}

func (annealer *ElapsedTimeTrackingAnnealer) Anneal() {
	annealer.startTime = time.Now()
	annealer.SimpleAnnealer.Anneal()
	annealer.finishTime = time.Now()

	annealer.LogHandler().Info(annealer.generateElapsedTimeString())
}

func (annealer *ElapsedTimeTrackingAnnealer) generateElapsedTimeString() string {
	return fmt.Sprintf("%s: total elapsed time of annealing = [%v]", annealer.Id(), annealer.ElapsedTime())
}

func (annealer *ElapsedTimeTrackingAnnealer) ElapsedTime() time.Duration {
	return annealer.finishTime.Sub(annealer.startTime)
}

func (annealer *ElapsedTimeTrackingAnnealer) Clone() annealing.Annealer {
	clone := *annealer
	explorerClone := annealer.SolutionExplorer().Clone()
	clone.SetSolutionExplorer(explorerClone)
	return &clone
}
