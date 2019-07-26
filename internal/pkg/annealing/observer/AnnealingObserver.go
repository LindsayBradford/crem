// Copyright (c) 2018 Australian Rivers Institute.

// observer package contains Reporting implementations  that capture annealing events and log them to various
// destinations in potentially very different formats (as dictated by loggers passed to a logger at build-time)
package observer

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/observer/filters"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/pkg/logging"
)

const AnnealerLogLevel logging.Level = "Annealer"

// AnnealingObserver is a base-implementation of an annealing logger.  It has a logHandler, but deliberately
// drops any AnnealingEvents received.
type AnnealingObserver struct {
	logHandler logging.Logger
	filter     filters.Filter
}

// Allows for the receipt of Event instances, but deliberately takes no action in observer those events.
func (l *AnnealingObserver) ObserveAnnealingEvent(event observer.Event) {}
