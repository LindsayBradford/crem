// (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

// modulators package supplies a number of logging modulators for managing the chattiness of loggers.
package filters

import . "github.com/LindsayBradford/crm/annealing/shared"

// LoggingFilter describes an interface to object that decides on how logging should be filtered
// (reduced in volume of entries logged).
type LoggingFilter interface {

	// ShouldFilter accepts an AnnealingEvent and decides whether it should be filtered (not logged).
	// This method returns true iff the logger is to ignore logging the supplied event.
	ShouldFilter(event AnnealingEvent) bool
}

// NullFilter is a default LoggingFilter that doesn't actually filter logging, allowing all events through.
type NullFilter struct{}

// ShouldFilter always returns false (do not filter the log of events)
func (nm *NullFilter) ShouldFilter(event AnnealingEvent) bool {
	return false
}
