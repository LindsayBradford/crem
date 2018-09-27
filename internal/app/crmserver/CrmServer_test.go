// Copyright (c) 2018 Australian Rivers Institute.

package main

import (
	"testing"

	"github.com/LindsayBradford/crm/internal/app/crmserver/components"
	. "github.com/onsi/gomega"
)

type testContext struct {
	name       string
	t          *testing.T
	configFile string
}

func TestDumbAnnealerIntegrationOneRun(t *testing.T) {
	context := testContext{
		name:       "Single run of Dumb annealer",
		t:          t,
		configFile: "testdata/DumbAnnealerTestConfig-OneRun.toml",
	}

	verifyDumbAnnealerRunsAgainstContext(context)
}

func TestDumbAnnealerIntegrationThreeRunsSequentially(t *testing.T) {
	context := testContext{
		name:       "Three sequential runs of Dumb annealer",
		t:          t,
		configFile: "testdata/DumbAnnealerTestConfig-ThreeRunsSequentially.toml",
	}

	verifyDumbAnnealerRunsAgainstContext(context)
}

func TestDumbAnnealerIntegrationThreeRunsConcurrently(t *testing.T) {
	context := testContext{
		name:       "Three concurrent runs of Dumb annealer",
		t:          t,
		configFile: "testdata/DumbAnnealerTestConfig-ThreeRunsConcurrently.toml",
	}

	verifyDumbAnnealerRunsAgainstContext(context)
}

func verifyDumbAnnealerRunsAgainstContext(context testContext) {
	if testing.Short() {
		context.t.Skip("skipping " + context.name + " in short mode")
	}
	g := NewGomegaWithT(context.t)

	simulatedMainCall := func() {
		components.RunScenarioFromConfigFile(context.configFile)
	}

	g.Expect(simulatedMainCall).To(Not(Panic()), context.name+" should not panic")
}
