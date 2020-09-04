// Copyright (c) 2018 Australian Rivers Institute.

package annealers

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer/null"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
	"github.com/LindsayBradford/crem/pkg/attributes"
	compositeErrors "github.com/LindsayBradford/crem/pkg/errors"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"github.com/LindsayBradford/crem/pkg/name"
	"github.com/pkg/errors"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
)

const (
	Id               = "Id"
	CurrentIteration = "CurrentIteration"
)

var _ observer.Observer = new(SimpleAnnealer)

type SimpleAnnealer struct {
	name.IdentifiableContainer

	explorer.ContainedExplorer
	model.ContainedModel

	loggers.ContainedLogger

	parameters Parameters

	observer.ContainedEventNotifier

	maximumIterations uint64
	currentIteration  uint64
}

func (sa *SimpleAnnealer) Initialise() {
	sa.SetSolutionExplorer(null.NullExplorer)
	sa.SetLogHandler(new(loggers.NullLogger))
	sa.SetEventNotifier(new(observer.SynchronousAnnealingEventNotifier))

	sa.SetId("Simple Annealer")

	sa.currentIteration = 0

	sa.parameters.Initialise()
	sa.assignStateFromParameters()
}

func (sa *SimpleAnnealer) SetId(title string) {
	sa.IdentifiableContainer.SetId(title)
	sa.SolutionExplorer().SetId(title)
}

func (sa *SimpleAnnealer) SetLogHandler(logHandler logging.Logger) {
	sa.ContainedLogger.SetLogHandler(logHandler)
	sa.SolutionExplorer().SetLogHandler(logHandler)
}

func (sa *SimpleAnnealer) SetModel(model model.Model) {
	model.SetId(sa.Id())
	sa.SolutionExplorer().SetModel(model)
}

func (sa *SimpleAnnealer) SetParameters(params parameters.Map) error {
	sa.parameters.AssignOnlyEnforcedUserValues(params)

	sa.SolutionExplorer().SetParameters(params)

	sa.assignStateFromParameters()
	return sa.parameters.ValidationErrors()
}

func (sa *SimpleAnnealer) assignStateFromParameters() {
	sa.maximumIterations = uint64(sa.parameters.GetInt64(MaximumIterations))
}

func (sa *SimpleAnnealer) ParameterErrors() error {
	mergedErrors := compositeErrors.New("Kirkpatrick Explorer Parameter Validation")

	mergedErrors.Add(sa.parameters.ValidationErrors())
	mergedErrors.Add(sa.SolutionExplorer().ParameterErrors())

	if mergedErrors.Size() > 0 {
		return mergedErrors
	}

	return nil
}

func (sa *SimpleAnnealer) DeepClone() annealing.Annealer {
	clone := *sa
	explorerClone := sa.SolutionExplorer().DeepClone()
	clone.SetSolutionExplorer(explorerClone)
	return &clone
}

func (sa *SimpleAnnealer) Model() model.Model {
	return sa.SolutionExplorer().Model()
}

func (sa *SimpleAnnealer) Anneal() {
	defer sa.handlePanicRecovery()

	sa.SolutionExplorer().Initialise()
	defer sa.SolutionExplorer().TearDown()

	sa.annealingStarted()

	for done := sa.initialDoneValue(); !done; {
		sa.iterationStarted()

		sa.SolutionExplorer().TryRandomChange()
		sa.SolutionExplorer().CoolDown()

		sa.iterationFinished()
		done = sa.checkIfDone()
	}

	sa.annealingFinished()
}

func (sa *SimpleAnnealer) handlePanicRecovery() {
	if r := recover(); r != nil {
		if rAsError, isError := r.(error); isError {
			wrappingError := errors.Wrap(rAsError, "Unrecoverable annealing failure. Exiting with error: ")
			sa.LogHandler().Error(wrappingError)
			panic(wrappingError)
		}
		panic(r)
	}
}

func (sa *SimpleAnnealer) annealingStarted() {
	event := sa.newEvent(observer.StartedAnnealing)
	sa.EventNotifier().NotifyObserversOfEvent(*event)
}

func (sa *SimpleAnnealer) iterationStarted() {
	sa.currentIteration++
	event := sa.newEvent(observer.StartedIteration)
	sa.EventNotifier().NotifyObserversOfEvent(*event)
}

func (sa *SimpleAnnealer) iterationFinished() {
	event := sa.newEvent(observer.FinishedIteration)
	sa.EventNotifier().NotifyObserversOfEvent(*event)
}

func (sa *SimpleAnnealer) annealingFinished() {
	event := sa.newEvent(observer.FinishedAnnealing)
	sa.EventNotifier().NotifyObserversOfEvent(*event)
}

func (sa *SimpleAnnealer) newEvent(eventType observer.EventType) *observer.Event {
	eventAttributes := sa.EventAttributes(eventType)
	return observer.NewEvent(eventType).JoiningAttributes(eventAttributes)
}

func (sa *SimpleAnnealer) EventAttributes(eventType observer.EventType) attributes.Attributes {
	baseAttributes := new(attributes.Attributes).
		Add(Id, sa.Id()).
		Add(MaximumIterations, sa.maximumIterations).
		Join(sa.SolutionExplorer().EventAttributes(eventType))
	switch eventType {
	case observer.StartedAnnealing:
		return baseAttributes
	case observer.StartedIteration:
		return new(attributes.Attributes).
			Add(Id, sa.Id()).
			Add(CurrentIteration, sa.currentIteration).
			Add(MaximumIterations, sa.maximumIterations).
			Join(sa.SolutionExplorer().EventAttributes(eventType))
	case observer.FinishedIteration:
		return new(attributes.Attributes).
			Add(Id, sa.Id()).
			Add(CurrentIteration, sa.currentIteration).
			Add(MaximumIterations, sa.maximumIterations).
			Join(sa.SolutionExplorer().EventAttributes(eventType))
	case observer.FinishedAnnealing:
		return new(attributes.Attributes).
			Add(Id, sa.Id()).
			Add(CurrentIteration, sa.currentIteration).
			Add(MaximumIterations, sa.maximumIterations).
			Join(sa.SolutionExplorer().EventAttributes(eventType))
	}
	return nil
}

func (sa *SimpleAnnealer) initialDoneValue() bool {
	return sa.parameters.GetInt64(MaximumIterations) == 0
}

func (sa *SimpleAnnealer) checkIfDone() bool {
	return sa.currentIteration >= uint64(sa.parameters.GetInt64(MaximumIterations))
}

func (sa *SimpleAnnealer) AddObserver(observer observer.Observer) error {
	return sa.EventNotifier().AddObserver(observer)
}

func (sa *SimpleAnnealer) Observers() []observer.Observer {
	return sa.EventNotifier().Observers()
}

func (sa *SimpleAnnealer) ObserveEvent(event observer.Event) {
	wrappingEvent := observer.NewEvent(event.EventType).WithId(sa.Id())

	if event.EventType.IsAnnealingIterationState() {
		wrappingEvent.
			WithAttribute(CurrentIteration, sa.currentIteration).
			WithAttribute(MaximumIterations, sa.maximumIterations)
	}

	wrappingEvent.JoiningAttributes(event.AllAttributes())
	sa.EventNotifier().NotifyObserversOfEvent(*wrappingEvent)
}

func (sa *SimpleAnnealer) SetSolutionExplorer(explorer explorer.Explorer) error {
	explorerError := sa.ContainedExplorer.SetSolutionExplorer(explorer)
	if explorerError != nil {
		return explorerError
	}
	explorer.SetId(sa.Id())
	return nil
}
