package variable

import (
	"encoding/json"
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

	var derivedData map[string]interface{}
	unmnarshalError := json.Unmarshal(jsonOfVariableUnderTest, &derivedData)

	g.Expect(unmnarshalError).To(BeNil())
	g.Expect(derivedData["Name"]).To(Equal("SimpleEncodeableDecisionVariable"))
	g.Expect(derivedData["Measure"]).To(Equal(Dollars.String()))
	g.Expect(derivedData["Value"]).To(Equal("42.42"))
}

func TestPerPlanningUnitDecisionVariable_MarshalJson_AsExpected(t *testing.T) {
	g := NewGomegaWithT(t)

	variableUnderTest := EncodeableDecisionVariable{
		Name:    "PerPlanningUnitEncodeableDecisionVariable",
		Measure: TonnesPerYear,
		Value:   84.84,
		ValuePerPlanningUnit: PlanningUnitValues{
			PlanningUnitValue{
				PlanningUnit: 18,
				Value:        43.43,
			},
			PlanningUnitValue{
				PlanningUnit: 19,
				Value:        41.41,
			},
		},
	}

	jsonOfVariableUnderTest, marshalError := variableUnderTest.MarshalJSON()

	g.Expect(marshalError).To(BeNil())
	t.Log(string(jsonOfVariableUnderTest))

	var derivedData map[string]interface{}
	unmnarshalError := json.Unmarshal(jsonOfVariableUnderTest, &derivedData)

	g.Expect(unmnarshalError).To(BeNil())
	g.Expect(derivedData["Name"]).To(Equal("PerPlanningUnitEncodeableDecisionVariable"))
	g.Expect(derivedData["Measure"]).To(Equal(TonnesPerYear.String()))
	g.Expect(derivedData["Value"]).To(Equal("84.840"))

	rawMap := derivedData["ValuePerPlanningUnit"]
	arrayMap, isArray := rawMap.([]interface{})
	if !isArray {
		g.Expect("").ToNot(BeEmpty(), "ValuePerPlanningUnit map didn't match expected type")
	}

	entry0 := arrayMap[0]
	entry0Map, is0Map := entry0.(map[string]interface{})
	if !is0Map {
		g.Expect("").ToNot(BeEmpty(), "entry0Map map didn't match expected type")
	}
	g.Expect(entry0Map["PlanningUnit"]).To(Equal("18"))
	g.Expect(entry0Map["Value"]).To(Equal("43.430"))

	entry1 := arrayMap[1]
	entry1Map, is1Map := entry1.(map[string]interface{})
	if !is1Map {
		g.Expect("").ToNot(BeEmpty(), "entry1Map map didn't match expected type")
	}
	g.Expect(entry1Map["PlanningUnit"]).To(Equal("19"))
	g.Expect(entry1Map["Value"]).To(Equal("41.410"))
}
