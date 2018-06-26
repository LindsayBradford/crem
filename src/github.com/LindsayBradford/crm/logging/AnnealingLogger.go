// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

// logging package contains AnnealingObserver implementations  that capture annealing events and log them to various
// destinations in potentially very different formats (as dictated by handlers passed to a logger at build-time)
package logging

import (
	. "github.com/LindsayBradford/crm/annealing"
	. "github.com/LindsayBradford/crm/logging/handlers"
)

// AnnealingLogger is a base-implementation of an annealing logger.  It has a logHandler, but deliberately
// drops any AnnealingEvents received.
type AnnealingLogger struct {
	logHandler LogHandler
}

// Allows for the receipt of AnnelingEvent instances, but deliberately takes no action in logging those events.
func (this *AnnealingLogger) ObserveAnnealingEvent(event AnnealingEvent) {}