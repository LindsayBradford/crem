// Copyright (c) 2018 Australian Rivers Institute.

package annealers

import (
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
	. "github.com/onsi/gomega"
)
import "testing"

func TestAnneal(t *testing.T) {
	g := NewGomegaWithT(t)

	builder := new(Builder)

	params := parameters.Map{
		MaximumIterations: int64(10000),
	}

	annealer, _ := builder.
		ElapsedTimeTrackingAnnealer().
		WithParameters(params).
		Build()

	elapsedTimeAnnealer, _ := annealer.(*ElapsedTimeTrackingAnnealer)

	g.Expect(
		elapsedTimeAnnealer.ElapsedTime()).To(BeZero(), "Annealer should recorded zero elapsed time")

	actualClone := annealer.DeepClone()

	g.Expect(actualClone).To(Equal(annealer), "Deep clone of annealer should equal clone")

	annealer.Anneal()

	g.Expect(
		elapsedTimeAnnealer.ElapsedTime().Nanoseconds()).To(Not(BeZero()), "Annealer should recorded elapsed time")
}
