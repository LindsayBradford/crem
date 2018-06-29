// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package modulators

import . "github.com/LindsayBradford/crm/annealing/shared"

// IterationModuloLoggingModulator modulates STARTED_ITERATION Annealing Event instances at a rate of 1 every modulo
// events. All other event types are allowed through to the LogHandler.
type IterationModuloLoggingModulator struct {
	iterationModulo uint
}

// WithModulo defines the modulo to apply against STARTED_ITERATION Annealing Event instances.
func (this *IterationModuloLoggingModulator) WithModulo(modulo uint) *IterationModuloLoggingModulator {
	this.iterationModulo = modulo
	return this
}

// ShouldModulate modulates only STARTED_ITERATION AnnealingEvent instances. Every modulo STARTED_ITERATION events
// received, one is allowed through to the LogHandler. The very first very last STARTED_ITERATION events are
// exceptions, and are also not modulated.
func (this *IterationModuloLoggingModulator) ShouldModulate(event AnnealingEvent) bool {
	if event.EventType != STARTED_ITERATION && event.EventType != FINISHED_ITERATION {
		return false
	}

	annealer := event.Annealer
	if event.EventType == FINISHED_ITERATION && (annealer.CurrentIteration() == 1 || annealer.CurrentIteration() == annealer.MaxIterations() ||
		annealer.CurrentIteration()%this.iterationModulo == 0) {
		return false
	}
	return true
}