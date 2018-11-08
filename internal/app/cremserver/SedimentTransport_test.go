// Copyright (c) 2018 Australian Rivers Institute.

package main

import (
	"testing"

	configTesting "github.com/LindsayBradford/crem/config/testing"
	"github.com/LindsayBradford/crem/internal/app/cremserver/components"
	"github.com/LindsayBradford/crem/internal/app/cremserver/components/scenario"
	"github.com/LindsayBradford/crem/logging/loggers"
)

func TestSedimentTransportAnnealerScenarioOneRun(t *testing.T) {
	context := configTesting.TestingContext{
		Name:           "Single run of sediment transport annealer",
		T:              t,
		ConfigFilePath: "testdata/SedimentTransportTestConfig-OneRun.toml",
		Runner:         components.RunScenarioFromConfigFile,
	}

	scenario.LogHandler = loggers.DefaultTestingLogger
	context.VerifyScenarioRunViaConfigFileDoesNotPanic()
}

func TestSedimentTransportAnnealerScenarioBadInputs(t *testing.T) {
	context := configTesting.TestingContext{
		Name:           "Attempted run of sediment transport annealer with bad inputs",
		T:              t,
		ConfigFilePath: "testdata/SedimentTransportTestConfig-BadInputs.toml",
		Runner:         components.RunScenarioFromConfigFile,
	}

	scenario.LogHandler = loggers.DefaultTestingLogger
	context.VerifyScenarioRunViaConfigFileDoesNotPanic()
}
