package variable

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestBasicEncodableDecisionVariable_MarshalJson_AsExpected(t *testing.T) {
	g := NewGomegaWithT(t)

	variableUnderTest := EncodeableDecisionVariable{
		Name:    "SimpleEncodeableDecisionVariable",
		Measure: Dollars,
		Value:   42.42,
	}

	jsonOfVariableUnderTest, marshalError := variableUnderTest.MarshalJSON()

	g.Expect(marshalError).To(BeNil())
	t.Log(string(jsonOfVariableUnderTest))
}

func TestPerPlanningUnitDecisionVariable_MarshalJson_AsExpected(t *testing.T) {
	g := NewGomegaWithT(t)

	variableUnderTest := EncodeableDecisionVariable{
		Name:    "PerPlanningUnitEncodeableDecisionVariable",
		Measure: Dollars,
		Value:   42.42,
		ValuePerPlanningUnit: PlanningUnitValues{
			PlanningUnitValue{
				PlanningUnit: 18,
				Value:        12.12,
			},
			PlanningUnitValue{
				PlanningUnit: 19,
				Value:        30.30,
			},
		},
	}

	jsonOfVariableUnderTest, marshalError := variableUnderTest.MarshalJSON()

	g.Expect(marshalError).To(BeNil())
	t.Log(string(jsonOfVariableUnderTest))
}
