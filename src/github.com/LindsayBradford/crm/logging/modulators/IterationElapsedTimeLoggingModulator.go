// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package modulators

import (
	. "github.com/LindsayBradford/crm/annealing/shared"
	"time"
)

// IterationElapsedTimeLoggingModulator is a LoggingModulator that will not modulate any AnnealingEvent types
// except STARTED_ITERATION. This Modulator will modulate STARTED_ITERATION events at a rate of one event per every
// elapsed wait duration specified. The very first very last events are exceptions, and are also not modulated.
type IterationElapsedTimeLoggingModulator struct {
	waitDuration    time.Duration
	lastTimeAllowed time.Time
}

// WithWait sets the wait duration between allowing STARTED_ITERATION AnnealingEvent instances through to a LogHandler.
func (this *IterationElapsedTimeLoggingModulator) WithWait(wait time.Duration) *IterationElapsedTimeLoggingModulator {
	this.waitDuration = wait
	return this
}

// ShouldModulate returns true for most STARTED_ITERATION AnnealingEvent instances. Those allowed through to the logger
// are either 1) the very first or very last STARTED_ITERATION event, or 2) the closest STARTED_ITERATION event to have
// occurred after the wait duration has passed after the previous was allowed through.
func (this *IterationElapsedTimeLoggingModulator) ShouldModulate(event AnnealingEvent) bool {
	if event.EventType != STARTED_ITERATION && event.EventType != FINISHED_ITERATION {
		return false
	}

	annealer := event.Annealer
	if event.EventType == FINISHED_ITERATION && (annealer.CurrentIteration() == 1 || annealer.CurrentIteration() == annealer.MaxIterations()) {
		this.lastTimeAllowed = time.Now()
		return false
	}

	if event.EventType == FINISHED_ITERATION && time.Now().Sub(this.lastTimeAllowed) >= this.waitDuration {
		this.lastTimeAllowed = this.lastTimeAllowed.Add(this.waitDuration)
		return false
	}

	return true
}
