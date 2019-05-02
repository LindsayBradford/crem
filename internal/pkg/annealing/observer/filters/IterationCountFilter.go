// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package filters

import (
	"github.com/LindsayBradford/crem/internal/pkg/observer"
)

// IterationCountFilter modulates FinishedIteration Annealing Event instances at a rate of 1 every modulo
// events. StartedIteration events are completely filtered out. All other event types are allowed through to the LogHandler.
type IterationCountFilter struct {
	iterationModulo uint64
}

// WithModulo defines the modulo to apply against FinishedIteration Annealing Event instances.
func (m *IterationCountFilter) WithModulo(modulo uint64) *IterationCountFilter {
	m.iterationModulo = modulo
	return m
}

// ShouldFilter modulates only FinishedIteration Event instances, and fully filters out all StartedIteration
// events. Every modulo FinishedIteration events received, one is allowed through to the LogHandler.
// The very first and very last FinishedIteration events are exceptions, and are also not filtered.
func (m *IterationCountFilter) ShouldFilter(event observer.Event) bool {
	if event.EventType != observer.StartedIteration && event.EventType != observer.FinishedIteration {
		return allowThroughFilter
	}

	return m.ShouldFilterAnnealerSource(event)
}

func (m *IterationCountFilter) ShouldFilterAnnealerSource(event observer.Event) bool {
	currentIteration := event.Attribute("CurrentIteration").(uint64)
	if (event.EventType == observer.FinishedIteration || event.EventType == observer.StartedIteration) &&
		(eventOnnFirstOrLastIteration(event) || currentIteration%m.iterationModulo == 0) {
		return allowThroughFilter
	}
	return blockAtFilter
}
