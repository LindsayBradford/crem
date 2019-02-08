// +build windows
// Copyright (c) 2018 Australian Rivers Institute.

package main

import (
	appTesting "github.com/LindsayBradford/crem/internal/pkg/config/testing"
	"testing"
)

func TestCremEngine_ScenarioOneRun_ExitWithSuccess(t *testing.T) {
	context := appTesting.BinaryTestingContext{
		T:                 t,
		Name:              "Single run of CremEngine CatchmentModel",
		ExecutablePath:    cremExceutablePath,
		TimeoutSeconds:    defaultDumbAnnealerTimeout,
		ConfigFilePath:    "testdata/CatchmentConfig-OneRun.toml",
		ExpectedErrorCode: appTesting.WithSuccess,
	}

	appTesting.TestExecutableAgainstConfigFile(context)
}

func TestCremEngine_ScenarioBadInputs_ExitWithError(t *testing.T) {
	context := appTesting.BinaryTestingContext{
		T:                 t,
		Name:              "Attempted run of CremEngine CatchmentModel with bad inputs",
		ExecutablePath:    cremExceutablePath,
		TimeoutSeconds:    defaultDumbAnnealerTimeout,
		ConfigFilePath:    "testdata/CatchmentConfig-BadInputs.toml",
		ExpectedErrorCode: appTesting.WithFailure,
	}

	appTesting.TestExecutableAgainstConfigFile(context)
}
