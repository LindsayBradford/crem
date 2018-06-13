package Annealer

import . "github.com/onsi/gomega"
import "testing"

func TestBuild(t *testing.T) {
	g := NewGomegaWithT(t)

	const expectedTemperature float64 = 1000.0
	const expectedIterations uint = 5000

	builder := new(AnnealerBuilder)

	annealer, _ := builder.
		SingleObjectiveAnnealer().
		WithStartingTemperature(expectedTemperature).
		WithIterations(expectedIterations).
		Build()

	g.Expect(annealer.Temperature()).To(BeIdenticalTo(expectedTemperature), "Annealer should have built with supplied Temperature")
	g.Expect(annealer.IterationsLeft()).To(BeIdenticalTo(expectedIterations), "Annealer should have built with supplied Iterations")
}
