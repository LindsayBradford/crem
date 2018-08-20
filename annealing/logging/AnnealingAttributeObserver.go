// Copyright (c) 2018 Australian Rivers Institute.

package logging

import (
	. "github.com/LindsayBradford/crm/annealing/shared"
	. "github.com/LindsayBradford/crm/logging/filters"
	. "github.com/LindsayBradford/crm/logging/handlers"
	. "github.com/LindsayBradford/crm/logging/shared"
)

// AnnealingAttributeObserver produces a relevant set of LogAttributes to match any AnnealingEvents received
// and passes those events to its LogHandler for whatever logging is appropriate.
type AnnealingAttributeObserver struct {
	AnnealingLogger
}

func (aao *AnnealingAttributeObserver) WithLogHandler(handler LogHandler) *AnnealingAttributeObserver {
	aao.logHandler = handler
	return aao
}

func (aao *AnnealingAttributeObserver) WithFilter(Filter LoggingFilter) *AnnealingAttributeObserver {
	aao.filter = Filter
	return aao
}

// ObserveAnnealingEvent captures and converts AnnealingEvent instances into a LogAttributes instance that
// captures key attributes associated with the event, and passes them to the LogHandler for processing.
func (aao *AnnealingAttributeObserver) ObserveAnnealingEvent(event AnnealingEvent) {
	if aao.logHandler.BeingDiscarded(AnnealerLogLevel) || aao.filter.ShouldFilter(event) {
		return
	}

	annealer := wrapAnnealer(event.Annealer)
	explorer := wrapSolutionExplorer(event.Annealer.SolutionExplorer())

	logAttributes := make(LogAttributes, 0)
	logAttributes = append(logAttributes, NameValuePair{Name: "Event", Value: event.EventType.String()})

	switch event.EventType {
	case StartedAnnealing:
		logAttributes = append(logAttributes,
			NameValuePair{Name: "MaximumIterations", Value: annealer.MaxIterations()},
			NameValuePair{Name: "Temperature", Value: annealer.Temperature()},
			NameValuePair{Name: "CoolingFactor", Value: annealer.CoolingFactor()},
		)
	case StartedIteration:
		logAttributes = append(logAttributes,
			NameValuePair{Name: "CurrentIteration", Value: annealer.CurrentIteration()},
			NameValuePair{Name: "Temperature", Value: annealer.Temperature()},
			NameValuePair{Name: "ObjectiveValue", Value: explorer.ObjectiveValue()},
		)
	case FinishedIteration:
		logAttributes = append(logAttributes,
			NameValuePair{Name: "CurrentIteration", Value: annealer.CurrentIteration()},
			NameValuePair{Name: "ObjectiveValue", Value: explorer.ObjectiveValue()},
			NameValuePair{Name: "ChangeInObjectiveValue", Value: explorer.ChangeInObjectiveValue()},
			NameValuePair{Name: "ChangeIsDesirable", Value: explorer.ChangeIsDesirable()},
			NameValuePair{Name: "AcceptanceProbability", Value: explorer.AcceptanceProbability()},
			NameValuePair{Name: "ChangeAccepted", Value: explorer.ChangeAccepted()},
		)
	case FinishedAnnealing:
		logAttributes = append(logAttributes,
			NameValuePair{Name: "CurrentIteration", Value: annealer.CurrentIteration()},
			NameValuePair{Name: "Temperature", Value: annealer.Temperature()},
		)
	case Note:
		logAttributes = append(logAttributes, NameValuePair{Name: "Note", Value: event.Note})
	default:
		// deliberately does nothing extra
	}
	aao.logHandler.LogAtLevel(AnnealerLogLevel, logAttributes)
}
