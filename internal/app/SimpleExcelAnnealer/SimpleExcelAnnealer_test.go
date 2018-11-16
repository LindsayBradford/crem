// +build windows
// Copyright (c) 2018 Australian Rivers Institute.

package main

import (
	"os"
	"testing"

	configTesting "github.com/LindsayBradford/crem/internal/pkg/config/testing"
)

const baseTestFilePath = "testdata/SimpleExcelAnnealerTestConfig-OneRun"
const configFileUnderTest = baseTestFilePath + ".toml"
const excelFileUnderTest = baseTestFilePath + ".xls"

func TestAnnealerIntegrationOneRun(t *testing.T) {
	context := configTesting.TestingContext{
		Name:           "Single run of Simple Excel Annealer",
		T:              t,
		ConfigFilePath: configFileUnderTest,
		Runner:         RunFromConfigFile,
	}

	context.VerifyScenarioRunViaConfigFileDoesNotPanic()
	os.Remove(excelFileUnderTest)
}
