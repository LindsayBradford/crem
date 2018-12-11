// Copyright (c) 2018 Australian Rivers Institute.

package main

import (
	"testing"

	"github.com/LindsayBradford/crem/cmd/cremengine/components"
	"github.com/LindsayBradford/crem/cmd/cremengine/components/scenario"
	configTesting "github.com/LindsayBradford/crem/internal/pkg/config/testing"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
)

func TestSedimentTransportAnnealerScenarioOneRun(t *testing.T) {
	context := configTesting.Context{
		Name:           "Single run of sediment transport annealer",
		T:              t,
		ConfigFilePath: "testdata/CatchmentConfig-OneRun.toml",
		Runner:         components.RunScenarioFromConfigFile,
	}

	scenario.LogHandler = loggers.DefaultTestingLogger
	context.VerifyScenarioRunViaConfigFile()
}

func TestSedimentTransportAnnealerScenarioBadInputs(t *testing.T) {
	context := configTesting.Context{
		Name:           "Attempted run of sediment transport annealer with bad inputs",
		T:              t,
		ConfigFilePath: "testdata/CatchmentConfig-BadInputs.toml",
		Runner:         components.RunScenarioFromConfigFile,
	}

	scenario.LogHandler = loggers.DefaultTestingLogger
	context.VerifyScenarioRunViaConfigFile()
}
