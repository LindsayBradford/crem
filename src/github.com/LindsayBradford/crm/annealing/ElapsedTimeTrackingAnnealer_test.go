// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package annealing

import . "github.com/onsi/gomega"
import "testing"

func TestAnneal(t *testing.T) {
	g := NewGomegaWithT(t)

	builder := new(AnnealerBuilder)

	annealer, _ := builder.
		ElapsedTimeTrackingAnnealer().
		WithStartingTemperature(10).
		WithCoolingFactor(.997).
		WithMaxIterations(100000).
		Build()

	elapsedTimeAnnealer, _:= annealer.(*ElapsedTimeTrackingAnnealer)

	g.Expect(
		elapsedTimeAnnealer.ElapsedTime()).To(BeZero(), "Annealer should recorded zero elapsed time")

	annealer.Anneal()

	g.Expect(
		elapsedTimeAnnealer.ElapsedTime().Nanoseconds()).To(Not(BeZero()), "Annealer should recorded elapsed time")
}