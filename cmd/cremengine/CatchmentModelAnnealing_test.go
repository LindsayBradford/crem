// +build windows
// Copyright (c) 2018 Australian Rivers Institute.

package main

import (
	"testing"

	"github.com/LindsayBradford/crem/cmd/cremengine/components"
	"github.com/LindsayBradford/crem/cmd/cremengine/components/scenario"
	configTesting "github.com/LindsayBradford/crem/internal/pkg/config/testing"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
)

func TestCatchmentModelAnnealerScenarioOneRun(t *testing.T) {
	context := configTesting.Context{
		Name:           "Single run of catchment model annealer",
		T:              t,
		ConfigFilePath: "testdata/CatchmentConfig-OneRun.toml",
		Runner:         components.RunExcelCompatibleScenarioFromConfigFile,
	}

	scenario.LogHandler = loggers.DefaultTestingLogger
	context.VerifyGoroutineScenarioRunViaConfigFileDoesNotPanic()
}

func TestCatchmentModelAnnealerScenarioBadInputs(t *testing.T) {
	context := configTesting.Context{
		Name:           "Attempted run of catchment model annealer with bad inputs",
		T:              t,
		ConfigFilePath: "testdata/CatchmentConfig-BadInputs.toml",
		Runner:         components.RunExcelCompatibleScenarioFromConfigFile,
	}

	scenario.LogHandler = loggers.DefaultTestingLogger
	context.VerifyGoroutineScenarioRunViaConfigFileDoesNotPanic()
}
