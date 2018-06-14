// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package Annealer

import . "github.com/onsi/gomega"
import "testing"

func TestBuild_usingDefaults(t *testing.T) {
	g := NewGomegaWithT(t)

	builder := new(AnnealerBuilder)

	annealer := builder.
		SingleObjectiveAnnealer().
		Build()

	g.Expect(
		annealer.Temperature()).To(BeIdenticalTo(float64(1)),
		"Annealer should have built with default Temperature of 1")

	g.Expect(
		annealer.CoolingFactor()).To(BeIdenticalTo(float64(1)),
		"Annealer should have built with default Cooling Factor of 1")

	g.Expect(
		annealer.MaxIterations()).To(BeZero(),
		"Annealer should have built with default iterations of 0")

	g.Expect(
		annealer.CurrentIteration()).To(BeZero(),
		"Annealer should have built with current iteration of 0")
}

func TestBuild_OverrdingDefaults(t *testing.T) {
	g := NewGomegaWithT(t)

	const expectedTemperature float64 = 1000
	const expectedCoolingFactor float64 = 0.5
	const expectedIterations uint = 5000

	builder := new(AnnealerBuilder)

	annealer := builder.
		SingleObjectiveAnnealer().
		WithStartingTemperature(expectedTemperature).
		WithCoolingFactor(expectedCoolingFactor).
		WithMaxIterations(expectedIterations).
		Build()

	g.Expect(
		annealer.Temperature()).To(BeIdenticalTo(expectedTemperature),
		"Annealer should have built with supplied Temperature")

	g.Expect(
		annealer.CoolingFactor()).To(BeIdenticalTo(expectedCoolingFactor),
		"Annealer should have built with supplied Cooling Factor")

	g.Expect(
		annealer.MaxIterations()).To(BeIdenticalTo(expectedIterations),
		"Annealer should have built with supplied Iterations")

	g.Expect(
		annealer.CurrentIteration()).To(BeZero(),
		"Annealer should have built with current iteration of 0")
}
