// (c) 2018 Australian Rivers Institute.

package observer

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/observer/filters"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	assert "github.com/LindsayBradford/crem/pkg/assert/debug"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/strings"
)

var (
	defaultConverter = strings.NewConverter().Localised().WithFloatingPointPrecision(6).PaddingZeros()
)

// AnnealingMessageObserver produces a stream of human-friendly, free-form text log entries from any observed
// Event instances received.
type AnnealingMessageObserver struct {
	AnnealingObserver
	invariantObserver *AnnealingInvariantObserver
}

func (amo *AnnealingMessageObserver) WithLogHandler(handler logging.Logger) *AnnealingMessageObserver {
	amo.logHandler = handler
	return amo
}

func (amo *AnnealingMessageObserver) WithFilter(filter filters.Filter) *AnnealingMessageObserver {
	amo.filter = filter
	return amo
}

func (amo *AnnealingMessageObserver) WithLoopInvariantObserver(watchLoopInvariant bool) *AnnealingMessageObserver {
	if watchLoopInvariant {
		assert.That(amo.logHandler != nil)
		amo.invariantObserver = new(AnnealingInvariantObserver).WithLogHandler(amo.logHandler)
	}
	return amo
}

// ObserveEvent captures and converts Event instances into free-form text strings that it
// then passes onto its relevant Logger as an Info call.
func (amo *AnnealingMessageObserver) ObserveEvent(event observer.Event) {
	if amo.invariantObserver != nil {
		amo.invariantObserver.ObserveEvent(event)
	}

	if amo.logHandler.BeingDiscarded(AnnealerLogLevel) || amo.filter.ShouldFilter(event) {
		return
	}

	var builder strings.FluentBuilder
	builder.
		Add("Id [", event.Id(), "], ").
		Add("Event [", event.EventType.String(), "]: ")

	if event.EventType.IsAnnealingState() {
		amo.observeAnnealerEvent(event, &builder)
	} else {
		amo.observeEvent(event, &builder)
	}
}

func (amo *AnnealingMessageObserver) observeAnnealerEvent(event observer.Event, builder *strings.FluentBuilder) {
	switch event.EventType {
	case observer.StartedAnnealing:
		builder.
			Add("Maximum Iterations [", format(event, "MaximumIterations"), "], ").
			Add("Temperature [", format(event, "Temperature"), "], ").
			Add("Cooling Factor [", format(event, "CoolingFactor"), "]")
		if event.HasAttribute("ObjectiveValue") {
			builder.Add(", Objective value [", format(event, "ObjectiveValue"), "]")
		}
		if event.HasAttribute("ArchiveSize") {
			builder.Add(", Archive size [", format(event, "ArchiveSize"), "]")
		}
	case observer.StartedIteration:
		builder.
			Add("Iteration [", format(event, "CurrentIteration"), "/", format(event, "MaximumIterations"), "], ").
			Add("Temperature [", format(event, "Temperature"), "]")
		if event.HasAttribute("ObjectiveValue") {
			builder.Add(", Objective value [", format(event, "ObjectiveValue"), "]")
		}
	case observer.FinishedIteration:
		builder.
			Add("Iteration [", format(event, "CurrentIteration"), "/", format(event, "MaximumIterations"), "], ").
			Add("Desirable? [", format(event, "ChangeIsDesirable"), "], ").
			Add("Acceptance Probability [", format(event, "AcceptanceProbability"), "], ").
			Add("Accepted? [", format(event, "ChangeAccepted"), "]")
		if event.HasAttribute("ObjectiveValue") {
			builder.Add(", Change [", format(event, "ChangeInObjectiveValue"), "], ").
				Add("Objective value [", format(event, "ObjectiveValue"), "]")
		}
		if event.HasAttribute("ArchiveSize") {
			builder.Add(", Archive size [", format(event, "ArchiveSize"), "], ").
				Add("Archive result [", format(event, "ArchiveResult"), "]")
		}

	case observer.FinishedAnnealing:
		builder.
			Add("Iteration [", format(event, "CurrentIteration"), "/", format(event, "CurrentIteration"), "], ").
			Add("Temperature [", format(event, "Temperature"), "]")
		if event.HasAttribute("ObjectiveValue") {
			builder.Add(", Objective value [", format(event, "ObjectiveValue"), "]")
		}
	default:
		// deliberately does nothing extra
	}

	amo.logHandler.LogAtLevel(AnnealerLogLevel, builder.String())
}

func (amo *AnnealingMessageObserver) observeEvent(event observer.Event, builder *strings.FluentBuilder) {
	switch event.EventType {
	case observer.Note:
		builder.Add("[", event.Note(), "]")
	case observer.ManagementAction:
		builder.
			Add("Type [", format(event, "Type"), "], ").
			Add("Planning Unit [", format(event, "PlanningUnit"), "], ").
			Add("Active [", format(event, "IsActive"), "]")

		if event.HasNote() {
			builder.Add(", Note [", event.Note(), "]")
		}
	case observer.DecisionVariable:
		builder.
			Add("Name [", format(event, "Name"), "], ").
			Add("Value [", format(event, "Value"), "], ")

		if event.HasNote() {
			builder.Add(", Note [", event.Note(), "]")
		}
	default:
		// deliberately does nothing extra
	}

	amo.logHandler.LogAtLevel(model.LogLevel, builder.String())
}

func format(event observer.Event, attributeName string) string {
	return defaultConverter.Convert(event.Attribute(attributeName))
}
