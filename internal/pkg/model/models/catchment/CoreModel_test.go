// Copyright (c) 2019 Australian Rivers Institute.

package catchment

import (
	model2 "github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/model/archive"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables/dissolvednitrogen"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables/opportunitycost"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables/particulatenitrogen"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables/totalnitrogen"
	"testing"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/csv"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/actions"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables/implementationcost"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables/sedimentproduction"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"github.com/LindsayBradford/crem/pkg/math"
	. "github.com/onsi/gomega"
)

const expectedName = "CatchmentModel"
const expectedMaximumImplementationCost = 65_000.0
const expectedMaximumSedimentProduction = 65_000.0

const equalTo = "=="

func TestCoreModel_NewCoreModel(t *testing.T) {
	g := NewGomegaWithT(t)

	model := NewCoreModel()
	actualName := model.Name()

	g.Expect(actualName).To(Equal(expectedName))

	actualActions := model.ManagementActions()
	expectedActionNumber := 0

	g.Expect(len(actualActions)).To(BeNumerically(equalTo, expectedActionNumber))

	actualVariables := model.DecisionVariables()
	expectedVariableNumber := 0

	g.Expect(len(*actualVariables)).To(BeNumerically(equalTo, expectedVariableNumber))
}

func TestCoreModel_Initialise_ValidDataSet_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	localExpectedName := "InitialiseTest"

	sourceDataSet := csv.NewDataSet("CatchmentModel")
	loadError := sourceDataSet.Load("testdata/ValidModel.csv")

	g.Expect(loadError).To(BeNil())

	model := NewCoreModel().
		WithSourceDataSet(sourceDataSet).
		WithName(localExpectedName)

	g.Expect(model.Name()).To(Equal(localExpectedName))

	model.Initialise(model2.Random)

	actualActions := model.ManagementActions()
	expectedActionNumber := 13

	g.Expect(len(actualActions)).To(BeNumerically(equalTo, expectedActionNumber))

	actualVariables := *model.DecisionVariables()

	g.Expect(actualVariables).To(HaveKey(implementationcost.VariableName))
	g.Expect(actualVariables[implementationcost.VariableName].Value()).To(BeNumerically(equalTo, 0))

	g.Expect(actualVariables).To(HaveKey(sedimentproduction.VariableName))
}

func TestCoreModel_InitialiseAndClone_ValidDataSet_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	localExpectedName := "InitialiseTest"

	sourceDataSet := csv.NewDataSet("CatchmentModel")
	loadError := sourceDataSet.Load("testdata/ValidModel.csv")

	g.Expect(loadError).To(BeNil())

	model := NewCoreModel().
		WithSourceDataSet(sourceDataSet).
		WithName(localExpectedName)

	g.Expect(model.Name()).To(Equal(localExpectedName))

	model.Initialise(model2.Random)
	model.SetManagementAction(0, true)
	model.AcceptAll()

	originalActions := *model.DecisionVariables()

	g.Expect(originalActions).To(HaveKey(implementationcost.VariableName))
	g.Expect(originalActions[implementationcost.VariableName].Value()).To(BeNumerically(">", 0))

	copiedModel := model.DeepClone()
	copiedModel.Initialise(model2.Unchanged)

	actualActions := copiedModel.ManagementActions()
	expectedActionNumber := 13

	g.Expect(len(actualActions)).To(BeNumerically(equalTo, expectedActionNumber))

	actualVariables := *copiedModel.DecisionVariables()

	g.Expect(actualVariables).To(HaveKey(implementationcost.VariableName))
	g.Expect(actualVariables[implementationcost.VariableName].Value()).To(BeNumerically(equalTo, 0))

	g.Expect(actualVariables).To(HaveKey(sedimentproduction.VariableName))
}

func TestCoreModel_Initialise_InvalidDataSet_Errors(t *testing.T) {
	g := NewGomegaWithT(t)

	sourceDataSet := csv.NewDataSet("CatchmentModel")
	loadError := sourceDataSet.Load("testdata/InvalidModel.csv")

	g.Expect(loadError).To(BeNil())

	newModelRunner := func() {
		NewCoreModel().WithSourceDataSet(sourceDataSet).Initialise(model2.AsIs)
	}

	g.Expect(newModelRunner).To(Panic())
}

