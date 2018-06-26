// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package logging

import (
	. "github.com/LindsayBradford/crm/annealing"
	"github.com/LindsayBradford/crm/strings"
	. "github.com/LindsayBradford/crm/logging/handlers"
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

	var builder strings.FluentBuilder
	builder.Add("{\"AnnealingEvent\": \"", event.EventType.String(), "\", ")

	switch event.EventType {
	case STARTED_ANNEALING:
		builder.
			Add("\"MaximumIterations\": ", annealer.MaxIterations(), ", ").
			Add("\"Temperature\": ", annealer.Temperature(), ", ").
			Add("\"CoolingFactor\": ", annealer.CoolingFactor())
	case STARTED_ITERATION, FINISHED_ANNEALING:
		builder.
			Add("\"CurrentIteration\": ", annealer.CurrentIteration(), ", ").
			Add("\"Temperature\": ", annealer.Temperature())
	case NOTE:
		builder.Add("\"Note\": \"", event.Note, "\"")
	default:
		// deliberately does nothing extra
	}

	builder.Add("}")

	this.logHandler.Info(builder.String())
}