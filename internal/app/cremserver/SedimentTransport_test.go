// Copyright (c) 2018 Australian Rivers Institute.

package main

import (
	"testing"

	configTesting "github.com/LindsayBradford/crem/config/testing"
	"github.com/LindsayBradford/crem/internal/app/cremserver/components"
)

func TestSedimentTransportAnnealerScenarioOneRun(t *testing.T) {
	context := configTesting.TestingContext{
		Name:           "Single run of sediment transport annealer",
		T:              t,
		ConfigFilePath: "testdata/SedimentTransportTestConfig-OneRun.toml",
		Runner:         components.RunScenarioFromConfigFile,
	}

	context.VerifyScenarioConfigFilesDoesNotPanic()
}
