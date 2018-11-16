// Copyright (c) 2018 Australian Rivers Institute.

// observer package contains Observer implementations  that capture annealing events and log them to various
// destinations in potentially very different formats (as dictated by loggers passed to a logger at build-time)
package observer

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/explorer"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/observer/filters"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/wrapper"
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
func (l *AnnealingObserver) ObserveAnnealingEvent(event annealing.Event) {}

func wrapAnnealer(eventAnnealer annealing.Annealer) *wrapper.FormatWrapper {
	wrapper := newAnnealerWrapper()
	wrapper.Wrap(eventAnnealer)
	return wrapper
}

func wrapSolutionExplorer(explorer explorer.Explorer) *explorer.FormatWrapper {
	wrapper := newSolutionExplorerWrapper()
	wrapper.Wrap(explorer)
	return wrapper
}

func newAnnealerWrapper() *wrapper.FormatWrapper {
	wrapper := wrapper.FormatWrapper{
		MethodFormats: map[string]string{
			"Temperature":       "%0.4f",
			"CoolingFactor":     "%0.3f",
			"MaximumIterations": "%03d",
			"CurrentIteration":  "%03d",
		},
	}
	return &wrapper
}

func newSolutionExplorerWrapper() *explorer.FormatWrapper {
	wrapper := explorer.FormatWrapper{
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
