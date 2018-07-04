// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford
package annealing

import (
	"fmt"

	"github.com/LindsayBradford/crm/annealing/objectives"
	"github.com/LindsayBradford/crm/annealing/shared"
	"github.com/LindsayBradford/crm/logging/handlers"
	. "github.com/onsi/gomega"
)
import "testing"

type dummyObserver struct {}
func (*dummyObserver) ObserveAnnealingEvent(event shared.AnnealingEvent) {}

func TestBuild_OverridingDefaults(t *testing.T) {
	g := NewGomegaWithT(t)

	const expectedTemperature float64 = 1000
	const expectedCoolingFactor float64 = 0.5
	const expectedIterations uint = 5000
	expectedLogHandler := new(handlers.BareBonesLogHandler)
	expectedObjectiveManager := new(objectives.DumbObjectiveManager)
	expectedObservers := []shared.AnnealingObserver{ new(dummyObserver)}

	builder := new(AnnealerBuilder)

	annealer, _ := builder.
		SimpleAnnealer().
		WithStartingTemperature(expectedTemperature).
		WithCoolingFactor(expectedCoolingFactor).
		WithMaxIterations(expectedIterations).
		WithLogHandler(expectedLogHandler).
		WithObjectiveManager(expectedObjectiveManager).
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
		annealer.ObjectiveManager()).To(BeIdenticalTo(expectedObjectiveManager),
		"Annealer should have built with supplied ObjectiveManager")

	g.Expect(
		annealer.Observers()).To(Equal(expectedObservers),
		"Annealer should have built with supplied AnnealingObservers")
}

func TestBuild_BadInputs(t *testing.T) {
	g := NewGomegaWithT(t)

	const badTemperature float64 = -1
	const badCoolingFactor float64 = 1.0000001
	badLogHandler := handlers.LogHandler(nil)
	badObjectiveManager := objectives.ObjectiveManager(nil)
	// badObserver := shared.AnnealingObserver(nil)

	expectedErrors := 5

	builder := new(AnnealerBuilder)

	annealer, err := builder.
		SimpleAnnealer().
		WithStartingTemperature(badTemperature).
		WithCoolingFactor(badCoolingFactor).
		WithLogHandler(badLogHandler).
		WithObjectiveManager(badObjectiveManager).
		WithObservers(nil).
		Build()

	g.Expect(
		err.Size()).To(BeIdenticalTo(expectedErrors),
		"Annealer should have built with " + fmt.Sprintf("%d", expectedErrors) + "errors")

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
		annealer.LogHandler()).To(Equal(handlers.NULL_LOG_HANDLER),
		"Annealer should have built with nullLogHandler")

	g.Expect(
		annealer.ObjectiveManager()).To(Equal(objectives.NULL_OBJECTIVE_MANAGER),
		"Annealer should have built with nullObjectiveManager")

	g.Expect(
		annealer.Observers()).To(BeNil(),
		"Annealer should have built with no AnnealerObservers")
}

func TestAnnealerBuilder_WithDumbObjectiveManager(t *testing.T) {
	g := NewGomegaWithT(t)

	expectedObjectiveValue := float64(10)

	expectedObjectiveManager := new(objectives.DumbObjectiveManager)
	expectedObjectiveManager.SetObjectiveValue(expectedObjectiveValue)

	builder := new(AnnealerBuilder)

	annealer, err := builder.
		SimpleAnnealer().
		WithDumbObjectiveManager(expectedObjectiveValue).
		Build()

	g.Expect(err).To(BeNil(),"Annealer should have built without errors")

	g.Expect(
		annealer.ObjectiveManager()).To(Equal(expectedObjectiveManager),
		"Annealer should have built with expected DumbObjectiveManager")

}