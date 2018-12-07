// Copyright (c) 2018 Australian Rivers Institute.

// Copyright (c) 2018 Australian Rivers Institute.

// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package annealers

import (
	"testing"

	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer/null"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
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
		annealer.MaximumIterations()).To(BeZero(),
		"Annealer should have built with default iterations of 0")

	g.Expect(
		annealer.CurrentIteration()).To(BeZero(),
		"Annealer should have built with current iteration of 0")

	g.Expect(
		annealer.LogHandler()).To(Equal(loggers.DefaultNullLogger),
		"Annealer should have built with NullLogger")

	g.Expect(
		annealer.SolutionExplorer()).To(Equal(null.NullExplorer),
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
		annealer.Temperature()).To(BeNumerically("==", 1),
		"Annealer should have ignored crap Temperature set attempt")

	coolingFactorParam := parameters.Map{CoolingFactor: 1.5}
	coolingFactorErr := annealer.SetParameters(coolingFactorParam)

	g.Expect(coolingFactorErr).To(Not(BeNil()))
	g.Expect(annealer.CoolingFactor()).To(BeNumerically("==", 1),
		"Annealer should have ignored crap CoolingFactor set attempt")

	logHandlerErr := annealer.SetLogHandler(nil)

	g.Expect(logHandlerErr).To(Not(BeNil()))
	g.Expect(annealer.LogHandler()).To(Equal(loggers.DefaultNullLogger),
		"Annealer should have ignored crap Logger set attempt")

	explorerErr := annealer.SetSolutionExplorer(nil)

	g.Expect(explorerErr).To(Not(BeNil()))
	g.Expect(annealer.SolutionExplorer()).To(Equal(null.NullExplorer),
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
	const expectedEndTemperature = ((startTemperature * coolingFactor) * coolingFactor) * coolingFactor

	annealer := new(SimpleAnnealer)
	annealer.Initialise()

	expectedParams := parameters.Map{
		StartingTemperature: startTemperature,
		CoolingFactor:       coolingFactor,
		MaximumIterations:   int64(iterations),
	}

	annealer.SetParameters(expectedParams)

	g.Expect(
		annealer.CurrentIteration()).To(BeZero(),
		"Annealer should have started with current iteration of 0")

	g.Expect(
		annealer.Temperature()).To(BeNumerically("==", startTemperature),
		"Annealer should have started with expected start temperature")

	annealer.Anneal()

	g.Expect(
		annealer.CurrentIteration()).To(BeNumerically("==", iterations),
		"Annealer should have ended with current iteration = max iterations")

	g.Expect(
		annealer.Temperature()).To(BeNumerically("==", expectedEndTemperature),
		"Annealer should have ended with temperature modified by cooling factor * iterations")
}

func TestSimpleAnnealer_AddObserver(t *testing.T) {
	g := NewGomegaWithT(t)

	annealer := new(SimpleAnnealer)
	annealer.Initialise()

	const expectedIterations = uint64(3)

	expectedParams := parameters.Map{
		StartingTemperature: 1000.0,
		CoolingFactor:       0.5,
		MaximumIterations:   int64(expectedIterations),
	}

	annealer.SetParameters(expectedParams)

	g.Expect(annealer.Observers()).To(BeNil(), "Annealer should start with no observers")

	observerError := annealer.AddObserver(nil)

	g.Expect(observerError).To(Not(BeNil()), "Annealer should have raised an error on adding nil AnnealerObserver")

	countingObserver := new(CountingObserver)
	countingObserver.eventCounts = make(map[annealing.EventType]uint64)

	observerError = annealer.AddObserver(countingObserver)

	g.Expect(observerError).To(BeNil())
	g.Expect(annealer.Observers()).To(ContainElement(countingObserver),
		"Annealer should have accepted CountingObserver as new AnnealerObserver")

	annealer.Anneal()

	g.Expect(countingObserver.eventCounts[annealing.StartedAnnealing]).To(BeNumerically("==", 1),
		"Annealer should have posted 1 StartedAnnealing event")

	g.Expect(countingObserver.eventCounts[annealing.FinishedAnnealing]).To(BeNumerically("==", 1),
		"Annealer should have posted 1 FinishedAnnealing event")

	g.Expect(countingObserver.eventCounts[annealing.StartedIteration]).To(BeNumerically("==", expectedIterations),
		"Annealer should have posted <expectedIterations> of  StartedIteration event")

	g.Expect(countingObserver.eventCounts[annealing.FinishedIteration]).To(BeNumerically("==", expectedIterations),
		"Annealer should have posted <expectedIterations> of  FinishedIteration event")
}