func TestCoreModel_WithDefaultParameters_NoErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	sourceDataSet := csv.NewDataSet("CatchmentModel")
	loadError := sourceDataSet.Load("testdata/ValidModel.csv")

	g.Expect(loadError).To(BeNil())

	parametersUnderTest := parameters.Map{}

	modelUnderTest := NewCoreModel().
		WithSourceDataSet(sourceDataSet).
		WithParameters(parametersUnderTest)

	parameterErrors := modelUnderTest.ParameterErrors()

	g.Expect(parameterErrors).To(BeNil())
}

func TestCoreModel_PlanningUnitValues_AsExpected(t *testing.T) {
	g := NewGomegaWithT(t)

	modelUnderTest := buildTestingModel(g)

	solution := new(solution.SolutionBuilder).
		WithId("testingBuilder").
		ForModel(modelUnderTest).
		Build()

	g.Expect(solution).To(Not(BeNil()))

	verifyPlanningUnitValues(g, solution, implementationcost.VariableName, 0)
	verifyPlanningUnitValues(g, solution, opportunitycost.VariableName, 0)

	verifyPlanningUnitValues(g, solution, sedimentproduction.VariableName, 1322.548)
	verifyPlanningUnitValues(g, solution, particulatenitrogen.VariableName, 2.754)
}

func TestCoreModel_AfterActionToggling_PlanningUnitValues_AsExpected(t *testing.T) {
	// given
	g := NewGomegaWithT(t)

	modelUnderTest := buildTestingModel(g)
	builder := new(solution.SolutionBuilder).
		WithId("testingBuilder").
		ForModel(modelUnderTest)

	solution := builder.Build()

	g.Expect(solution).To(Not(BeNil()))

	// when

	planningUnit := planningunit.Id(18)

	modelUnderTest.ToggleAction(planningUnit, actions.RiverBankRestorationType)
	modelUnderTest.AcceptChange()
	modelUnderTest.ToggleAction(planningUnit, actions.HillSlopeRestorationType)
	modelUnderTest.AcceptChange()
	modelUnderTest.ToggleAction(planningUnit, actions.RiverBankRestorationType)
	modelUnderTest.AcceptChange()
	modelUnderTest.ToggleAction(planningUnit, actions.HillSlopeRestorationType)
	modelUnderTest.AcceptChange()

	// then
	newSolution := builder.Build()

	g.Expect(newSolution).To(Not(BeNil()))

	verifyPlanningUnitValues(g, newSolution, implementationcost.VariableName, 0)
	verifyPlanningUnitValues(g, newSolution, opportunitycost.VariableName, 0)

	verifyPlanningUnitValues(g, newSolution, sedimentproduction.VariableName, 1322.548)
	verifyPlanningUnitValues(g, newSolution, particulatenitrogen.VariableName, 2.754)
}

func TestCoreModel_IsEquivalentTo_AsExpected(t *testing.T) {
	// given
	g := NewGomegaWithT(t)
	const planningUnitUnderTest = 18
	planningUnit := planningunit.Id(planningUnitUnderTest)

	// when

	firstModelUnderTest := buildTestingModel(g)
	secondModelUnderTest := buildTestingModel(g)

	// then

	g.Expect(firstModelUnderTest.IsEquivalentTo(secondModelUnderTest)).To(BeTrue())

	// when

	firstModelUnderTest.ToggleAction(planningUnit, actions.RiverBankRestorationType)
	firstModelUnderTest.AcceptChange()

	secondModelUnderTest.ToggleAction(planningUnit, actions.RiverBankRestorationType)
	secondModelUnderTest.AcceptChange()

	// then

	g.Expect(firstModelUnderTest.IsEquivalentTo(secondModelUnderTest)).To(BeTrue())

	// when

	secondModelUnderTest.ToggleAction(planningUnit, actions.RiverBankRestorationType)
	secondModelUnderTest.AcceptChange()

	// then

	g.Expect(firstModelUnderTest.IsEquivalentTo(secondModelUnderTest)).To(BeFalse())
}

