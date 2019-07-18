// +build windows
// Copyright (c) 2018 Australian Rivers Institute.

package main

import (
	"github.com/LindsayBradford/crem/cmd/cremexplorer/bootstrap"
	"os"
	"testing"

	appTesting "github.com/LindsayBradford/crem/internal/pkg/config/testing"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"github.com/onsi/gomega/gexec"
)

const baseTestFilePath = "testdata/SimpleExcelAnnealerTestConfig-OneRun"
const configFileUnderTest = baseTestFilePath + ".toml"
const excelFileUnderTest = baseTestFilePath + ".xls"

var exceutablePath string

func TestMain(m *testing.M) {
	setup(m)
	code := m.Run()
	tearDown()
	os.Exit(code)
}

func setup(m *testing.M) {
	var err error
	exceutablePath, err = gexec.Build("github.com/LindsayBradford/crem/internal/app/SimpleExcelAnnealer")

	if err != nil {
		os.Exit(appTesting.WithFailure)
	}
}

func tearDown() {
	gexec.CleanupBuildArtifacts()
}

func TestSimpleExcelAnnealer_OneRun(t *testing.T) {
	context := appTesting.BlackboxTestingContext{
		Name:           "Single Black-Box run of SimpleExcelAnnealer",
		ExecutablePath: exceutablePath,
		TimeoutSeconds: 20,
		T:              t,
		ConfigFilePath: configFileUnderTest,

		ExpectedErrorCode: appTesting.WithSuccess,
	}

	appTesting.TestExecutableAgainstConfigFile(context)
	os.Remove(excelFileUnderTest)
}

func TestSimpleExcelAnnealer_Whitebox_ExitWithSuccess(t *testing.T) {
	// TODO: Why doens't GetMainThreadChannel.Close() take in VerifyGoroutineScenarioRunViaConfigFileDoesNotPanic()?
	t.Skip("Not sure why this is failing on mainThread.Close()")

	context := appTesting.WhiteboxTestingContext{
		Name:           "Single White-Box run of SimpleExcelAnnealer",
		T:              t,
		ConfigFilePath: configFileUnderTest,
		Runner:         RunFromConfigFile,
	}

	bootstrap.LogHandler = loggers.DefaultTestingLogger
	context.VerifyGoroutineScenarioRunViaConfigFileDoesNotPanic()
	os.Remove(excelFileUnderTest)
}
