// Copyright (c) 2019 Australian Rivers Institute.

package catchment

import (
	"testing"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/annealers"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/csv"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/actions"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables/implementationcost"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables/sedimentproduction"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	"github.com/LindsayBradford/crem/internal/pkg/model/variableNew"
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
	"github.com/LindsayBradford/crem/pkg/math"
	. "github.com/onsi/gomega"
)

const expectedName = "CatchmentModel"

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

	model.Initialise()

	actualActions := model.ManagementActions()
	expectedActionNumber := 16

	g.Expect(len(actualActions)).To(BeNumerically(equalTo, expectedActionNumber))

	actualVariables := *model.DecisionVariables()

	g.Expect(actualVariables).To(HaveKey(implementationcost.ImplementationCostVariableName))
	g.Expect(actualVariables[implementationcost.ImplementationCostVariableName].Value()).To(BeNumerically(equalTo, 0))

	g.Expect(actualVariables).To(HaveKey(sedimentproduction.SedimentProductionVariableName))
}

func TestCoreModel_Initialise_InvalidDataSet_Errors(t *testing.T) {
	g := NewGomegaWithT(t)

	sourceDataSet := csv.NewDataSet("CatchmentModel")
	loadError := sourceDataSet.Load("testdata/InvalidModel.csv")

	g.Expect(loadError).To(BeNil())

	newModelRunner := func() {
		NewCoreModel().WithSourceDataSet(sourceDataSet).Initialise()
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

	solution := new(annealers.SolutionBuilder).
		WithId("testingBuilder").
		ForModel(modelUnderTest).
		Build()

	g.Expect(solution).To(Not(BeNil()))

	verifyPlanningUnitValues(g, solution, implementationcost.ImplementationCostVariableName, 0)
	verifyPlanningUnitValues(g, solution, sedimentproduction.SedimentProductionVariableName, 38310.166)
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

func buildTestingModel(g *GomegaWithT) *CoreModel {
	sourceDataSet := csv.NewDataSet("CatchmentModel")
	loadError := sourceDataSet.Load("testdata/TestingModel.csv")

	g.Expect(loadError).To(BeNil())

	parametersUnderTest := parameters.Map{}

	modelUnderTest := NewCoreModel().
		WithSourceDataSet(sourceDataSet).
		WithParameters(parametersUnderTest)

	parameterErrors := modelUnderTest.ParameterErrors()
	g.Expect(parameterErrors).To(BeNil())

	modelUnderTest.Initialise()
	return modelUnderTest
}

func verifyActionToggle(t *testing.T, modelUnderTest *CoreModel, planningUnit planningunit.Id, actionType action.ManagementActionType, g *GomegaWithT) {
	firstSolution := new(annealers.SolutionBuilder).
		WithId("testingBuilder").
		ForModel(modelUnderTest).
		Build()

	modelUnderTest.ToggleAction(planningUnit, actionType)
	modelUnderTest.ToggleAction(planningUnit, actionType)

	secondSolution := new(annealers.SolutionBuilder).
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

func solutionVariable(solution *solution.Solution, variableName string) *variableNew.EncodeableDecisionVariable {
	for _, currSolution := range solution.DecisionVariables {
		if currSolution.Name == variableName {
			return &currSolution
		}
	}
	return nil
}
