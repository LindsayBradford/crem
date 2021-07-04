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

	poolUnderTest := NewModelPool(modelUnderTest)

	g.Expect(poolUnderTest.Size()).To(BeNumerically("==", 2))
	g.Expect(poolUnderTest.Model(AsIs)).To(Not(BeNil()))
	g.Expect(poolUnderTest.Model(Scratchpad)).To(Not(BeNil()))

	modelUnderTest.Initialise(model.AsIs)

	poolAsIsModel := poolUnderTest.Model(AsIs)
	g.Expect(poolAsIsModel.IsEquivalentTo(modelUnderTest)).To(BeTrue())

	poolScratchpadModel := poolUnderTest.Model(Scratchpad)
	g.Expect(poolScratchpadModel.IsEquivalentTo(modelUnderTest)).To(BeTrue())
}

func TestModePool_ScratchpadVariesFromAsIs(t *testing.T) {
	g := NewGomegaWithT(t)
	referenceModel := buildTestModel(g)
	poolUnderTest := NewModelPool(referenceModel)

	scratchpadModel := poolUnderTest.Model(Scratchpad)
	scratchpadModel.SetManagementAction(0, true)

	g.Expect(poolUnderTest.Model(Scratchpad).IsEquivalentTo(poolUnderTest.Model(AsIs))).To(BeFalse())
}

func TestModePool_InstantiateModel(t *testing.T) {
	g := NewGomegaWithT(t)
	referenceModel := buildTestModel(g)
	poolUnderTest := NewModelPool(referenceModel)

	poolUnderTest.InstantiateModel("testModel", "A1", "a test model") // A1 = 100001010000 binary

	g.Expect(poolUnderTest.Model("testModel").IsEquivalentTo(poolUnderTest.Model(AsIs))).To(BeFalse())

	testModel := poolUnderTest.Model("testModel")
	g.Expect(testModel.ManagementActions()[0].IsActive()).To(BeTrue())
	g.Expect(testModel.ManagementActions()[1].IsActive()).To(BeFalse())
	g.Expect(testModel.ManagementActions()[2].IsActive()).To(BeFalse())
	g.Expect(testModel.ManagementActions()[3].IsActive()).To(BeFalse())
	g.Expect(testModel.ManagementActions()[4].IsActive()).To(BeFalse())
	g.Expect(testModel.ManagementActions()[5].IsActive()).To(BeTrue())
	g.Expect(testModel.ManagementActions()[6].IsActive()).To(BeFalse())
	g.Expect(testModel.ManagementActions()[7].IsActive()).To(BeTrue())
	g.Expect(testModel.ManagementActions()[8].IsActive()).To(BeFalse())
	g.Expect(testModel.ManagementActions()[9].IsActive()).To(BeFalse())
	g.Expect(testModel.ManagementActions()[10].IsActive()).To(BeFalse())
	g.Expect(testModel.ManagementActions()[11].IsActive()).To(BeFalse())
}

func TestModeContainer_New(t *testing.T) {
	g := NewGomegaWithT(t)

	modelUnderTest := buildTestModel(g)
	solutionUnderTest := solutionOfModel(modelUnderTest)
	modelContainerUnderTest := NewModelContainer(modelUnderTest)

	g.Expect(modelContainerUnderTest.LastUpdated).To(BeTemporally("~", time.Now(), timeThreshold))
	g.Expect(modelContainerUnderTest.Model.IsEquivalentTo(modelUnderTest)).To(BeTrue())
	g.Expect(modelContainerUnderTest.Solution).To(Equal(solutionUnderTest))
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
