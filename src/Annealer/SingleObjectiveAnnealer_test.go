// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package Annealer

import . "github.com/onsi/gomega"
import "testing"

func TestAnneal(t *testing.T) {
	g := NewGomegaWithT(t)

	const startTemperature float64 = 1000.0
	const coolingFactor float64 = 0.5
	const iterations uint = 2
	const expectedEndTemperature float64 = (startTemperature * coolingFactor) * coolingFactor

	builder := new(AnnealerBuilder)

	annealer := builder.
		SingleObjectiveAnnealer().
		WithStartingTemperature(startTemperature).
		WithCoolingFactor(coolingFactor).
		WithMaxIterations(iterations).
		Build()

	g.Expect(
		annealer.CurrentIteration()).To(BeZero(),
		"Annealer should have started with current iteration of 0")

	g.Expect(
		annealer.Temperature()).To(BeIdenticalTo(startTemperature),
		"Annealer should have started with expected start temperature")

	annealer.Anneal()

	g.Expect(
		annealer.CurrentIteration()).To(BeIdenticalTo(annealer.MaxIterations()),
		"Annealer should have ended with current iteration = max iterations")

	g.Expect(
		annealer.Temperature()).To(BeIdenticalTo(expectedEndTemperature),
		"Annealer should have ended with tempperature modified by cooling factor * iterations")
}
