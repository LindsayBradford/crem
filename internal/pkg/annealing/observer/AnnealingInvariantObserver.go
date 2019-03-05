// Copyright (c) 2019 Australian Rivers Institute.

package observer

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/observer/filters"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/math"
	"github.com/LindsayBradford/crem/pkg/strings"
)

// TODO:  This cannot be run concurently as-is.  Needs a closer look at Annealer state cloning.

const decimalPrecisionRequired = 6

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

func (amo *AnnealingInvariantObserver) ObserveEvent(event observer.Event) {
	if amo.loopInvariantUpheld(event) {
		return
	}

	var builder strings.FluentBuilder
	builder.
		Add("Id [", event.Id(), "], ").
		Add("Event [", event.EventType.String(), "]: ").
		Add("Loop Invariant Broken.")
	amo.logHandler.LogAtLevel(AnnealerLogLevel, builder.String())
	panic(builder.String())
}

func (amo *AnnealingInvariantObserver) loopInvariantUpheld(event observer.Event) bool {
	if annealer, isAnnealer := event.Source().(annealing.Observable); isAnnealer {
		switch event.EventType {
		case observer.StartedAnnealing:
			amo.previousObjectiveValue = annealer.ObservableExplorer().ObjectiveValue()
			return true
		case observer.FinishedIteration:
			var expectedObjectiveValue float64
			if annealer.ObservableExplorer().ChangeAccepted() {
				expectedObjectiveValue = amo.previousObjectiveValue + annealer.ObservableExplorer().ChangeInObjectiveValue()
				amo.previousObjectiveValue = annealer.ObservableExplorer().ObjectiveValue()
			} else {
				expectedObjectiveValue = amo.previousObjectiveValue
			}

			roundedExpectedObjectiveValue := math.RoundFloat(expectedObjectiveValue, decimalPrecisionRequired)
			roundedActualObjectiveValue := math.RoundFloat(annealer.ObservableExplorer().ObjectiveValue(), decimalPrecisionRequired)

			invariantUpheld := roundedExpectedObjectiveValue == roundedActualObjectiveValue
			return invariantUpheld
		default:
			return true
		}
	}
	return true
}
