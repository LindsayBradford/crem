package api

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment"
	"github.com/LindsayBradford/crem/internal/pkg/parameters"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

const timeThreshold = time.Millisecond

func TestModePool_New(t *testing.T) {
	g := NewGomegaWithT(t)
	modelUnderTest := buildTestModel(g)

	poolUnderTest := NewSolutionPool(modelUnderTest)

	g.Expect(poolUnderTest.Size()).To(BeNumerically("==", 1))
	g.Expect(poolUnderTest.Solution(AsIs)).To(Not(BeNil()))

	modelUnderTest.Initialise(model.AsIs)
	modelUnderTest.ReplaceAttribute("ParetoFrontMember", false)
	solutionUnderTest := poolUnderTest.deriveSolutionFrom(modelUnderTest)

	poolAsIsSolution := poolUnderTest.Solution(AsIs)
	g.Expect(poolAsIsSolution).To(Equal(solutionUnderTest))
}

func TestModePool_AddSolution(t *testing.T) {
	g := NewGomegaWithT(t)
	referenceModel := buildTestModel(g)
	poolUnderTest := NewSolutionPool(referenceModel)

	referenceModel.Initialise(model.AsIs)
	referenceModel.SetManagementAction(0, true)
	referenceModel.SetManagementAction(5, true)
	referenceModel.SetManagementAction(7, true)
	referenceModel.ReplaceAttribute("ParetoFrontMember", true)

	expectedSolution := poolUnderTest.deriveSolutionFrom(referenceModel)

	labelUnderTest := SolutionPoolLabel("test")
	expectedSummary := "a test solution"

	poolUnderTest.AddSolution(labelUnderTest, "A1", expectedSummary) // A1 = 100001010000 binary

	g.Expect(poolUnderTest.Summary(labelUnderTest)).To(Equal(expectedSummary))
	g.Expect(poolUnderTest.Solution(labelUnderTest)).To(Equal(expectedSolution))
}

func TestModeContainer_New(t *testing.T) {
	g := NewGomegaWithT(t)

	modelUnderTest := buildTestModel(g)
	solutionUnderTest := solutionOfModel(modelUnderTest)
	modelContainerUnderTest := NewSolutionContainer(solutionUnderTest, "testSummary")

	g.Expect(modelContainerUnderTest.LastUpdated).To(BeTemporally("~", time.Now(), timeThreshold))
	g.Expect(modelContainerUnderTest.Solution).To(Equal(solutionUnderTest))
	g.Expect(modelContainerUnderTest.Summary).To(Equal("testSummary"))
}

func solutionOfModel(model model.Model) *solution.Solution {
	return solutionBuilder.WithId(model.Id()).ForModel(model).Build()
}

func buildTestModel(g *GomegaWithT) *catchment.Model {
	parametersUnderTest := parameters.Map{"DataSourcePath": "testdata/ValidModel.csv"}

	modelUnderTest := catchment.NewModel().WithParameters(parametersUnderTest)
	modelUnderTest.SetId("ModelUnderTest")

	parameterErrors := modelUnderTest.ParameterErrors()
	g.Expect(parameterErrors).To(BeNil())

	//modelUnderTest.AddObserver(loggers.DefaultTestingAnnealingObserver)

	modelUnderTest.Initialise(model.Random)
	return modelUnderTest
}
