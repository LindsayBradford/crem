// +build windows
// Copyright (c) 2018 Australian Rivers Institute.

package main

import (
	"github.com/LindsayBradford/crem/cmd/cremengine/components/scenario"
	appTesting "github.com/LindsayBradford/crem/internal/pkg/config/testing"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"

	"testing"
)

const defaultCatchmentModelAnnealerTimeout = 10

func TestCremEngine_ScenarioOneRunBlackbox_ExitWithSuccess(t *testing.T) {
	context := appTesting.BlackboxTestingContext{
		T:                 t,
		Name:              "Single run of CremEngine CatchmentModel",
		ExecutablePath:    cremExceutablePath,
		TimeoutSeconds:    defaultCatchmentModelAnnealerTimeout,
		ConfigFilePath:    "testdata/CatchmentConfig-BlackBox-OneRun.toml",
		ExpectedErrorCode: appTesting.WithSuccess,
	}

	appTesting.TestExecutableAgainstConfigFile(context)
}

func TestCremEngine_ScenarioBadInputsBlackbox_ExitWithError(t *testing.T) {
	context := appTesting.BlackboxTestingContext{
		T:                 t,
		Name:              "Attempted run of CremEngine CatchmentModel with bad inputs",
		ExecutablePath:    cremExceutablePath,
		TimeoutSeconds:    defaultCatchmentModelAnnealerTimeout,
		ConfigFilePath:    "testdata/CatchmentConfig-BadInputs.toml",
		ExpectedErrorCode: appTesting.WithFailure,
	}

	appTesting.TestExecutableAgainstConfigFile(context)
}

func TestCremEngine_ScenarioOneRunWhitebox_ExitWithSuccess(t *testing.T) {
	context := appTesting.WhiteboxTestingContext{
		Name: "Single run of catchment model annealer",
		T:    t,
		//ConfigFilePath:    "testdata/CatchmentConfig-BadInputs.toml",
		ConfigFilePath: "testdata/CatchmentConfig-WhiteBox-OneRun.toml",
		Runner:         scenario.RunExcelCompatibleScenarioFromConfigFile,
	}

	scenario.LogHandler = loggers.DefaultTestingLogger
	context.VerifyGoroutineScenarioRunViaConfigFileDoesNotPanic()
}
