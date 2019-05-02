// Copyright (c) 2019 Australian Rivers Institute.

package observer

import (
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
	switch event.EventType {
	case observer.StartedAnnealing:
		amo.previousObjectiveValue = event.Attribute("ObjectiveValue").(float64)
		return true
	case observer.FinishedIteration:
		actualObjectiveValue := event.Attribute("ObjectiveValue").(float64)
		changeInObjectiveValue := event.Attribute("ChangeInObjectiveValue").(float64)
		changeAccepted := event.Attribute("ChangeAccepted").(bool)

		var expectedObjectiveValue float64

		if changeAccepted {
			expectedObjectiveValue = amo.previousObjectiveValue + changeInObjectiveValue
			amo.previousObjectiveValue = actualObjectiveValue
		} else {
			expectedObjectiveValue = amo.previousObjectiveValue
		}

		roundedExpectedObjectiveValue := math.RoundFloat(expectedObjectiveValue, decimalPrecisionRequired)
		roundedActualObjectiveValue := math.RoundFloat(actualObjectiveValue, decimalPrecisionRequired)

		invariantUpheld := roundedExpectedObjectiveValue == roundedActualObjectiveValue
		return invariantUpheld
	default:
		return true
	}
}
