// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package logging

import (
	. "github.com/LindsayBradford/crm/annealing"
	. "github.com/LindsayBradford/crm/logging/handlers"
	. "github.com/LindsayBradford/crm/logging/shared"
	)

// JsonMessageAnnealingLogger produces a stream of text log entries from any observed
// AnnealingEvent, logging the event as a JSON encoding of the event's content.
type JsonMessageAnnealingLogger struct{
	AnnealingLogger
}

func (this *JsonMessageAnnealingLogger) WithLogHandler(handler LogHandler) *JsonMessageAnnealingLogger {
	this.logHandler = handler
	return this
}

// ObserveAnnealingEvent captures and converts AnnealingEvent instances into a JSON encoding of the event that it
// then passes onto its relevant LogHandler as an Info call.
func (this *JsonMessageAnnealingLogger) ObserveAnnealingEvent(event AnnealingEvent) {
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