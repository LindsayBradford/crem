// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package logging

import (
	. "github.com/LindsayBradford/crm/annealing"
	. "github.com/LindsayBradford/crm/logging/handlers"
	. "github.com/LindsayBradford/crm/logging/shared"
	)

// AnnealingAttributeLogger produces a relevant set of LogAttributes to match any AnnealingEvents received
// and passes those events to its LogHandler for whatever logging is appropriate.
type AnnealingAttributeLogger struct{
	AnnealingLogger
}

func (this *AnnealingAttributeLogger) WithLogHandler(handler LogHandler) *AnnealingAttributeLogger {
	this.logHandler = handler
	return this
}

// ObserveAnnealingEvent captures and converts AnnealingEvent instances into a LogAttributes instance that
// captures key attributes associated with the event, and passes them to the LogHandler for processing.
func (this *AnnealingAttributeLogger) ObserveAnnealingEvent(event AnnealingEvent) {
	annealer := wrap(event.Annealer)

	logAttributes := make(LogAttributes, 0)
	logAttributes = append(logAttributes, NameValuePair{"AnnealingEvent", event.EventType.String()})

	switch event.EventType {
	case STARTED_ANNEALING:
		logAttributes = append(logAttributes,
			NameValuePair{"MaximumIterations", annealer.MaxIterations()},
			NameValuePair{"Temperature", annealer.Temperature()},
			NameValuePair{"CoolingFactor", annealer.CoolingFactor()},
		)
	case STARTED_ITERATION, FINISHED_ANNEALING:
		logAttributes = append(logAttributes,
			NameValuePair{"CurrentIteration", annealer.CurrentIteration()},
			NameValuePair{"Temperature", annealer.Temperature()},
		)
	case NOTE:
		logAttributes = append(logAttributes, NameValuePair{"Note", event.Note})
	default:
		// deliberately does nothing extra
	}
	this.logHandler.InfoWithAttributes(logAttributes)
}