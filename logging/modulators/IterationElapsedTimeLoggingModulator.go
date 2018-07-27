// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package modulators

import (
	. "github.com/LindsayBradford/crm/annealing/shared"
	"time"
)

// IterationElapsedTimeLoggingModulator is a LoggingModulator that will not modulate any AnnealingEvent types
// except StartedIteration & FinishedIteration. It completely filters out all StartedIteration events, and modulates
// FinishedIteration events at a rate of one event per every lapsed wait duration specified.
// The very first very last events are exceptions, and are also not modulated.
type IterationElapsedTimeLoggingModulator struct {
	waitDuration    time.Duration
	lastTimeAllowed time.Time
}

// WithWait sets the wait duration between allowing FinishedIteration AnnealingEvent instances through to a LogHandler.
func (m *IterationElapsedTimeLoggingModulator) WithWait(wait time.Duration) *IterationElapsedTimeLoggingModulator {
	m.waitDuration = wait
	return m
}

// ShouldModulate returns true for most FinishedIteration AnnealingEvent instances. Those allowed through to the logger
// are either 1) the very first or very last event, or 2) the closest FinishedIteration event to have
// occurred after the wait duration has passed since the last previous event allowed through.
func (m *IterationElapsedTimeLoggingModulator) ShouldModulate(event AnnealingEvent) bool {
	if event.EventType != StartedIteration && event.EventType != FinishedIteration {
		return false
	}

	annealer := event.Annealer
	if event.EventType == FinishedIteration &&
		(annealer.CurrentIteration() == 1 || annealer.CurrentIteration() == annealer.MaxIterations()) {
		m.lastTimeAllowed = time.Now()
		return false
	}

	if event.EventType == FinishedIteration && time.Now().Sub(m.lastTimeAllowed) >= m.waitDuration {
		m.lastTimeAllowed = m.lastTimeAllowed.Add(m.waitDuration)
		return false
	}

	return true
}
