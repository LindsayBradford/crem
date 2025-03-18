//go:build windows

// Copyright (c) 2019 Australian Rivers Institute.

package main

import (
	"testing"

	"github.com/LindsayBradford/crem/cmd/cremexplorer/bootstrap"
	configTesting "github.com/LindsayBradford/crem/internal/pkg/config/testing"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
)

const defaultDumbAnnealerTimeout = 10

func TestDumbAnnealerIntegrationOneRun(t *testing.T) {
	context := configTesting.BlackboxTestingContext{
		T:                 t,
		Name:              "Single run of Dumb Annealer",
		ExecutablePath:    dumbAnnealerExecutablePath,
		TimeoutSeconds:    defaultDumbAnnealerTimeout,
		ConfigFilePath:    "testdata/DumbAnnealerTestConfig-OneRun.toml",
		ExpectedErrorCode: configTesting.WithSuccess,
	}

	configTesting.TestExecutableAgainstConfigFile(context)
}

func TestDumbAnnealerIntegrationThreeRunsSequentially(t *testing.T) {
	context := configTesting.BlackboxTestingContext{
		T:                 t,
		Name:              "Three sequential runs of Dumb Annealer",
		ExecutablePath:    dumbAnnealerExecutablePath,
		TimeoutSeconds:    15,
		ConfigFilePath:    "testdata/DumbAnnealerTestConfig-ThreeRunsSequentially.toml",
		ExpectedErrorCode: configTesting.WithSuccess,
	}

	configTesting.TestExecutableAgainstConfigFile(context)
}

func TestDumbAnnealerIntegrationThreeRunsConcurrently(t *testing.T) {
	context := configTesting.BlackboxTestingContext{
		T:                 t,
		Name:              "Three concurrent runs of Dumb Annealer",
		ExecutablePath:    dumbAnnealerExecutablePath,
		TimeoutSeconds:    defaultDumbAnnealerTimeout,
		ConfigFilePath:    "testdata/DumbAnnealerTestConfig-ThreeRunsConcurrently.toml",
		ExpectedErrorCode: configTesting.WithSuccess,
	}

	configTesting.TestExecutableAgainstConfigFile(context)
}

func TestKirkpatrickDumbAnnealerIntegrationOneRun(t *testing.T) {
	context := configTesting.BlackboxTestingContext{
		Name:              "Single run of Kirkpatrick Dumb Annealer",
		ExecutablePath:    dumbAnnealerExecutablePath,
		TimeoutSeconds:    defaultDumbAnnealerTimeout,
		T:                 t,
		ConfigFilePath:    "testdata/KirkpatrickDumbAnnealerTestConfig-OneRun.toml",
		ExpectedErrorCode: configTesting.WithSuccess,
	}

	configTesting.TestExecutableAgainstConfigFile(context)
}

func TestKirkpatrickDumbAnnealer_ScenarioOneRunWhitebox_ExitWithSuccess(t *testing.T) {
	context := configTesting.WhiteboxTestingContext{
		Name:           "Single run of catchment model annealer",
		T:              t,
		ConfigFilePath: "testdata/KirkpatrickDumbAnnealerTestConfig-Whitebox-OneRun.toml",
		// ConfigFilePath:    "testdata/KirkpatrickDumbAnnealerTestConfig-BadConfig.toml",
		Runner: RunFromConfigFile,
	}

	bootstrap.LogHandler = loggers.DefaultTestingLogger
	context.VerifyGoroutineScenarioRunViaConfigFileDoesNotPanic()
}

func TestKirkpatrickDumbAnnealerBadConfig_Fails(t *testing.T) {
	context := configTesting.BlackboxTestingContext{
		T:                 t,
		Name:              "Kirkpatrick Dumb Annealer with Bad Config",
		ExecutablePath:    dumbAnnealerExecutablePath,
		TimeoutSeconds:    defaultDumbAnnealerTimeout,
		ConfigFilePath:    "testdata/KirkpatrickDumbAnnealerTestConfig-BadConfig.toml",
		ExpectedErrorCode: configTesting.WithFailure,
	}

	configTesting.TestExecutableAgainstConfigFile(context)
}
