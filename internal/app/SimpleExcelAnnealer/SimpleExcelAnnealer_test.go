// +build windows
// Copyright (c) 2018 Australian Rivers Institute.

package main

import (
	"os"
	"testing"

	appTesting "github.com/LindsayBradford/crem/internal/pkg/config/testing"
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

func TestAnnealerIntegrationOneRun(t *testing.T) {
	context := appTesting.BlackboxTestingContext{
		Name:           "Single run of CremEngine CatchmentModel",
		ExecutablePath: exceutablePath,
		TimeoutSeconds: 20,
		T:              t,
		ConfigFilePath: configFileUnderTest,

		ExpectedErrorCode: appTesting.WithSuccess,
	}

	appTesting.TestExecutableAgainstConfigFile(context)
	os.Remove(excelFileUnderTest)
}
