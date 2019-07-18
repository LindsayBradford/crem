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

func TestCREMExplorer_BlackBox_ExitWithSuccess(t *testing.T) {
	context := appTesting.BlackboxTestingContext{
		T:                 t,
		Name:              "CREMEngine - Black Box",
		ExecutablePath:    exceutablePath,
		TimeoutSeconds:    defaultCatchmentModelAnnealerTimeout,
		ConfigFilePath:    "testdata/TestCREMEngine-BlackBox.toml",
		ExpectedErrorCode: appTesting.WithSuccess,
	}

	appTesting.TestExecutableAgainstConfigFile(context)
}

func TestCREMExplorer_BlackBox_ExitWithError(t *testing.T) {
	context := appTesting.BlackboxTestingContext{
		T:                 t,
		Name:              "CREMEngine - Black Box with bad inputs",
		ExecutablePath:    exceutablePath,
		TimeoutSeconds:    defaultCatchmentModelAnnealerTimeout,
		ConfigFilePath:    "testdata/TestCREMEngine-BadInputs.toml",
		ExpectedErrorCode: appTesting.WithFailure,
	}

	appTesting.TestExecutableAgainstConfigFile(context)
}

func TestCREMExplorer_WhiteBox_ExitWithSuccess(t *testing.T) {
	context := appTesting.WhiteboxTestingContext{
		Name:           "CREMEngine - White Box",
		T:              t,
		ConfigFilePath: "testdata/TestCREMEngine-WhiteBox.toml",
		Runner:         bootstrap.RunExcelCompatibleScenarioFromConfigFile,
	}

	bootstrap.LogHandler = loggers.DefaultTestingLogger
	context.VerifyGoroutineScenarioRunViaConfigFileDoesNotPanic()
}
