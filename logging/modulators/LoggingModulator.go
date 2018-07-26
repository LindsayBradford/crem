// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

// modulators package supplies a number of logging modulators for managing the chattiness of loggers.
package modulators

import . "github.com/LindsayBradford/crm/annealing/shared"

// LoggingModulator describes an interface to object that decides on how logging should be modulated
// (reduced in volume of entries logged).
type LoggingModulator interface {

	// ShouldModulate accepts an AnnealingEvent and decides whether it should be modulated (not logged).
	// This method returns true iff the logger is to modulate/ignore logging the supplied event.
	ShouldModulate(event AnnealingEvent) bool
}

// NullModulator is a default LoggingModulator that doesn't actually modulate logging (lets all events through).
type NullModulator struct{}

// ShouldModulate always returns false (do not modulate the log for event)
func (this *NullModulator) ShouldModulate(event AnnealingEvent) bool {
	return false
}