func TestSimpleAnnealer_ConcurrentEventNotifier(t *testing.T) {
	g := NewGomegaWithT(t)

	annealer := new(SimpleAnnealer)
	annealer.Initialise()

	notifier := new(annealing.ConcurrentAnnealingEventNotifier)
	annealer.SetEventNotifier(notifier)
	g.Expect(annealer.eventNotifier).To(Equal(notifier),
		"Annealer should use the event notifier assigned to it")

	const expectedIterations = uint64(3)

	expectedParams := parameters.Map{
		StartingTemperature: 1000.0,
		CoolingFactor:       0.5,
		MaximumIterations:   int64(expectedIterations),
	}

	annealer.SetParameters(expectedParams)

	g.Expect(annealer.Observers()).To(BeNil(), "Annealer should start with no observers")

	observerError := annealer.AddObserver(nil)
	g.Expect(observerError).To(Not(BeNil()), "Annealer should have raised an error on adding nil AnnealerObserver")

	countingObserver := new(CountingObserver)
	countingObserver.eventCounts = make(map[annealing.EventType]uint64)

	observerError = annealer.AddObserver(countingObserver)

	g.Expect(observerError).To(BeNil())
	g.Expect(annealer.Observers()).To(ContainElement(countingObserver),
		"Annealer should have accepted CountingObserver as new AnnealerObserver")

	annealer.Anneal()

	// Poll on last expected event with gomega's Eventually().  Expect rest to hold (without polling).

	g.Eventually(
		func() uint64 {
			return countingObserver.eventCounts[annealing.FinishedAnnealing]
		}).Should(BeNumerically("==", 1),
		"Annealer should have posted 1 FinishedAnnealing event")

	g.Expect(countingObserver.eventCounts[annealing.StartedAnnealing]).To(BeNumerically("==", 1),
		"Annealer should have posted 1 StartedAnnealing event")

	g.Expect(countingObserver.eventCounts[annealing.StartedIteration]).To(BeNumerically("==", expectedIterations),
		"Annealer should have posted <expectedIterations> of  StartedIteration event")

	g.Expect(countingObserver.eventCounts[annealing.FinishedIteration]).To(BeNumerically("==", expectedIterations),
		"Annealer should have posted <expectedIterations> of  FinishedIteration event")
}

type CountingObserver struct {
	eventCounts map[annealing.EventType]uint64
}

func (co *CountingObserver) ObserveAnnealingEvent(event annealing.Event) {
	co.eventCounts[event.EventType] += 1
}

func TestSimpleAnnealer_SetSolutionExplorer(t *testing.T) {
	g := NewGomegaWithT(t)

	annealer := new(SimpleAnnealer)
	annealer.Initialise()

	const expectedTryCount = uint64(3)

	expectedParams := parameters.Map{
		StartingTemperature: 1000.0,
		CoolingFactor:       0.5,
		MaximumIterations:   int64(expectedTryCount),
	}

	annealer.SetParameters(expectedParams)

	expectedSolutionExplorer := new(TryCountingSolutionExplorer)

	explorerErr := annealer.SetSolutionExplorer(expectedSolutionExplorer)

	g.Expect(explorerErr).To(BeNil())
	g.Expect(annealer.solutionExplorer).To(BeIdenticalTo(expectedSolutionExplorer),
		"Annealer should have accepted CountingObserver as new Explorer")

	annealer.Anneal()

	g.Expect(expectedSolutionExplorer.changesTried).To(BeNumerically("==", expectedTryCount),
		"Annealer should have tried same number of changes as iterations")
}

type TryCountingSolutionExplorer struct {
	null.Explorer
	changesTried uint64
}

func (tcse *TryCountingSolutionExplorer) TryRandomChange(temperature float64) {
	tcse.changesTried += 1
}

func TestSimpleAnnealer_SetLogHandler(t *testing.T) {
	g := NewGomegaWithT(t)

	annealer := new(SimpleAnnealer)
	annealer.Initialise()

	expectedParams := parameters.Map{
		StartingTemperature: 1000.0,
		CoolingFactor:       0.5,
		MaximumIterations:   int64(3),
	}

	annealer.SetParameters(expectedParams)

	expectedLogHandler := new(DummyLogHandler)

	logHandlerErr := annealer.SetLogHandler(expectedLogHandler)

	g.Expect(logHandlerErr).To(BeNil())
	g.Expect(annealer.LogHandler()).To(BeIdenticalTo(expectedLogHandler),
		"Annealer should have accepted DummyLogHandler as new logger")
}

type DummyLogHandler struct {
	loggers.NullLogger
}
