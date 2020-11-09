// Copyright (c) 2018 Australian Rivers Institute.

package annealers

import (
	"fmt"
	"time"

	"github.com/LindsayBradford/crem/internal/pkg/annealing"
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
	return fmt.Sprintf("Scenario [%s]: total elapsed time of annealing = [%v]", annealer.Id(), annealer.ElapsedTime())
}

func (annealer *ElapsedTimeTrackingAnnealer) ElapsedTime() time.Duration {
	return annealer.finishTime.Sub(annealer.startTime)
}

func (annealer *ElapsedTimeTrackingAnnealer) DeepClone() annealing.Annealer {
	clone := *annealer
	explorerClone := annealer.SolutionExplorer().DeepClone()
	clone.SetSolutionExplorer(explorerClone)

	return &clone
}