func TestCoreModel_Compression_AsExpected(t *testing.T) {
	// given
	g := NewGomegaWithT(t)
	const planningUnitUnderTest = 18
	modelArchive := new(archive.NonDominanceModelArchive).Initialise()

	// when

	modelUnderTest := buildTestingModel(g)
	compressedModelUnderTest := modelArchive.Compress(modelUnderTest)

	decompressedModel := modelUnderTest.DeepClone()
	decompressedModel.Initialise(model2.Unchanged)
	modelArchive.Decompress(compressedModelUnderTest, decompressedModel)

	// then

	g.Expect(modelUnderTest.IsEquivalentTo(decompressedModel)).To(BeTrue())
}

func TestCoreModel_CompressionOfChanged_AsExpected(t *testing.T) {
	// given
	g := NewGomegaWithT(t)
	const planningUnitUnderTest = 18
	planningUnit := planningunit.Id(planningUnitUnderTest)
	modelArchive := new(archive.NonDominanceModelArchive).Initialise()

	// when

	modelUnderTest := buildTestingModel(g)
	modelUnderTest.ToggleAction(planningUnit, actions.RiverBankRestorationType)
	modelUnderTest.AcceptChange()

	compressedModelUnderTest := modelArchive.Compress(modelUnderTest)

	decompressedModel := modelUnderTest.DeepClone()
	decompressedModel.Initialise(model2.AsIs)
	modelArchive.Decompress(compressedModelUnderTest, decompressedModel)

	// then

	g.Expect(modelUnderTest.IsEquivalentTo(decompressedModel)).To(BeTrue())
}

func TestCoreModel_ParticulateNitrogen_NoRoundingErrors(t *testing.T) {
	// given
	g := NewGomegaWithT(t)
	const planningUnitUnderTest = 18

	modelUnderTest := buildTestingModel(g)
	builder := new(solution.SolutionBuilder).
		WithId("testingBuilder").
		ForModel(modelUnderTest)

	solution := builder.Build()

	g.Expect(solution).To(Not(BeNil()))

	// when

	planningUnit := planningunit.Id(planningUnitUnderTest)

	for index := 0; index < 1_000; index++ {
		modelUnderTest.ToggleAction(planningUnit, actions.RiverBankRestorationType)
		modelUnderTest.AcceptChange()

		modelUnderTest.ToggleAction(planningUnit, actions.HillSlopeRestorationType)
		modelUnderTest.AcceptChange()

		modelUnderTest.ToggleAction(planningUnit, actions.RiverBankRestorationType)
		modelUnderTest.AcceptChange()

		modelUnderTest.ToggleAction(planningUnit, actions.HillSlopeRestorationType)
		modelUnderTest.AcceptChange()
	}
	// then
	newSolution := builder.Build()

	g.Expect(newSolution).To(Not(BeNil()))

	variableUnderTest := solutionVariable(solution, particulatenitrogen.VariableName)
	planningUnit18Entry := variableUnderTest.ValuePerPlanningUnit[1]

	g.Expect(variableUnderTest.Value).To(BeNumerically(equalTo, 2.754))
	g.Expect(planningUnit18Entry.PlanningUnit).To(BeNumerically(equalTo, planningUnitUnderTest))
	g.Expect(planningUnit18Entry.Value).To(BeNumerically(equalTo, 2.291))
}

func TestCoreModel_ParticulateNitrogen_HillSlopeRiverbankDependency_NoRoundingErrors(t *testing.T) {
	// given
	g := NewGomegaWithT(t)
	const planningUnitUnderTest = 112

	modelUnderTest := buildTestingModel(g)
	builder := new(solution.SolutionBuilder).
		WithId("testingBuilder").
		ForModel(modelUnderTest)

	solution := builder.Build()

	g.Expect(solution).To(Not(BeNil()))

	// when

	planningUnit := planningunit.Id(planningUnitUnderTest)

	for index := 0; index < 1_000; index++ {
		modelUnderTest.ToggleAction(planningUnit, actions.RiverBankRestorationType)
		modelUnderTest.AcceptChange()

		modelUnderTest.ToggleAction(planningUnit, actions.HillSlopeRestorationType)
		modelUnderTest.AcceptChange()

		modelUnderTest.ToggleAction(planningUnit, actions.RiverBankRestorationType)
		modelUnderTest.AcceptChange()

		modelUnderTest.ToggleAction(planningUnit, actions.HillSlopeRestorationType)
		modelUnderTest.AcceptChange()
	}

	// then
	newSolution := builder.Build()

	g.Expect(newSolution).To(Not(BeNil()))

	variableUnderTest := solutionVariable(solution, particulatenitrogen.VariableName)
	planningUnitEntry := variableUnderTest.ValuePerPlanningUnit[3]

	g.Expect(planningUnitEntry.PlanningUnit).To(BeNumerically(equalTo, planningUnitUnderTest))
	g.Expect(planningUnitEntry.Value).To(BeNumerically(equalTo, 0.404))

	g.Expect(variableUnderTest.Value).To(BeNumerically(equalTo, 2.754))
}

