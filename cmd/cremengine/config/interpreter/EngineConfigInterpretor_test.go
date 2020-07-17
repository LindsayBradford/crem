// Copyright (c) 2019 Australian Rivers Institute.

package interpreter

import (
	data "github.com/LindsayBradford/crem/cmd/cremengine/config/data"
	"github.com/LindsayBradford/crem/cmd/cremengine/engine"
	"testing"

	. "github.com/onsi/gomega"
)

func TestScenarioConfigInterpreter_NewMScenarioConfigInterpreter_NullScenarioNoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	// when
	interpreterUnderTest := NewEngineConfigInterpreter()

	// then
	g.Expect(interpreterUnderTest.Engine()).To(BeAssignableToTypeOf(engine.NullEngine))
	g.Expect(interpreterUnderTest.Errors()).To(BeNil())
}

func TestConfigInterpreter_EmptyEngine_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	configUnderTest := data.EngineConfig{}

	// when
	interpreterUnderTest := NewEngineConfigInterpreter().Interpret(configUnderTest.Engine)

	// then
	if interpreterUnderTest.Errors() != nil {
		t.Log(interpreterUnderTest.Errors())
	}
	g.Expect(interpreterUnderTest.Errors()).To(BeNil())
}
