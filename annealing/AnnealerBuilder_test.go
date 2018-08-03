// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package annealing

import (
	"fmt"

	"github.com/LindsayBradford/crm/annealing/shared"
	"github.com/LindsayBradford/crm/annealing/solution"
	"github.com/LindsayBradford/crm/logging/handlers"
	. "github.com/onsi/gomega"
)
import "testing"

type dummyObserver struct{}

func (*dummyObserver) ObserveAnnealingEvent(event shared.AnnealingEvent) {}

func TestBuild_OverridingDefaults(t *testing.T) {
	g := NewGomegaWithT(t)

	const expectedTemperature float64 = 1000
	const expectedCoolingFactor float64 = 0.5
	const expectedIterations uint64 = 5000
	expectedLogHandler := new(handlers.BareBonesLogHandler)
	expectedSolutionExplorer := new(solution.DumbSolutionExplorer)
	expectedObservers := []shared.AnnealingObserver{new(dummyObserver)}

	builder := new(AnnealerBuilder)

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
		"Annealer should have built with supplied SolutionExplorer")

	g.Expect(
		annealer.Observers()).To(Equal(expectedObservers),
		"Annealer should have built with supplied Observers")
}

func TestBuild_BadInputs(t *testing.T) {
	g := NewGomegaWithT(t)

	const badTemperature float64 = -1
	const badCoolingFactor float64 = 1.0000001
	badLogHandler := handlers.LogHandler(nil)
	badExplorer := solution.SolutionExplorer(nil)

	expectedErrors := 5

	builder := new(AnnealerBuilder)

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

	g.Expect(
		annealer.LogHandler()).To(Equal(handlers.DefaultNullLogHandler),
		"Annealer should have built with nullLogHandler")

	g.Expect(
		annealer.SolutionExplorer()).To(Equal(solution.NULL_SOLUTION_EXPLORER),
		"Annealer should have built with Null Solution Explorer")

	g.Expect(
		annealer.Observers()).To(BeNil(),
		"Annealer should have built with no AnnealerObservers")
}

func TestAnnealerBuilder_WithDumbSolutionExplorer(t *testing.T) {
	g := NewGomegaWithT(t)

	expectedObjectiveValue := float64(10)

	expectedSolutionExplorer := new(solution.DumbSolutionExplorer)
	expectedSolutionExplorer.SetObjectiveValue(expectedObjectiveValue)

	builder := new(AnnealerBuilder)

	annealer, err := builder.
		SimpleAnnealer().
		WithDumbSolutionExplorer(expectedObjectiveValue).
		Build()

	g.Expect(err).To(BeNil(), "Annealer should have built without errors")

	g.Expect(
		annealer.SolutionExplorer()).To(Equal(expectedSolutionExplorer),
		"Annealer should have built with expected DumbSolutionExplorer")

}
