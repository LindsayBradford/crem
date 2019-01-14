// Copyright (c) 2019 Australian Rivers Institute.

package observer

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/observer/filters"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/strings"
)

// TODO:  This cannot be run concurently as-is.  Needs a closer look at Annealer state cloning.

type AnnealingInvariantObserver struct {
	AnnealingObserver
	previousObjectiveValue float64
}

func (amo *AnnealingInvariantObserver) WithLogHandler(handler logging.Logger) *AnnealingInvariantObserver {
	amo.logHandler = handler
	return amo
}

func (amo *AnnealingInvariantObserver) WithFilter(filter filters.Filter) *AnnealingInvariantObserver {
	amo.filter = filter
	return amo
}

func (amo *AnnealingInvariantObserver) ObserveAnnealingEvent(event annealing.Event) {
	if amo.loopInvariantUpheld(event) {
		return
	}

	var builder strings.FluentBuilder
	builder.Add(
		"Id [", event.Annealer.Id(), "], ",
		"Event [", event.EventType.String(),
		"]: Loop Invariant Broken",
	)
	amo.logHandler.LogAtLevel(AnnealerLogLevel, builder.String())
	panic(builder.String())
}

func (amo *AnnealingInvariantObserver) loopInvariantUpheld(event annealing.Event) bool {
	switch event.EventType {
	case annealing.StartedAnnealing:
		amo.previousObjectiveValue = event.Annealer.ObservableExplorer().ObjectiveValue()
		return true
	case annealing.FinishedIteration:
		var expectedObjectiveValue float64
		if event.Annealer.ObservableExplorer().ChangeAccepted() {
			expectedObjectiveValue = amo.previousObjectiveValue + event.Annealer.ObservableExplorer().ChangeInObjectiveValue()
			amo.previousObjectiveValue = event.Annealer.ObservableExplorer().ObjectiveValue()
		} else {
			expectedObjectiveValue = amo.previousObjectiveValue
		}
		return expectedObjectiveValue == event.Annealer.ObservableExplorer().ObjectiveValue()
	default:
		return true
	}
}
