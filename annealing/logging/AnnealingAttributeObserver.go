// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

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

func (this *AnnealingAttributeObserver) WithLogHandler(handler LogHandler) *AnnealingAttributeObserver {
	this.logHandler = handler
	return this
}

func (this *AnnealingAttributeObserver) WithFilter(Filter LoggingFilter) *AnnealingAttributeObserver {
	this.filter = Filter
	return this
}

// ObserveAnnealingEvent captures and converts AnnealingEvent instances into a LogAttributes instance that
// captures key attributes associated with the event, and passes them to the LogHandler for processing.
func (this *AnnealingAttributeObserver) ObserveAnnealingEvent(event AnnealingEvent) {
	if this.logHandler.BeingDiscarded(AnnealerLogLevel) || this.filter.ShouldFilter(event) {
		return
	}

	annealer := wrapAnnealer(event.Annealer)
	explorer := wrapSolutionExplorer(event.Annealer.SolutionExplorer())

	logAttributes := make(LogAttributes, 0)
	logAttributes = append(logAttributes, NameValuePair{"Event", event.EventType.String()})

	switch event.EventType {
	case StartedAnnealing:
		logAttributes = append(logAttributes,
			NameValuePair{"MaximumIterations", annealer.MaxIterations()},
			NameValuePair{"Temperature", annealer.Temperature()},
			NameValuePair{"CoolingFactor", annealer.CoolingFactor()},
		)
	case StartedIteration:
		logAttributes = append(logAttributes,
			NameValuePair{"CurrentIteration", annealer.CurrentIteration()},
			NameValuePair{"Temperature", annealer.Temperature()},
			NameValuePair{"ObjectiveValue", explorer.ObjectiveValue()},
		)
	case FinishedIteration:
		logAttributes = append(logAttributes,
			NameValuePair{"CurrentIteration", annealer.CurrentIteration()},
			NameValuePair{"ObjectiveValue", explorer.ObjectiveValue()},
			NameValuePair{"ChangeInObjectiveValue", explorer.ChangeInObjectiveValue()},
			NameValuePair{"ChangeIsDesirable", explorer.ChangeIsDesirable()},
			NameValuePair{"AcceptanceProbability", explorer.AcceptanceProbability()},
			NameValuePair{"ChangeAccepted", explorer.ChangeAccepted()},
		)
	case FinishedAnnealing:
		logAttributes = append(logAttributes,
			NameValuePair{"CurrentIteration", annealer.CurrentIteration()},
			NameValuePair{"Temperature", annealer.Temperature()},
		)
	case Note:
		logAttributes = append(logAttributes, NameValuePair{"Note", event.Note})
	default:
		// deliberately does nothing extra
	}
	this.logHandler.LogAtLevel(AnnealerLogLevel, logAttributes)
}
