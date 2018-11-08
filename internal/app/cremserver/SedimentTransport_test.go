// Copyright (c) 2018 Australian Rivers Institute.

package main

import (
	"testing"

	"github.com/LindsayBradford/crem/internal/app/cremserver/components"
	. "github.com/onsi/gomega"
)

func TestSedimentTransportAnnealerScenarioOneRun(t *testing.T) {
	context := testContext{
		name:       "Single run of sediment transport annealer",
		t:          t,
		configFile: "testdata/SedimentTransportTestConfig-OneRun.toml",
	}

	verifyScenarioRunsAgainstContext(context)
}

func verifyScenarioRunsAgainstContext(context testContext) {
	if testing.Short() {
		context.t.Skip("skipping " + context.name + " in short mode")
	}
	g := NewGomegaWithT(context.t)

	simulatedMainCall := func() {
		components.RunScenarioFromConfigFile(context.configFile)
	}

	g.Expect(simulatedMainCall).To(Not(Panic()), context.name+" should not panic")
}
