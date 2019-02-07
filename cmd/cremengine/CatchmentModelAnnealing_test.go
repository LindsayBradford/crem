// +build windows
// Copyright (c) 2018 Australian Rivers Institute.

package main

import (
	"testing"
)

func TestCremEngine_ScenarioOneRun_ExitWithSuccess(t *testing.T) {
	context := BinaryTestingContext{
		Name:              "Single run of CremEngine CatchmentModel",
		T:                 t,
		ConfigFilePath:    "testdata/CatchmentConfig-OneRun.toml",
		ExpectedErrorCode: withSuccess,
	}

	testCremExecutableAgainstConfigFile(context)
}

func TestCremEngine_ScenarioBadInputs_ExitWithError(t *testing.T) {
	context := BinaryTestingContext{
		Name:              "Attempted run of CremEngine CatchmentModel with bad inputs",
		T:                 t,
		ConfigFilePath:    "testdata/CatchmentConfig-BadInputs.toml",
		ExpectedErrorCode: withFailure,
	}

	testCremExecutableAgainstConfigFile(context)
}
