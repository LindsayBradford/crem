// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package filters

import . "github.com/LindsayBradford/crem/annealing/shared"

// IterationCountLoggingFilter modulates FinishedIteration Annealing Event instances at a rate of 1 every modulo
// events. StartedIteration events are completely filtered out. All other event types are allowed through to the LogHandler.
type IterationCountLoggingFilter struct {
	iterationModulo uint64
}

// WithModulo defines the modulo to apply against FinishedIteration Annealing Event instances.
func (m *IterationCountLoggingFilter) WithModulo(modulo uint64) *IterationCountLoggingFilter {
	m.iterationModulo = modulo
	return m
}

// ShouldFilter modulates only FinishedIteration AnnealingEvent instances, and fully filters out all StartedIteration
// events. Every modulo FinishedIteration events received, one is allowed through to the LogHandler.
// The very first and very last FinishedIteration events are exceptions, and are also not filtered.
func (m *IterationCountLoggingFilter) ShouldFilter(event AnnealingEvent) bool {
	if event.EventType != StartedIteration && event.EventType != FinishedIteration {
		return false
	}

	annealer := event.Annealer
	if event.EventType == FinishedIteration &&
		(annealer.CurrentIteration() == 1 || annealer.CurrentIteration() == annealer.MaxIterations() ||
			annealer.CurrentIteration()%m.iterationModulo == 0) {
		return false
	}

	return true
}