func TestCoreModel_TotalNitrogen_NoRoundingErrors(t *testing.T) {
	// given
	g := NewGomegaWithT(t)
	const planningUnitUnderTest = 18
	const wetlandsPlanningUnitUnderTest = 21

	modelUnderTest := buildTestingModel(g)
	builder := new(solution.SolutionBuilder).
		WithId("testingBuilder").
		ForModel(modelUnderTest)

	solution := builder.Build()

	g.Expect(solution).To(Not(BeNil()))

	// when

	planningUnit := planningunit.Id(planningUnitUnderTest)
	wetlandsPlanningUnit := planningunit.Id(wetlandsPlanningUnitUnderTest)

	for index := 0; index < 1_000; index++ {
		modelUnderTest.ToggleAction(planningUnit, actions.RiverBankRestorationType)
		modelUnderTest.AcceptChange()

		modelUnderTest.ToggleAction(planningUnit, actions.HillSlopeRestorationType)
		modelUnderTest.AcceptChange()

		modelUnderTest.ToggleAction(planningUnit, actions.RiverBankRestorationType)
		modelUnderTest.AcceptChange()

		modelUnderTest.ToggleAction(planningUnit, actions.HillSlopeRestorationType)
		modelUnderTest.AcceptChange()

		modelUnderTest.ToggleAction(wetlandsPlanningUnit, actions.WetlandsEstablishmentType)
		modelUnderTest.AcceptChange()

		modelUnderTest.ToggleAction(wetlandsPlanningUnit, actions.WetlandsEstablishmentType)
		modelUnderTest.AcceptChange()
	}
	// then
	newSolution := builder.Build()

	g.Expect(newSolution).To(Not(BeNil()))

	variableUnderTest := solutionVariable(solution, totalnitrogen.VariableName)
	baseParticulateVariable := solutionVariable(solution, particulatenitrogen.VariableName)
	baseDissolvedVariable := solutionVariable(solution, dissolvednitrogen.VariableName)

	g.Expect(variableUnderTest.Value).To(BeNumerically("==", baseParticulateVariable.Value+baseDissolvedVariable.Value))

	for _, planningUnitValue := range variableUnderTest.ValuePerPlanningUnit {
		actualValue := planningUnitValue

		expectedParticulateValue := float64(0)
		for _, particulateValue := range baseParticulateVariable.ValuePerPlanningUnit {
			if particulateValue.PlanningUnit == actualValue.PlanningUnit {
				expectedParticulateValue = particulateValue.Value
			}
		}

		expectedDissolvedValue := float64(0)
		for _, dissolvedValue := range baseDissolvedVariable.ValuePerPlanningUnit {
			if dissolvedValue.PlanningUnit == actualValue.PlanningUnit {
				expectedDissolvedValue = dissolvedValue.Value
			}
		}

		expectedValue := math.RoundFloat(expectedParticulateValue+expectedDissolvedValue, 3)
		g.Expect(actualValue.Value).To(BeNumerically("==", expectedValue))
	}
}

func verifyPlanningUnitValues(g *GomegaWithT, solution *solution.Solution, variableName string, expectedValue float64) {
	variableUnderTest := solutionVariable(solution, variableName)
	g.Expect(variableUnderTest.Value).To(BeNumerically(equalTo, expectedValue))

	var planningUnitValues float64
	for _, currValue := range variableUnderTest.ValuePerPlanningUnit {
		planningUnitValues += currValue.Value
	}
	precisionOfVariable := math.DerivePrecision(variableUnderTest.Value)
	roundedPlanningUnitValues := math.RoundFloat(planningUnitValues, precisionOfVariable)

	g.Expect(variableUnderTest.Value).To(BeNumerically(equalTo, roundedPlanningUnitValues))
}

