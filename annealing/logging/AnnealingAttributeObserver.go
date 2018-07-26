// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package logging

import (
	. "github.com/LindsayBradford/crm/annealing/shared"
	. "github.com/LindsayBradford/crm/logging/handlers"
	. "github.com/LindsayBradford/crm/logging/modulators"
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

func (this *AnnealingAttributeObserver) WithModulator(modulator LoggingModulator) *AnnealingAttributeObserver {
	this.modulator = modulator
	return this
}

// ObserveAnnealingEvent captures and converts AnnealingEvent instances into a LogAttributes instance that
// captures key attributes associated with the event, and passes them to the LogHandler for processing.
func (this *AnnealingAttributeObserver) ObserveAnnealingEvent(event AnnealingEvent) {
	if this.logHandler.BeingDiscarded(ANNEALER) || this.modulator.ShouldModulate(event) {
		return
	}

	annealer := wrapAnnealer(event.Annealer)
	explorer := wrapSolutionExplorer(event.Annealer.SolutionExplorer())

	logAttributes := make(LogAttributes, 0)
	logAttributes = append(logAttributes, NameValuePair{"Event", event.EventType.String()})

	switch event.EventType {
	case STARTED_ANNEALING:
		logAttributes = append(logAttributes,
			NameValuePair{"MaximumIterations", annealer.MaxIterations()},
			NameValuePair{"Temperature", annealer.Temperature()},
			NameValuePair{"CoolingFactor", annealer.CoolingFactor()},
		)
	case STARTED_ITERATION:
		logAttributes = append(logAttributes,
			NameValuePair{"CurrentIteration", annealer.CurrentIteration()},
			NameValuePair{"Temperature", annealer.Temperature()},
			NameValuePair{"ObjectiveValue", explorer.ObjectiveValue()},
		)
	case FINISHED_ITERATION:
		logAttributes = append(logAttributes,
			NameValuePair{"CurrentIteration", annealer.CurrentIteration()},
			NameValuePair{"ObjectiveValue", explorer.ObjectiveValue()},
			NameValuePair{"ChangeInObjectiveValue", explorer.ChangeInObjectiveValue()},
			NameValuePair{"ChangeIsDesirable", explorer.ChangeIsDesirable()},
			NameValuePair{"AcceptanceProbability", explorer.AcceptanceProbability()},
			NameValuePair{"ChangeAccepted", explorer.ChangeAccepted()},
		)
	case FINISHED_ANNEALING:
		logAttributes = append(logAttributes,
			NameValuePair{"CurrentIteration", annealer.CurrentIteration()},
			NameValuePair{"Temperature", annealer.Temperature()},
		)
	case NOTE:
		logAttributes = append(logAttributes, NameValuePair{"Note", event.Note})
	default:
		// deliberately does nothing extra
	}
	this.logHandler.LogAtLevelWithAttributes(ANNEALER, logAttributes)
}
