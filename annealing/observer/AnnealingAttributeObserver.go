// Copyright (c) 2018 Australian Rivers Institute.

package observer

import (
	"github.com/LindsayBradford/crem/annealing"
	"github.com/LindsayBradford/crem/annealing/observer/filters"
	"github.com/LindsayBradford/crem/logging"
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

// ObserveAnnealingEvent captures and converts Event instances into a Attributes instance that
// captures key attributes associated with the event, and passes them to the Logger for processing.
func (aao *AnnealingAttributeObserver) ObserveAnnealingEvent(event annealing.Event) {
	if aao.logHandler.BeingDiscarded(AnnealerLogLevel) || aao.filter.ShouldFilter(event) {
		return
	}

	annealer := wrapAnnealer(event.Annealer)
	explorer := wrapSolutionExplorer(event.Annealer.SolutionExplorer())

	logAttributes := make(logging.Attributes, 0)
	logAttributes = append(logAttributes, logging.NameValuePair{Name: "Id", Value: annealer.Id()})
	logAttributes = append(logAttributes, logging.NameValuePair{Name: "Event", Value: event.EventType.String()})

	switch event.EventType {
	case annealing.StartedAnnealing:
		logAttributes = append(logAttributes,
			logging.NameValuePair{Name: "MaximumIterations", Value: annealer.MaxIterations()},
			logging.NameValuePair{Name: "Temperature", Value: annealer.Temperature()},
			logging.NameValuePair{Name: "CoolingFactor", Value: annealer.CoolingFactor()},
		)
	case annealing.StartedIteration:
		logAttributes = append(logAttributes,
			logging.NameValuePair{Name: "CurrentIteration", Value: annealer.CurrentIteration()},
			logging.NameValuePair{Name: "Temperature", Value: annealer.Temperature()},
			logging.NameValuePair{Name: "ObjectiveValue", Value: explorer.ObjectiveValue()},
		)
	case annealing.FinishedIteration:
		logAttributes = append(logAttributes,
			logging.NameValuePair{Name: "CurrentIteration", Value: annealer.CurrentIteration()},
			logging.NameValuePair{Name: "ObjectiveValue", Value: explorer.ObjectiveValue()},
			logging.NameValuePair{Name: "ChangeInObjectiveValue", Value: explorer.ChangeInObjectiveValue()},
			logging.NameValuePair{Name: "ChangeIsDesirable", Value: explorer.ChangeIsDesirable()},
			logging.NameValuePair{Name: "AcceptanceProbability", Value: explorer.AcceptanceProbability()},
			logging.NameValuePair{Name: "ChangeAccepted", Value: explorer.ChangeAccepted()},
		)
	case annealing.FinishedAnnealing:
		logAttributes = append(logAttributes,
			logging.NameValuePair{Name: "CurrentIteration", Value: annealer.CurrentIteration()},
			logging.NameValuePair{Name: "Temperature", Value: annealer.Temperature()},
		)
	case annealing.Note:
		logAttributes = append(logAttributes, logging.NameValuePair{Name: "Note", Value: event.Note})
	default:
		// deliberately does nothing extra
	}
	aao.logHandler.LogAtLevel(AnnealerLogLevel, logAttributes)
}
