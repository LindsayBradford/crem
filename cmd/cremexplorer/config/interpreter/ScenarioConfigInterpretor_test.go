// Copyright (c) 2019 Australian Rivers Institute.

package interpreter

import (
	"testing"

	"github.com/LindsayBradford/crem/cmd/cremexplorer/config/data"
	"github.com/LindsayBradford/crem/internal/pkg/scenario"

	. "github.com/onsi/gomega"
)

func TestScenarioConfigInterpreter_NewMScenarioConfigInterpreter_NullScenarioNoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	// when
	interpreterUnderTest := NewScenarioConfigInterpreter()

	// then
	g.Expect(interpreterUnderTest.Scenario()).To(BeAssignableToTypeOf(scenario.NullScenario))
	g.Expect(interpreterUnderTest.Errors()).To(BeNil())
}

func TestConfigInterpreter_UnnamedScenario_Errors(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	configUnderTest := data.ScenarioConfig{Name: ""}

	// when
	interpreterUnderTest := NewScenarioConfigInterpreter().Interpret(&configUnderTest)

	// then
	if interpreterUnderTest.Errors() != nil {
		t.Log(interpreterUnderTest.Errors())
	}
	g.Expect(interpreterUnderTest.Errors()).To(Not(BeNil()))
}

func TestConfigInterpreter_ProfilingScenario_HasProfilingRunner(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	configUnderTest := data.ScenarioConfig{
		Name:           "Profiling Scenario test",
		CpuProfilePath: "testdata",
	}

	// when
	interpreterUnderTest := NewScenarioConfigInterpreter().Interpret(&configUnderTest)

	// then
	if interpreterUnderTest.Errors() != nil {
		t.Log(interpreterUnderTest.Errors())
	}
	g.Expect(interpreterUnderTest.Errors()).To(BeNil())
}