func TestCoreModel_ToggleRiverBankRestoration_AsExpected(t *testing.T) {
	g := NewGomegaWithT(t)

	modelUnderTest := buildTestingModel(g)
	planningUnit := planningunit.Id(18)

	verifyActionToggle(t, modelUnderTest, planningUnit, actions.RiverBankRestorationType, g)
}

func TestCoreModel_ToggleGullyRestoration_AsExpected(t *testing.T) {
	g := NewGomegaWithT(t)

	modelUnderTest := buildTestingModel(g)
	planningUnit := planningunit.Id(18)

	verifyActionToggle(t, modelUnderTest, planningUnit, actions.GullyRestorationType, g)
}

func TestCoreModel_ToggleHillSlopeRestoration_AsExpected(t *testing.T) {
	g := NewGomegaWithT(t)

	modelUnderTest := buildTestingModel(g)
	planningUnit := planningunit.Id(18)

	verifyActionToggle(t, modelUnderTest, planningUnit, actions.HillSlopeRestorationType, g)
}

func TestCoreModel_Bounded_InitialisationValid(t *testing.T) {
	g := NewGomegaWithT(t)

	modelUnderTest := buildBoundedTestingModel(g)

	changeState, changeErrors := modelUnderTest.ChangeIsValid()

	g.Expect(changeState).To(BeTrue())
	if changeErrors != nil {
		t.Log(changeErrors)
	}
	g.Expect(changeErrors).To(BeNil())
}

func TestCoreModel_Bounded_RandomisationStaysValid(t *testing.T) {
	g := NewGomegaWithT(t)

	modelUnderTest := buildBoundedTestingModel(g)

	modelUnderTest.InitialiseActions(model2.Random)

	changeState, changeErrors := modelUnderTest.StateIsValid()

	if changeErrors != nil {
		t.Log(changeErrors)
	}
	g.Expect(changeState).To(BeTrue())
	g.Expect(changeErrors).To(BeNil())

	implementationCost := modelUnderTest.DecisionVariable("ImplementationCost")
	t.Logf("Against bound of %f, ImplementationCost = %f", expectedMaximumImplementationCost, implementationCost.Value())

	g.Expect(implementationCost.Value()).To(BeNumerically("<", expectedMaximumImplementationCost))
}

func TestCoreModel_Bounded_ValidityAsExpected(t *testing.T) {
	g := NewGomegaWithT(t)

	modelUnderTest := buildBoundedTestingModel(g)
	changeState, changeErrors := modelUnderTest.ChangeIsValid()

	if changeErrors != nil {
		t.Log(changeErrors)
	}
	g.Expect(changeState).To(BeTrue())
	g.Expect(changeErrors).To(BeNil())

	modelUnderTest.ToggleAction(17, actions.GullyRestorationType)

	changeState, changeErrors = modelUnderTest.ChangeIsValid()

	if changeErrors != nil {
		t.Log(changeErrors)
	}
	g.Expect(changeState).To(BeTrue())
	g.Expect(changeErrors).To(BeNil())

	modelUnderTest.AcceptChange()

	modelSnapshot := new(solution.SolutionBuilder).
		WithId("modelUnderTest").
		ForModel(modelUnderTest).
		Build()

	state, stateErrors := modelUnderTest.StateIsValid()

	if stateErrors != nil {
		t.Log(stateErrors)
	}
	g.Expect(state).To(BeTrue())
	g.Expect(stateErrors).To(BeNil())

	modelUnderTest.ToggleAction(18, actions.GullyRestorationType)

	changeState, changeErrors = modelUnderTest.ChangeIsValid()

	if changeErrors != nil {
		t.Log(changeErrors)
	}
	g.Expect(changeState).To(BeFalse())
	g.Expect(changeErrors).To(Not(BeNil()))

	modelUnderTest.RevertChange()

	state, stateErrors = modelUnderTest.StateIsValid()

	if stateErrors != nil {
		t.Log(stateErrors)
	}
	g.Expect(state).To(BeTrue())
	g.Expect(stateErrors).To(BeNil())

	currentSnapshot := new(solution.SolutionBuilder).
		WithId("modelUnderTest").
		ForModel(modelUnderTest).
		Build()

	verifySolutionsMatch(t, g, modelSnapshot, currentSnapshot)
}

