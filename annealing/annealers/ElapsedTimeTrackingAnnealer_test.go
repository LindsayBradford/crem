// Copyright (c) 2018 Australian Rivers Institute.

package annealers

import (
	"github.com/LindsayBradford/crem/annealing/parameters"
	. "github.com/onsi/gomega"
)
import "testing"

func TestAnneal(t *testing.T) {
	g := NewGomegaWithT(t)

	builder := new(Builder)

	params := parameters.Map {
		StartingTemperature: float64(10),
		CoolingFactor:       0.997,
		MaximumIterations:   int64(10000),
	}

	annealer, _ := builder.
		ElapsedTimeTrackingAnnealer().
		WithParameters(params).
		Build()

	elapsedTimeAnnealer, _ := annealer.(*ElapsedTimeTrackingAnnealer)

	g.Expect(
		elapsedTimeAnnealer.ElapsedTime()).To(BeZero(), "Annealer should recorded zero elapsed time")

	annealer.Anneal()

	g.Expect(
		elapsedTimeAnnealer.ElapsedTime().Nanoseconds()).To(Not(BeZero()), "Annealer should recorded elapsed time")
}
