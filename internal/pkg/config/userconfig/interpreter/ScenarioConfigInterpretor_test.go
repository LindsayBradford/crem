package interpreter

import (
	"github.com/LindsayBradford/crem/internal/pkg/scenario"
	"testing"

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