func TestCoreModel_MoreThanOneBoundedParameter_ParameterErrors(t *testing.T) {
	g := NewGomegaWithT(t)

	errors := buildInvalidBoundedTestingModel(g)

	if errors != nil {
		t.Log(errors)
	}
	g.Expect(errors).To(Not(BeNil()))
}

func buildTestingModel(g *GomegaWithT) *CoreModel {
	sourceDataSet := buildTestingModelDataSet(g)

	parametersUnderTest := parameters.Map{}

	modelUnderTest := buildModelUnderTest(sourceDataSet, parametersUnderTest, g)
	return modelUnderTest
}

func buildBoundedTestingModel(g *GomegaWithT) *CoreModel {
	sourceDataSet := buildTestingModelDataSet(g)

	parametersUnderTest := parameters.Map{
		"MaximumImplementationCost": expectedMaximumImplementationCost,
	}

	modelUnderTest := buildModelUnderTest(sourceDataSet, parametersUnderTest, g)
	return modelUnderTest
}

func buildInvalidBoundedTestingModel(g *GomegaWithT) error {
	sourceDataSet := buildTestingModelDataSet(g)

	parametersUnderTest := parameters.Map{
		"MaximumImplementationCost": expectedMaximumImplementationCost,
		"MaximumSedimentProduction": expectedMaximumSedimentProduction,
	}

	errors := buildInvalidModelUnderTest(sourceDataSet, parametersUnderTest, g)
	return errors
}

func buildInvalidModelUnderTest(sourceDataSet *csv.DataSet, parametersUnderTest parameters.Map, g *GomegaWithT) error {
	modelUnderTest := NewCoreModel().
		WithSourceDataSet(sourceDataSet).
		WithParameters(parametersUnderTest)

	parameterErrors := modelUnderTest.ParameterErrors()
	g.Expect(parameterErrors).To(Not(BeNil()))

	return parameterErrors
}

func buildModelUnderTest(sourceDataSet *csv.DataSet, parametersUnderTest parameters.Map, g *GomegaWithT) *CoreModel {
	modelUnderTest := NewCoreModel().
		WithSourceDataSet(sourceDataSet).
		WithParameters(parametersUnderTest)
	modelUnderTest.SetId("ModelUnderTest")

	parameterErrors := modelUnderTest.ParameterErrors()
	g.Expect(parameterErrors).To(BeNil())

	modelUnderTest.AddObserver(loggers.DefaultTestingAnnealingObserver)

	modelUnderTest.Initialise(model2.AsIs)
	return modelUnderTest
}

func buildTestingModelDataSet(g *GomegaWithT) *csv.DataSet {
	sourceDataSet := csv.NewDataSet("CatchmentModel")
	loadError := sourceDataSet.Load("testdata/TestingModel.csv")

	g.Expect(loadError).To(BeNil())
	return sourceDataSet
}

func verifyActionToggle(t *testing.T, modelUnderTest *CoreModel, planningUnit planningunit.Id, actionType action.ManagementActionType, g *GomegaWithT) {
	firstSolution := new(solution.SolutionBuilder).
		WithId("testingBuilder").
		ForModel(modelUnderTest).
		Build()

	modelUnderTest.ToggleAction(planningUnit, actionType)
	modelUnderTest.AcceptChange()

	modelUnderTest.ToggleAction(planningUnit, actionType)
	modelUnderTest.AcceptChange()

	secondSolution := new(solution.SolutionBuilder).
		WithId("testingBuilder").
		ForModel(modelUnderTest).
		Build()

	verifySolutionsMatch(t, g, firstSolution, secondSolution)
}

func verifySolutionsMatch(t *testing.T, g *GomegaWithT, firstSolution *solution.Solution, secondSolution *solution.Solution) {
	matchErrors := firstSolution.MatchErrors(secondSolution)
	if matchErrors != nil {
		t.Log(matchErrors)
	}

	g.Expect(matchErrors).To(BeNil())
}

func solutionVariable(solution *solution.Solution, variableName string) *variable.EncodeableDecisionVariable {
	for _, currSolution := range solution.DecisionVariables {
		if currSolution.Name == variableName {
			return &currSolution
		}
	}
	return nil
}
