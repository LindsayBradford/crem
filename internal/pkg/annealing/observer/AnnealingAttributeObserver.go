// Copyright (c) 2018 Australian Rivers Institute.

package observer

import (
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

	if event.EventType.IsAnnealingState() {
		aao.observeAnnealingEvent(event, logAttributes)
	} else {
		aao.observeEvent(event, logAttributes)
	}
}

func (aao *AnnealingAttributeObserver) observeAnnealingEvent(event observer.Event, logAttributes attributes.Attributes) {
	switch event.EventType {
	case observer.StartedAnnealing:
		logAttributes = append(logAttributes,
			event.AttributesNamed("MaximumIterations", "Temperature", "CoolingFactor")...)
	case observer.StartedIteration:
		logAttributes = append(logAttributes,
			event.AttributesNamed("CurrentIteration", "Temperature", "ObjectiveValue")...)
	case observer.FinishedIteration:
		logAttributes = append(logAttributes,
			event.AttributesNamed(
				"CurrentIteration", "ObjectiveValue", "ChangeInObjectiveValue",
				"ChangeIsDesirable", "AcceptanceProbability", "ChangeAccepted")...)
	case observer.FinishedAnnealing:
		logAttributes = event.AllAttributes()
	default:
		// deliberately does nothing extra
	}
	aao.logHandler.LogAtLevelWithAttributes(AnnealerLogLevel, logAttributes)
}

func (aao *AnnealingAttributeObserver) observeEvent(event observer.Event, logAttributes attributes.Attributes) {
	switch event.EventType {
	case observer.Note:
		logAttributes = append(logAttributes, event.AttributesNamed("Note")...)
	case observer.ManagementAction:
		logAttributes = append(logAttributes, event.AttributesNamed("Type", "PlanningUnit", "Active", "Note")...)
	case observer.DecisionVariable:
		logAttributes = append(logAttributes, event.AttributesNamed("Name", "Value", "Note")...)
	default:
		// deliberately does nothing extra
	}
	aao.logHandler.LogAtLevelWithAttributes(model.LogLevel, logAttributes)
}
