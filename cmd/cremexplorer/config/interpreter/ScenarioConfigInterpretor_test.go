// Copyright (c) 2019 Australian Rivers Institute.

package interpreter

import (
	"testing"

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
