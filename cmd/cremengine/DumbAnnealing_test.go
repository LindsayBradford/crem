// +build windows
// Copyright (c) 2019 Australian Rivers Institute.

package main

import (
	appTesting "github.com/LindsayBradford/crem/internal/pkg/config/testing"
	"testing"
)

const defaultDumbAnnealerTimeout = 5

func TestDumbAnnealerIntegrationOneRun(t *testing.T) {
	context := appTesting.BinaryTestingContext{
		T:                 t,
		Name:              "Single run of Dumb Annealer",
		ExecutablePath:    cremExceutablePath,
		TimeoutSeconds:    defaultDumbAnnealerTimeout,
		ConfigFilePath:    "testdata/DumbAnnealerTestConfig-OneRun.toml",
		ExpectedErrorCode: appTesting.WithSuccess,
	}

	appTesting.TestExecutableAgainstConfigFile(context)
}

func TestDumbAnnealerIntegrationThreeRunsSequentially(t *testing.T) {
	context := appTesting.BinaryTestingContext{
		T:                 t,
		Name:              "Three sequential runs of Dumb Annealer",
		ExecutablePath:    cremExceutablePath,
		TimeoutSeconds:    defaultDumbAnnealerTimeout,
		ConfigFilePath:    "testdata/DumbAnnealerTestConfig-ThreeRunsSequentially.toml",
		ExpectedErrorCode: appTesting.WithSuccess,
	}

	appTesting.TestExecutableAgainstConfigFile(context)
}

func TestDumbAnnealerIntegrationThreeRunsConcurrently(t *testing.T) {
	context := appTesting.BinaryTestingContext{
		T:                 t,
		Name:              "Three concurrent runs of Dumb Annealer",
		ExecutablePath:    cremExceutablePath,
		TimeoutSeconds:    defaultDumbAnnealerTimeout,
		ConfigFilePath:    "testdata/DumbAnnealerTestConfig-ThreeRunsConcurrently.toml",
		ExpectedErrorCode: appTesting.WithSuccess,
	}

	appTesting.TestExecutableAgainstConfigFile(context)
}

func TestKirkpatrickDumbAnnealerIntegrationOneRun(t *testing.T) {
	context := appTesting.BinaryTestingContext{
		Name:              "Single run of Kirkpatrick Dumb Annealer",
		ExecutablePath:    cremExceutablePath,
		TimeoutSeconds:    defaultDumbAnnealerTimeout,
		T:                 t,
		ConfigFilePath:    "testdata/KirkpatrickDumbAnnealerTestConfig-OneRun.toml",
		ExpectedErrorCode: appTesting.WithSuccess,
	}

	appTesting.TestExecutableAgainstConfigFile(context)
}
