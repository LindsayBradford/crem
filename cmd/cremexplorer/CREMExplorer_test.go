// +build windows
// Copyright (c) 2018 Australian Rivers Institute.

package main

import (
	"github.com/LindsayBradford/crem/cmd/cremexplorer/bootstrap"
	appTesting "github.com/LindsayBradford/crem/internal/pkg/config/testing"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"

	"testing"
)

const defaultCatchmentModelAnnealerTimeout = 10

func TestCremEngine_ScenarioOneRunBlackbox_ExitWithSuccess(t *testing.T) {
	context := appTesting.BlackboxTestingContext{
		T:                 t,
		Name:              "Single run of CremEngine CatchmentModel",
		ExecutablePath:    exceutablePath,
		TimeoutSeconds:    defaultCatchmentModelAnnealerTimeout,
		ConfigFilePath:    "testdata/TestCREMEngine-BlackBox.toml",
		ExpectedErrorCode: appTesting.WithSuccess,
	}

	appTesting.TestExecutableAgainstConfigFile(context)
}

func TestCremEngine_ScenarioBadInputsBlackbox_ExitWithError(t *testing.T) {
	context := appTesting.BlackboxTestingContext{
		T:                 t,
		Name:              "Attempted run of CremEngine CatchmentModel with bad inputs",
		ExecutablePath:    exceutablePath,
		TimeoutSeconds:    defaultCatchmentModelAnnealerTimeout,
		ConfigFilePath:    "testdata/TestCREMEngine-BadInputs.toml",
		ExpectedErrorCode: appTesting.WithFailure,
	}

	appTesting.TestExecutableAgainstConfigFile(context)
}

func TestCremEngine_ScenarioOneRunWhitebox_ExitWithSuccess(t *testing.T) {
	context := appTesting.WhiteboxTestingContext{
		Name:           "Single run of catchment model annealer",
		T:              t,
		ConfigFilePath: "testdata/TestCREMEngine-WhiteBox.toml",
		Runner:         bootstrap.RunExcelCompatibleScenarioFromConfigFile,
	}

	bootstrap.LogHandler = loggers.DefaultTestingLogger
	context.VerifyGoroutineScenarioRunViaConfigFileDoesNotPanic()
}
