// +build windows
// Copyright (c) 2018 Australian Rivers Institute.

package main

import (
	. "github.com/onsi/gomega"
	"os"
	"testing"
)

const baseTestFilePath = "testdata/SimpleExcelAnnealerTestConfig-OneRun"
const configFileUnderTest = baseTestFilePath + ".toml"
const excelFileUnderTest = baseTestFilePath + ".xls"

type testContext struct {
	name       string
	t          *testing.T
	configFile string
}

func TestAnnealerIntegrationOneRun(t *testing.T) {
	context := testContext{
		name:       "Single run of Simple Excel Annealer",
		t:          t,
		configFile: configFileUnderTest,
	}

	verifyAnnealerRunsAgainstContext(context)
	os.Remove(excelFileUnderTest)
}

func verifyAnnealerRunsAgainstContext(context testContext) {
	if testing.Short() {
		context.t.Skip("skipping " + context.name + " in short mode")
	}
	g := NewGomegaWithT(context.t)

	simulatedMainCall := func() {
		RunFromConfigFile(context.configFile)
	}

	g.Expect(simulatedMainCall).To(Not(Panic()), context.name+" should not panic")
}
