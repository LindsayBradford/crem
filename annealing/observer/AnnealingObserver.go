// Copyright (c) 2018 Australian Rivers Institute.

// observer package contains Observer implementations  that capture annealing events and log them to various
// destinations in potentially very different formats (as dictated by handlers passed to a logger at build-time)
package observer

import (
	. "github.com/LindsayBradford/crem/annealing"
	. "github.com/LindsayBradford/crem/annealing/explorer"
	"github.com/LindsayBradford/crem/annealing/wrapper"
	. "github.com/LindsayBradford/crem/logging/filters"
	. "github.com/LindsayBradford/crem/logging/handlers"
	. "github.com/LindsayBradford/crem/logging/shared"
)

const AnnealerLogLevel LogLevel = "Annealer"

// AnnealingObserver is a base-implementation of an annealing logger.  It has a logHandler, but deliberately
// drops any AnnealingEvents received.
type AnnealingObserver struct {
	logHandler LogHandler
	filter     LoggingFilter
}

// Allows for the receipt of Event instances, but deliberately takes no action in observer those events.
func (l *AnnealingObserver) ObserveAnnealingEvent(event Event) {}

func wrapAnnealer(eventAnnealer Annealer) *wrapper.FormatWrapper {
	wrapper := newAnnealerWrapper()
	wrapper.Wrap(eventAnnealer)
	return wrapper
}

func wrapSolutionExplorer(explorer Explorer) *FormatWrapper {
	wrapper := newSolutionExplorerWrapper()
	wrapper.Wrap(explorer)
	return wrapper
}

func newAnnealerWrapper() *wrapper.FormatWrapper {
	wrapper := wrapper.FormatWrapper{
		MethodFormats: map[string]string{
			"Temperature":      "%0.4f",
			"CoolingFactor":    "%0.3f",
			"MaxIterations":    "%03d",
			"CurrentIteration": "%03d",
		},
	}
	return &wrapper
}

func newSolutionExplorerWrapper() *FormatWrapper {
	wrapper := FormatWrapper{
		MethodFormats: map[string]string{
			"ObjectiveValue":         "%0.4f",
			"ChangeInObjectiveValue": "%0.4f",
			"ChangeIsDesirable":      "%t",
			"AcceptanceProbability":  "%0.6f",
			"ChangeAccepted":         "%t",
		},
	}
	return &wrapper
}
