// Copyright (c) 2018 Australian Rivers Institute.

// Copyright (c) 2018 Australian Rivers Institute.

// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package annealers

import (
	"fmt"

	"github.com/LindsayBradford/crem/annealing"
	"github.com/LindsayBradford/crem/annealing/explorer"
	"github.com/LindsayBradford/crem/logging/handlers"
	. "github.com/onsi/gomega"
)
import "testing"

type dummyObserver struct{}

func (*dummyObserver) ObserveAnnealingEvent(event annealing.Event) {}

func TestBuild_OverridingDefaults(t *testing.T) {
	g := NewGomegaWithT(t)

	const expectedTemperature float64 = 1000
	const expectedCoolingFactor float64 = 0.5
	const expectedIterations uint64 = 5000
	expectedLogHandler := new(handlers.BareBonesLogHandler)
	expectedSolutionExplorer := new(explorer.DumbExplorer)
	expectedObservers := []annealing.Observer{new(dummyObserver)}

	builder := new(Builder)

	annealer, _ := builder.
		SimpleAnnealer().
		WithStartingTemperature(expectedTemperature).
		WithCoolingFactor(expectedCoolingFactor).
		WithMaxIterations(expectedIterations).
		WithLogHandler(expectedLogHandler).
		WithSolutionExplorer(expectedSolutionExplorer).
		WithObservers(expectedObservers...).
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

	g.Expect(
		annealer.LogHandler()).To(BeIdenticalTo(expectedLogHandler),
		"Annealer should have built with supplied LogHandler")

	g.Expect(
		annealer.SolutionExplorer()).To(BeIdenticalTo(expectedSolutionExplorer),
		"Annealer should have built with supplied Explorer")

	g.Expect(
		annealer.Observers()).To(Equal(expectedObservers),
		"Annealer should have built with supplied Observers")
}

func TestBuild_BadInputs(t *testing.T) {
	g := NewGomegaWithT(t)

	const badTemperature float64 = -1
	const badCoolingFactor float64 = 1.0000001
	badLogHandler := handlers.LogHandler(nil)
	badExplorer := explorer.Explorer(nil)

	expectedErrors := 5

	builder := new(Builder)

	annealer, err := builder.
		SimpleAnnealer().
		WithStartingTemperature(badTemperature).
		WithCoolingFactor(badCoolingFactor).
		WithLogHandler(badLogHandler).
		WithSolutionExplorer(badExplorer).
		WithObservers(nil).
		Build()

	g.Expect(
		err.Size()).To(BeIdenticalTo(expectedErrors),
		"Annealer should have built with "+fmt.Sprintf("%d", expectedErrors)+"errors")

	g.Expect(
		annealer).To(BeNil(),
		"Annealer should not have been built")
}

func TestAnnealerBuilder_WithDumbSolutionExplorer(t *testing.T) {
	g := NewGomegaWithT(t)

	expectedObjectiveValue := float64(10)

	expectedSolutionExplorer := new(explorer.DumbExplorer)
	expectedSolutionExplorer.SetObjectiveValue(expectedObjectiveValue)
	expectedSolutionExplorer.SetScenarioId("Simple Annealer")

	builder := new(Builder)

	annealer, err := builder.
		SimpleAnnealer().
		WithDumbSolutionExplorer(expectedObjectiveValue).
		Build()

	g.Expect(err).To(BeNil(), "Annealer should have built without errors")

	g.Expect(
		annealer.SolutionExplorer()).To(Equal(expectedSolutionExplorer),
		"Annealer should have built with expected DumbExplorer")

}
