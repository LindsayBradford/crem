// Copyright (c) 2019 Australian Rivers Institute.

package catchment

import (
	"fmt"
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
	. "github.com/onsi/gomega"
)

const expectedName = "CatchmentModel"

const equalTo = "=="
const approx = "~"

const desiredPrecision = 1e-2

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

	verifyPlanningUnitValues(g, solution, implementationcost.ImplementationCost2VariableName, 0)
	verifyPlanningUnitValues(g, solution, sedimentproduction.SedimentProduction2VariableName, 38310.166)
}

func verifyPlanningUnitValues(g *GomegaWithT, solution *solution.Solution, variableName string, expectedValue float64) {
	actualImplementationCost := solutionVariable(solution, variableName)
	g.Expect(actualImplementationCost.Value).To(BeNumerically(equalTo, expectedValue))

	var planningUnitImplementationCost float64
	for _, currValue := range actualImplementationCost.ValuePerPlanningUnit {
		planningUnitImplementationCost += currValue.Value
	}
	g.Expect(actualImplementationCost.Value).To(BeNumerically(approx, planningUnitImplementationCost, desiredPrecision))
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

func TestCoreModel_RiverBankVsHillSlopeRestoration_AsExpected(t *testing.T) {
	g := NewGomegaWithT(t)

	modelUnderTest := buildTestingModel(g)
	planningUnit := planningunit.Id(17)

	modelUnderTest.note("Toggling River Bank Restoration 17 On")
	modelUnderTest.ToggleAction(planningUnit, actions.RiverBankRestorationType)

	solution := new(annealers.SolutionBuilder).
		WithId("testingBuilder").
		ForModel(modelUnderTest).
		Build()

	verifyVariablesMatch(t, g, solution, implementationcost.ImplementationCostVariableName, implementationcost.ImplementationCost2VariableName)
	verifyVariablesMatch(t, g, solution, sedimentproduction.SedimentProductionVariableName, sedimentproduction.SedimentProduction2VariableName)

	modelUnderTest.note("Toggling Hill Slope Restoration 17 On")
	modelUnderTest.ToggleAction(planningUnit, actions.HillSlopeRestorationType)

	solution = new(annealers.SolutionBuilder).
		WithId("testingBuilder").
		ForModel(modelUnderTest).
		Build()

	verifyVariablesMatch(t, g, solution, implementationcost.ImplementationCostVariableName, implementationcost.ImplementationCost2VariableName)
	verifyVariablesMatch(t, g, solution, sedimentproduction.SedimentProductionVariableName, sedimentproduction.SedimentProduction2VariableName)

	modelUnderTest.note("Toggling Hill Slope Restoration 17 Off")
	modelUnderTest.ToggleAction(planningUnit, actions.HillSlopeRestorationType)

	solution = new(annealers.SolutionBuilder).
		WithId("testingBuilder").
		ForModel(modelUnderTest).
		Build()

	verifyVariablesMatch(t, g, solution, implementationcost.ImplementationCostVariableName, implementationcost.ImplementationCost2VariableName)
	verifyVariablesMatch(t, g, solution, sedimentproduction.SedimentProductionVariableName, sedimentproduction.SedimentProduction2VariableName)

	modelUnderTest.note("Toggling River Bank Restoration 17 Off")
	modelUnderTest.ToggleAction(planningUnit, actions.RiverBankRestorationType)

	solution = new(annealers.SolutionBuilder).
		WithId("testingBuilder").
		ForModel(modelUnderTest).
		Build()

	verifyVariablesMatch(t, g, solution, implementationcost.ImplementationCostVariableName, implementationcost.ImplementationCost2VariableName)
	verifyVariablesMatch(t, g, solution, sedimentproduction.SedimentProductionVariableName, sedimentproduction.SedimentProduction2VariableName)
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
	solution := new(annealers.SolutionBuilder).
		WithId("testingBuilder").
		ForModel(modelUnderTest).
		Build()

	modelUnderTest.ToggleAction(planningUnit, actionType)

	verifyVariablesMatch(t, g, solution, implementationcost.ImplementationCostVariableName, implementationcost.ImplementationCost2VariableName)
	verifyVariablesMatch(t, g, solution, sedimentproduction.SedimentProductionVariableName, sedimentproduction.SedimentProduction2VariableName)

	modelUnderTest.ToggleAction(planningUnit, actionType)

	verifyVariablesMatch(t, g, solution, implementationcost.ImplementationCostVariableName, implementationcost.ImplementationCost2VariableName)
	verifyVariablesMatch(t, g, solution, sedimentproduction.SedimentProductionVariableName, sedimentproduction.SedimentProduction2VariableName)
}

func verifyVariablesMatch(t *testing.T, g *GomegaWithT, solution *solution.Solution, firstVariableName string, secondVariableName string) {
	firstVariable := solutionVariable(solution, firstVariableName)
	secondVariable := solutionVariable(solution, secondVariableName)

	for planningUnit := range firstVariable.ValuePerPlanningUnit {
		if firstVariable.ValuePerPlanningUnit[planningUnit].Value != secondVariable.ValuePerPlanningUnit[planningUnit].Value {
			difference := firstVariable.ValuePerPlanningUnit[planningUnit].Value - secondVariable.ValuePerPlanningUnit[planningUnit].Value
			t.Logf("Planning Unit[%v], [%s].Value = [%f], [%s].Value = [%f], Difference = [%f]", firstVariable.ValuePerPlanningUnit[planningUnit].PlanningUnit,
				firstVariableName, firstVariable.ValuePerPlanningUnit[planningUnit].Value,
				secondVariableName, secondVariable.ValuePerPlanningUnit[planningUnit].Value, difference)
		}
		g.Expect(firstVariable.ValuePerPlanningUnit[planningUnit].Value).To(BeNumerically(approx, secondVariable.ValuePerPlanningUnit[planningUnit].Value, desiredPrecision))
	}

	g.Expect(firstVariable.Value).To(BeNumerically(approx, secondVariable.Value, desiredPrecision))
}

func TestCoreModel_RandomVariableValuesMatch_AsExpected(t *testing.T) {
	g := NewGomegaWithT(t)

	modelUnderTest := buildTestingModel(g)

	const loops = 100
	for i := 0; i < loops; i++ {
		iterationMsg := fmt.Sprintf("Invoking Model loop # %d", i)
		modelUnderTest.note(iterationMsg)

		iterationMsg = fmt.Sprintf("Model loop # %d, doing random change", i)
		modelUnderTest.note(iterationMsg)
		modelUnderTest.DoRandomChange()

		doSolution := new(annealers.SolutionBuilder).
			WithId("testingBuilder").
			ForModel(modelUnderTest).
			Build()

		verifyVariablesMatch(t, g, doSolution, implementationcost.ImplementationCostVariableName, implementationcost.ImplementationCost2VariableName)
		verifyVariablesMatch(t, g, doSolution, sedimentproduction.SedimentProductionVariableName, sedimentproduction.SedimentProduction2VariableName)

		iterationMsg = fmt.Sprintf("Model loop # %d, trying and accepting random change", i)
		modelUnderTest.note(iterationMsg)

		modelUnderTest.TryRandomChange()
		modelUnderTest.AcceptChange()

		tryAcceptSolution := new(annealers.SolutionBuilder).
			WithId("testingBuilder").
			ForModel(modelUnderTest).
			Build()

		verifyVariablesMatch(t, g, tryAcceptSolution, implementationcost.ImplementationCostVariableName, implementationcost.ImplementationCost2VariableName)
		verifyVariablesMatch(t, g, tryAcceptSolution, sedimentproduction.SedimentProductionVariableName, sedimentproduction.SedimentProduction2VariableName)

		iterationMsg = fmt.Sprintf("Model loop # %d, trying and reverting random change", i)
		modelUnderTest.note(iterationMsg)

		modelUnderTest.TryRandomChange()
		modelUnderTest.RevertChange()

		tryRevertSolution := new(annealers.SolutionBuilder).
			WithId("testingBuilder").
			ForModel(modelUnderTest).
			Build()

		verifyVariablesMatch(t, g, tryRevertSolution, implementationcost.ImplementationCostVariableName, implementationcost.ImplementationCost2VariableName)
		verifyVariablesMatch(t, g, tryRevertSolution, sedimentproduction.SedimentProductionVariableName, sedimentproduction.SedimentProduction2VariableName)
	}
}

func solutionVariable(solution *solution.Solution, variableName string) *variableNew.EncodeableDecisionVariable {
	for _, currSolution := range solution.DecisionVariables {
		if currSolution.Name == variableName {
			return &currSolution
		}
	}
	return nil
}

func TestCoreModel_InitialisedValuesMatch_AsExpected(t *testing.T) {
	g := NewGomegaWithT(t)

	modelUnderTest := buildTestingModel(g)

	modelUnderTest.initialising = true
	for _, action := range modelUnderTest.managementActions.Actions() {
		modelUnderTest.managementActions.RandomlyInitialiseAction(action)
	}
	modelUnderTest.initialising = false

	solution := new(annealers.SolutionBuilder).
		WithId("testingBuilder").
		ForModel(modelUnderTest).
		Build()

	verifyVariablesMatch(t, g, solution, implementationcost.ImplementationCostVariableName, implementationcost.ImplementationCost2VariableName)
	verifyVariablesMatch(t, g, solution, sedimentproduction.SedimentProductionVariableName, sedimentproduction.SedimentProduction2VariableName)
}
