// Copyright (c) 2018 Australian Rivers Institute.

package observer

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/observer/filters"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/pkg/attributes"
	"github.com/LindsayBradford/crem/pkg/logging"
)

// AnnealingAttributeObserver produces a relevant set of Attributes to match any AnnealingEvents received
// and passes those events to its Logger for whatever observer is appropriate.
type AnnealingAttributeObserver struct {
	AnnealingObserver
}

func (aao *AnnealingAttributeObserver) WithLogHandler(handler logging.Logger) *AnnealingAttributeObserver {
	aao.logHandler = handler
	return aao
}

func (aao *AnnealingAttributeObserver) WithFilter(Filter filters.Filter) *AnnealingAttributeObserver {
	aao.filter = Filter
	return aao
}

// ObserveEvent captures and converts Event instances into a Attributes instance that
// captures key attributes associated with the event, and passes them to the Logger for processing.
func (aao *AnnealingAttributeObserver) ObserveEvent(event observer.Event) {
	if aao.logHandler.BeingDiscarded(AnnealerLogLevel) || aao.filter.ShouldFilter(event) {
		return
	}

	logAttributes := make(attributes.Attributes, 0)
	logAttributes = append(logAttributes, attributes.NameValuePair{Name: "Id", Value: event.Id()})
	logAttributes = append(logAttributes, attributes.NameValuePair{Name: "Event", Value: event.EventType.String()})

	if observableAnnealer, isAnnealer := event.Source().(annealing.Observable); isAnnealer {
		aao.observeAnnealingEvent(observableAnnealer, event, logAttributes)
	} else {
		aao.observeEvent(event, logAttributes)
	}
}

func (aao *AnnealingAttributeObserver) observeAnnealingEvent(observableAnnealer annealing.Observable, event observer.Event, logAttributes attributes.Attributes) {

	annealer := wrapAnnealer(observableAnnealer)
	explorer := wrapSolutionExplorer(observableAnnealer.ObservableExplorer())

	switch event.EventType {
	case observer.StartedAnnealing:
		logAttributes = append(logAttributes,
			attributes.NameValuePair{Name: "MaximumIterations", Value: annealer.MaximumIterations()},
			attributes.NameValuePair{Name: "Temperature", Value: annealer.Temperature()},
			attributes.NameValuePair{Name: "CoolingFactor", Value: annealer.CoolingFactor()},
		)
	case observer.StartedIteration:
		logAttributes = append(logAttributes,
			attributes.NameValuePair{Name: "CurrentIteration", Value: annealer.CurrentIteration()},
			attributes.NameValuePair{Name: "Temperature", Value: annealer.Temperature()},
			attributes.NameValuePair{Name: "ObjectiveValue", Value: explorer.ObjectiveValue()},
		)
	case observer.FinishedIteration:
		logAttributes = append(logAttributes,
			attributes.NameValuePair{Name: "CurrentIteration", Value: annealer.CurrentIteration()},
			attributes.NameValuePair{Name: "ObjectiveValue", Value: explorer.ObjectiveValue()},
			attributes.NameValuePair{Name: "ChangeInObjectiveValue", Value: explorer.ChangeInObjectiveValue()},
			attributes.NameValuePair{Name: "ChangeIsDesirable", Value: explorer.ChangeIsDesirable()},
			attributes.NameValuePair{Name: "AcceptanceProbability", Value: explorer.AcceptanceProbability()},
			attributes.NameValuePair{Name: "ChangeAccepted", Value: explorer.ChangeAccepted()},
		)
	case observer.FinishedAnnealing:
		logAttributes = append(logAttributes,
			attributes.NameValuePair{Name: "CurrentIteration", Value: annealer.CurrentIteration()},
			attributes.NameValuePair{Name: "Temperature", Value: annealer.Temperature()},
		)
	default:
		// deliberately does nothing extra
	}
	aao.logHandler.LogAtLevelWithAttributes(AnnealerLogLevel, logAttributes)
}

func (aao *AnnealingAttributeObserver) observeEvent(event observer.Event, logAttributes attributes.Attributes) {
	switch event.EventType {
	case observer.Note:
		logAttributes = append(logAttributes, event.Entries("Note")...)
	case observer.ManagementAction:
		logAttributes = append(logAttributes, event.Entries("Type", "PlanningUnit", "Active", "Note")...)
	case observer.DecisionVariable:
		logAttributes = append(logAttributes, event.Entries("Name", "Value", "Note")...)
	default:
		// deliberately does nothing extra
	}
	aao.logHandler.LogAtLevelWithAttributes(model.LogLevel, logAttributes)
}
