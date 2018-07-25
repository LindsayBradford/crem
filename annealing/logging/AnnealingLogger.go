// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

// logging package contains AnnealingObserver implementations  that capture annealing events and log them to various
// destinations in potentially very different formats (as dictated by handlers passed to a logger at build-time)
package logging

import (
	. "github.com/LindsayBradford/crm/annealing/shared"
	. "github.com/LindsayBradford/crm/annealing/objectives"
	. "github.com/LindsayBradford/crm/logging/handlers"
	. "github.com/LindsayBradford/crm/logging/modulators"
	. "github.com/LindsayBradford/crm/logging/shared"
)

const ANNEALER LogLevel = "Annealer"

// AnnealingLogger is a base-implementation of an annealing logger.  It has a logHandler, but deliberately
// drops any AnnealingEvents received.
type AnnealingLogger struct {
	logHandler LogHandler
	modulator  LoggingModulator
}

// Allows for the receipt of AnnealingEvent instances, but deliberately takes no action in logging those events.
func (this *AnnealingLogger) ObserveAnnealingEvent(event AnnealingEvent) {}

func wrapAnnealer(eventAnnealer Annealer) *AnnealerStateFormatWrapper {
	wrapper := newAnnealingWrapper()
	wrapper.Wrap(eventAnnealer)
	return wrapper
}

func wrapObjectiveManager(eventObjectiveManager ObjectiveManager) *ObjectiveManagerStateFormatWrapper {
	wrapper := newObjectiveManagerWrapper()
	wrapper.Wrap(eventObjectiveManager)
	return wrapper
}

func newAnnealingWrapper() *AnnealerStateFormatWrapper  {
	annealingWrapper := AnnealerStateFormatWrapper{
		MethodFormats: map[string]string{
			"Temperature":      "%0.4f",
			"CoolingFactor":    "%0.3f",
			"MaxIterations":    "%03d",
			"CurrentIteration": "%03d",
		},
	}
	return &annealingWrapper
}

func newObjectiveManagerWrapper() *ObjectiveManagerStateFormatWrapper  {
	objectiveManagerWrapper := ObjectiveManagerStateFormatWrapper{
		MethodFormats: map[string]string{
			"ObjectiveValue":         "%0.4f",
			"ChangeInObjectiveValue": "%0.4f",
			"ChangeIsDesirable":      "%t",
			"AcceptanceProbability":  "%0.6f",
			"ChangeAccepted":         "%t",
		},
	}
	return &objectiveManagerWrapper
}