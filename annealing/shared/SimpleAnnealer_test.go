// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package shared

import (
	"testing"

	"github.com/LindsayBradford/crm/annealing/solution"
	"github.com/LindsayBradford/crm/logging/handlers"
	. "github.com/onsi/gomega"
)

func TestSimpleAnnealer_Initialise(t *testing.T) {
	g := NewGomegaWithT(t)

	annealer := new(SimpleAnnealer)
	annealer.Initialise()

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
		"Annealer should have built with NullLogHandler")

	g.Expect(
		annealer.SolutionExplorer()).To(Equal(solution.NULL_EXPLORER),
		"Annealer should have built with Null Solution Explorer")

	g.Expect(
		annealer.Observers()).To(BeNil(),
		"Annealer should have built with no AnnealerObservers")
}

func TestSimpleAnnealer_Errors(t *testing.T) {
	g := NewGomegaWithT(t)

	annealer := new(SimpleAnnealer)
	annealer.Initialise()

	tempErr := annealer.SetTemperature(-1)

	g.Expect(tempErr).To(Not(BeNil()))
	g.Expect(
		annealer.Temperature()).To(BeIdenticalTo(float64(1)),
		"Annealer should have ignored crap Temperature set attempt")

	coolingFactorErr := annealer.SetCoolingFactor(1.000001)

	g.Expect(coolingFactorErr).To(Not(BeNil()))
	g.Expect(annealer.CoolingFactor()).To(BeIdenticalTo(float64(1)),
		"Annealer should have ignored crap CoolingFactor set attempt")

	logHandlerErr := annealer.SetLogHandler(nil)

	g.Expect(logHandlerErr).To(Not(BeNil()))
	g.Expect(annealer.LogHandler()).To(Equal(handlers.DefaultNullLogHandler),
		"Annealer should have ignored crap LogHandler set attempt")

	explorerErr := annealer.SetSolutionExplorer(nil)

	g.Expect(explorerErr).To(Not(BeNil()))
	g.Expect(annealer.SolutionExplorer()).To(Equal(solution.NULL_EXPLORER),
		"Annealer should have ignored crap Solution Explorer set attempt")

	observersErr := annealer.AddObserver(nil)

	g.Expect(observersErr).To(Not(BeNil()))
	g.Expect(annealer.Observers()).To(BeNil(),
		"Annealer should have ignored crap AnnealerObserver set attempt")
}

func TestSimpleAnnealer_Anneal(t *testing.T) {
	g := NewGomegaWithT(t)

	const startTemperature float64 = 1000.0
	const coolingFactor float64 = 0.5
	const iterations uint64 = 3
	const expectedEndTemperature float64 = ((startTemperature * coolingFactor) * coolingFactor) * coolingFactor

	annealer := new(SimpleAnnealer)
	annealer.Initialise()

	annealer.SetTemperature(startTemperature)
	annealer.SetCoolingFactor(coolingFactor)
	annealer.SetMaxIterations(iterations)

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
		"Annealer should have ended with temperature modified by cooling factor * iterations")
}

type CountingObserver struct {
	eventCounts map[AnnealingEventType]uint64
}

func (this *CountingObserver) ObserveAnnealingEvent(event AnnealingEvent) {
	this.eventCounts[event.EventType] += 1
}

func TestSimpleAnnealer_AddObserver(t *testing.T) {
	g := NewGomegaWithT(t)

	annealer := new(SimpleAnnealer)
	annealer.Initialise()

	const expectedIterations = uint64(3)

	annealer.SetTemperature(1000.0)
	annealer.SetCoolingFactor(0.5)
	annealer.SetMaxIterations(expectedIterations)

	countingObserver := new(CountingObserver)
	countingObserver.eventCounts = make(map[AnnealingEventType]uint64)

	observerError := annealer.AddObserver(countingObserver)

	g.Expect(observerError).To(BeNil())
	g.Expect(annealer.Observers()).To(ContainElement(countingObserver),
		"Annealer should have accepted CountingObserver as new AnnealerObserver")

	annealer.Anneal()

	g.Expect(countingObserver.eventCounts[StartedAnnealing]).To(BeNumerically("==", 1),
		"Annealer should have posted 1 StartedAnnealing event")

	g.Expect(countingObserver.eventCounts[FinishedAnnealing]).To(BeNumerically("==", 1),
		"Annealer should have posted 1 FinishedAnnealing event")

	g.Expect(countingObserver.eventCounts[StartedIteration]).To(BeIdenticalTo(expectedIterations),
		"Annealer should have posted <expectedIterations> of  StartedIteration event")

	g.Expect(countingObserver.eventCounts[FinishedIteration]).To(BeIdenticalTo(expectedIterations),
		"Annealer should have posted <expectedIterations> of  FinishedIteration event")
}

type TryCountingSolutionExplorer struct {
	solution.NullExplorer
	changesTried uint64
}

func (this *TryCountingSolutionExplorer) TryRandomChange(temperature float64) {
	this.changesTried += 1
}

func TestSimpleAnnealer_SetSolutionExplorer(t *testing.T) {
	g := NewGomegaWithT(t)

	annealer := new(SimpleAnnealer)
	annealer.Initialise()

	const expectedTryCount = uint64(3)

	annealer.SetTemperature(1000.0)
	annealer.SetCoolingFactor(0.5)
	annealer.SetMaxIterations(expectedTryCount)

	expectedSolutionExplorer := new(TryCountingSolutionExplorer)

	explorerErr := annealer.SetSolutionExplorer(expectedSolutionExplorer)

	g.Expect(explorerErr).To(BeNil())
	g.Expect(annealer.solutionExplorer).To(BeIdenticalTo(expectedSolutionExplorer),
		"Annealer should have accepted CountingObserver as new Explorer")

	annealer.Anneal()

	g.Expect(expectedSolutionExplorer.changesTried).To(BeIdenticalTo(expectedTryCount),
		"Annealer should have tried same number of changes as iterations")
}

type DummyLogHandler struct {
	handlers.NullLogHandler
}

func TestSimpleAnnealer_SetLogHandler(t *testing.T) {
	g := NewGomegaWithT(t)

	annealer := new(SimpleAnnealer)
	annealer.Initialise()

	annealer.SetTemperature(1000.0)
	annealer.SetCoolingFactor(0.5)
	annealer.SetMaxIterations(3)

	expectedLogHandler := new(DummyLogHandler)

	logHandlerErr := annealer.SetLogHandler(expectedLogHandler)

	g.Expect(logHandlerErr).To(BeNil())
	g.Expect(annealer.logger).To(BeIdenticalTo(expectedLogHandler),
		"Annealer should have accepted DummyLogHandler as new logger")
}
