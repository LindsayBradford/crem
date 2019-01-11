// Copyright (c) 2019 Australian Rivers Institute.

package main

import (
	"testing"

	"github.com/LindsayBradford/crem/cmd/cremengine/components"
	"github.com/LindsayBradford/crem/cmd/cremengine/components/scenario"
	testing2 "github.com/LindsayBradford/crem/internal/pkg/config/testing"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
)

func TestDumbAnnealerIntegrationThreeRunsSequentially(t *testing.T) {
	context := testing2.Context{
		Name:           "Three sequential runs of Dumb annealer",
		T:              t,
		ConfigFilePath: "testdata/DumbAnnealerTestConfig-ThreeRunsSequentially.toml",
		Runner:         components.RunScenarioFromConfigFile,
	}

	scenario.LogHandler = loggers.DefaultTestingLogger
	context.VerifyScenarioRunViaConfigFileDoesNotPanic()
}

func TestDumbAnnealerIntegrationOneRun(t *testing.T) {
	context := testing2.Context{
		Name:           "Single run of Dumb annealer",
		T:              t,
		ConfigFilePath: "testdata/DumbAnnealerTestConfig-OneRun.toml",
		Runner:         components.RunScenarioFromConfigFile,
	}

	scenario.LogHandler = loggers.DefaultTestingLogger
	context.VerifyScenarioRunViaConfigFileDoesNotPanic()
}

func TestDumbAnnealerIntegrationThreeRunsConcurrently(t *testing.T) {
	context := testing2.Context{
		Name:           "Three concurrent runs of Dumb annealer",
		T:              t,
		ConfigFilePath: "testdata/DumbAnnealerTestConfig-ThreeRunsConcurrently.toml",
		Runner:         components.RunScenarioFromConfigFile,
	}

	scenario.LogHandler = loggers.DefaultTestingLogger
	context.VerifyScenarioRunViaConfigFileDoesNotPanic()
}

func TestKirkpatrickDumbAnnealerIntegrationOneRun(t *testing.T) {
	context := testing2.Context{
		Name:           "Single run of Kirkpatrick Dumb annealer",
		T:              t,
		ConfigFilePath: "testdata/KirkpatrickDumbAnnealerTestConfig-OneRun.toml",
		Runner:         components.RunScenarioFromConfigFile,
	}

	context.VerifyScenarioRunViaConfigFileDoesNotPanic()
}
