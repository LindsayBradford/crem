// (c) 2018 Australian Rivers Institute.

package observer

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/observer/filters"
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

	if amo.logHandler.BeingDiscarded(AnnealingLogLevel) || amo.filter.ShouldFilter(event) {
		return
	}

	var builder strings.FluentBuilder
	builder.
		Add("Id [", event.Id(), "], ").
		Add("Event [", event.EventType.String(), "]: ")

	amo.observeEvent(event, &builder)

}

func (amo *AnnealingMessageObserver) observeEvent(event observer.Event, builder *strings.FluentBuilder) {
	switch event.EventType {
	case observer.StartedAnnealing:
		amo.stringifyEvent(event, builder)
	case observer.FinishedAnnealing:
		event.RemoveAttribute("Solution")
		fusedIterationsEvent := fuseIterationAttributes(event)
		amo.stringifyEvent(fusedIterationsEvent, builder)
	default:
		fusedIterationsEvent := fuseIterationAttributes(event)
		amo.stringifyEvent(fusedIterationsEvent, builder)
	}

	amo.logHandler.LogAtLevel(AnnealingLogLevel, builder.String())
}

const leftBrace = " ["
const rightBrace = "]"
const comma = ", "

func (amo *AnnealingMessageObserver) stringifyEvent(event observer.Event, builder *strings.FluentBuilder) {
	for index, attrib := range event.AllAttributes() {
		builder.Add(attrib.Name, leftBrace, format(event, attrib.Name), rightBrace)
		if index < len(event.AllAttributes())-1 {
			builder.Add(comma)
		}
	}
}

func format(event observer.Event, attributeName string) string {
	return defaultConverter.Convert(event.Attribute(attributeName))
}

func fuseIterationAttributes(event observer.Event) observer.Event {
	if !event.HasAttribute("CurrentIteration") && !event.HasAttribute("MaximumIterations") {
		return event
	}

	if event.Attribute("CurrentIteration").(uint64) == 0 {
		duplicateEvent := &event
		duplicateEvent.RemoveAttribute("CurrentIteration")
		duplicateEvent.RemoveAttribute("MaximumIterations")
		return *duplicateEvent
	}

	fusedIterationsValue := new(strings.FluentBuilder).
		Add(format(event, "CurrentIteration"), "/", format(event, "MaximumIterations")).String()

	duplicateEvent := &event
	duplicateEvent.RemoveAttribute("MaximumIterations")
	duplicateEvent.ReplaceAttribute("CurrentIteration", fusedIterationsValue)
	duplicateEvent.RenameAttribute("CurrentIteration", "Iteration")

	return *duplicateEvent
}
