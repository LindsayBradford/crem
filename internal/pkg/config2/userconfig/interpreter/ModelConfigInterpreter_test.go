// Copyright (c) 2019 Australian Rivers Institute.

package interpreter

import (
	"testing"

	"github.com/LindsayBradford/crem/internal/pkg/config2/userconfig/data"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/dumb"
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
	. "github.com/onsi/gomega"
)

const equalTo = "=="

func TestConfigInterpreter_NewModelConfigInterpreter_NullModelNoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	// when
	interpreterUnderTest := NewModelConfigInterpreter()

	// then
	g.Expect(interpreterUnderTest.Model()).To(BeAssignableToTypeOf(model.NullModel))
	g.Expect(interpreterUnderTest.Errors()).To(BeNil())
}

func TestConfigInterpreter_NullModel_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	// when
	configUnderTest := data.ModelConfig{
		Type: "NullModel",
	}
	interpreterUnderTest := NewModelConfigInterpreter().Interpret(&configUnderTest)

	// then
	g.Expect(interpreterUnderTest.Model()).To(BeAssignableToTypeOf(model.NullModel))
	g.Expect(interpreterUnderTest.Errors()).To(BeNil())
}

func TestConfigInterpreter_CatchmentModel_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	// when
	configUnderTest := data.ModelConfig{
		Type: "CatchmentModel",
	}
	interpreterUnderTest := NewModelConfigInterpreter().Interpret(&configUnderTest)

	// then
	g.Expect(interpreterUnderTest.Model()).To(BeAssignableToTypeOf(&catchment.Model{}))
	g.Expect(interpreterUnderTest.Errors()).To(BeNil())
}

func TestConfigInterpreter_ValidDumbModelWithParameter_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	// when
	expectedObjectiveValue := 4242.42
	parametersUnderTest := parameters.Map{
		"InitialObjectiveValue": expectedObjectiveValue,
	}
	configUnderTest := data.ModelConfig{
		Type:       "DumbModel",
		Parameters: parametersUnderTest,
	}
	interpreterUnderTest := NewModelConfigInterpreter().Interpret(&configUnderTest)

	// then
	actualModel := interpreterUnderTest.Model()
	g.Expect(actualModel).To(BeAssignableToTypeOf(dumb.NewModel()))
	g.Expect(actualModel.DecisionVariable("ObjectiveValue").Value()).To(BeNumerically(equalTo, expectedObjectiveValue))
	g.Expect(interpreterUnderTest.Errors()).To(BeNil())
}

func TestModelConfigInterpreter_RegisteringValidDummyModel_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	configUnderTest := data.ModelConfig{Type: "dummyModel"}

	dummyModelConfigFunctions := func(config data.ModelConfig) model.Model {
		newModel := new(dummyModel)
		newModel.Initialise()
		return newModel
	}

	// when
	interpreterUnderTest := NewModelConfigInterpreter().
		RegisteringModel("dummyModel", dummyModelConfigFunctions).
		Interpret(&configUnderTest)

	// then
	g.Expect(interpreterUnderTest.Model()).To(BeAssignableToTypeOf(&dummyModel{}))
	g.Expect(interpreterUnderTest.Errors()).To(BeNil())
}

type dummyModel struct {
	dumb.Model
}

func TestModelConfigInterpreter_NoSuchModel_Error(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	configUnderTest := data.ModelConfig{Type: "noSuchModel"}

	// when
	interpreterUnderTest := NewModelConfigInterpreter().Interpret(&configUnderTest)

	// then
	g.Expect(interpreterUnderTest.Model()).To(Equal(model.NullModel))
	g.Expect(interpreterUnderTest.Errors()).To(Not(BeNil()))
	t.Log(interpreterUnderTest.Errors())
}

func TestModelConfigInterpreter_BadModelParameter_Error(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	parametersUnderTest := parameters.Map{
		"NotAParameterForTheModel": "really... it doesn't exist",
	}
	configUnderTest := data.ModelConfig{Type: "MultiObjectiveDumbModel", Parameters: parametersUnderTest}

	// when
	interpreterUnderTest := NewModelConfigInterpreter().Interpret(&configUnderTest)

	// then
	g.Expect(interpreterUnderTest.Model()).To(Equal(model.NullModel))
	g.Expect(interpreterUnderTest.Errors()).To(Not(BeNil()))
	t.Log(interpreterUnderTest.Errors())
}
