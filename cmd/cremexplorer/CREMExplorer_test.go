// +build windows
// Copyright (c) 2018 Australian Rivers Institute.

package main

import (
	"github.com/LindsayBradford/crem/cmd/cremexplorer/bootstrap"
	configTesting "github.com/LindsayBradford/crem/internal/pkg/config/testing"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"

	"testing"
)

const defaultCatchmentModelAnnealerTimeout = 10

func TestCREMExplorer_BlackBox_ExitWithSuccess(t *testing.T) {
	context := configTesting.BlackboxTestingContext{
		T:                 t,
		Name:              "CREMExplorer - Black Box",
		ExecutablePath:    exceutablePath,
		TimeoutSeconds:    defaultCatchmentModelAnnealerTimeout,
		ConfigFilePath:    "testdata/TestCREMExplorer-BlackBox.toml",
		ExpectedErrorCode: configTesting.WithSuccess,
	}

	configTesting.TestExecutableAgainstConfigFile(context)
}

func TestCREMExplorer_BlackBox_ExitWithError(t *testing.T) {
	context := configTesting.BlackboxTestingContext{
		T:                 t,
		Name:              "CREMExplorer - Black Box with bad inputs",
		ExecutablePath:    exceutablePath,
		TimeoutSeconds:    defaultCatchmentModelAnnealerTimeout,
		ConfigFilePath:    "testdata/TestCREMExplorer-BadInputs.toml",
		ExpectedErrorCode: configTesting.WithFailure,
	}

	configTesting.TestExecutableAgainstConfigFile(context)
}

func TestCREMExplorer_Kirkpatrick_WhiteBox_ExitWithSuccess(t *testing.T) {
	context := configTesting.WhiteboxTestingContext{
		Name:           "CREMExplorer - Kirkpatrick - White Box",
		T:              t,
		ConfigFilePath: "testdata/TestCREMExplorer-Kirkpatrick-WhiteBox.toml",
		Runner:         bootstrap.RunExcelCompatibleScenarioFromConfigFile,
	}

	bootstrap.LogHandler = loggers.DefaultTestingLogger
	context.VerifyGoroutineScenarioRunViaConfigFileDoesNotPanic()
}

func TestCREMExplorer_Suppapitnarm_WhiteBox_ExitWithSuccess(t *testing.T) {
	context := configTesting.WhiteboxTestingContext{
		Name:           "CREMExplorer - Suppapitnarm - White Box",
		T:              t,
		ConfigFilePath: "testdata/TestCREMExplorer-Suppapitnarm-WhiteBox.toml",
		Runner:         bootstrap.RunExcelCompatibleScenarioFromConfigFile,
	}

	bootstrap.LogHandler = loggers.DefaultTestingLogger
	context.VerifyGoroutineScenarioRunViaConfigFileDoesNotPanic()
}
