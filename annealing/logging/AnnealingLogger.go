// Copyright (c) 2018 Australian Rivers Institute.

// logging package contains AnnealingObserver implementations  that capture annealing events and log them to various
// destinations in potentially very different formats (as dictated by handlers passed to a logger at build-time)
package logging

import (
	. "github.com/LindsayBradford/crem/annealing/shared"
	. "github.com/LindsayBradford/crem/annealing/solution"
	. "github.com/LindsayBradford/crem/logging/filters"
	. "github.com/LindsayBradford/crem/logging/handlers"
	. "github.com/LindsayBradford/crem/logging/shared"
)

const AnnealerLogLevel LogLevel = "Annealer"

// AnnealingLogger is a base-implementation of an annealing logger.  It has a logHandler, but deliberately
// drops any AnnealingEvents received.
type AnnealingLogger struct {
	logHandler LogHandler
	filter     LoggingFilter
}

// Allows for the receipt of AnnealingEvent instances, but deliberately takes no action in logging those events.
func (l *AnnealingLogger) ObserveAnnealingEvent(event AnnealingEvent) {}

func wrapAnnealer(eventAnnealer Annealer) *AnnealerFormatWrapper {
	wrapper := newAnnealerWrapper()
	wrapper.Wrap(eventAnnealer)
	return wrapper
}

func wrapSolutionExplorer(explorer Explorer) *ExplorerFormatWrapper {
	wrapper := newSolutionExplorerWrapper()
	wrapper.Wrap(explorer)
	return wrapper
}

func newAnnealerWrapper() *AnnealerFormatWrapper {
	wrapper := AnnealerFormatWrapper{
		MethodFormats: map[string]string{
			"Temperature":      "%0.4f",
			"CoolingFactor":    "%0.3f",
			"MaxIterations":    "%03d",
			"CurrentIteration": "%03d",
		},
	}
	return &wrapper
}

func newSolutionExplorerWrapper() *ExplorerFormatWrapper {
	wrapper := ExplorerFormatWrapper{
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
