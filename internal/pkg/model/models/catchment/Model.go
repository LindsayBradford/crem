// Copyright (c) 2019 Australian Rivers Institute.

package catchment

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/variables"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"github.com/LindsayBradford/crem/pkg/strings"

	"github.com/LindsayBradford/crem/internal/pkg/dataset/excel"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	baseParameters "github.com/LindsayBradford/crem/internal/pkg/parameters"
	"github.com/LindsayBradford/crem/pkg/threading"
)

func NewModel() *Model {
	newModel := new(Model)

	newModel.SetName("CatchmentModel")
	newModel.SetEventNotifier(loggers.DefaultTestingEventNotifier)

	newModel.parameters.Initialise()
	newModel.managementActions.Initialise()
	newModel.ContainedDecisionVariables.Initialise()

	return newModel
}

type Model struct {
	sourceDataSet      excel.DataSet
	oleFunctionWrapper threading.MainThreadFunctionWrapper
	CoreModel
}

func (m *Model) WithOleFunctionWrapper(wrapper threading.MainThreadFunctionWrapper) *Model {
	m.oleFunctionWrapper = wrapper
	return m
}

func (m *Model) WithParameters(params baseParameters.Map) *Model {
	m.CoreModel.SetParameters(params)
	return m
}

func (m *Model) Initialise() {
	m.note("Initialising")

	m.loadSourceDataSet()
	m.CoreModel.Initialise()
	m.RandomlyInitialiseActions()

	// the SedimentVsCost variable needs values in both SedimentProduced and ImplementationCost in order to
	// build an effective weighted unit-scaled vector.  If one of them is 0 (where implementation cost starts),
	// we get numbers that report as infinity, breaking stuff.  I can't build that vector UNTIL I can guarantee
	// that I can build a unit-scaled vector without issue.  For now, this is a spot where I can keep it, though it's
	// conceptually a very poor location.

	m.buildSedimentVsCostDecisionVariable()
}

func (m *Model) buildSedimentVsCostDecisionVariable() {
	sedimentProduction := m.ContainedDecisionVariables.Variable(variables.SedimentProductionVariableName)
	implementationCost := m.ContainedDecisionVariables.Variable(variables.ImplementationCostVariableName)

	sedimentWeight := m.parameters.GetFloat64(parameters.SedimentProductionDecisionWeight)
	implementationCostWeight := m.parameters.GetFloat64(parameters.ImplementationCostDecisionWeight)

	sedimentVsCost, buildError := new(variables.SedimentVsCost).
		Initialise().
		WithObservers(m).
		WithWeightedVariable(sedimentProduction, sedimentWeight).
		WithWeightedVariable(implementationCost, implementationCostWeight).
		Build()

	if buildError != nil {
		panic(buildError)
	}

	noteBuilder := new(strings.FluentBuilder).
		Add(sedimentProduction.Name(), " weight = ", strconv.FormatFloat(sedimentWeight, 'f', 3, 64), ", ").
		Add(implementationCost.Name(), " weight = ", strconv.FormatFloat(implementationCostWeight, 'f', 3, 64))

	m.ObserveDecisionVariableWithNote(sedimentProduction, " Initial Value")
	m.ObserveDecisionVariableWithNote(implementationCost, " Initial Value")
	m.ObserveDecisionVariableWithNote(sedimentVsCost, noteBuilder.String())

	m.ContainedDecisionVariables.Add(sedimentVsCost)
}

func (m *Model) loadSourceDataSet() {
	m.sourceDataSet = *excel.NewDataSet("CatchmentDataSet", m.oleFunctionWrapper)
	dataSourcePath := m.deriveDataSourcePath()
	m.sourceDataSet.Load(dataSourcePath)
	m.WithSourceDataSet(&m.sourceDataSet)
}

func (m *Model) deriveDataSourcePath() string {
	relativeFilePath := m.parameters.GetString(parameters.DataSourcePath)
	workingDirectory, _ := os.Getwd()
	return filepath.Join(workingDirectory, relativeFilePath)
}

func (m *Model) DeepClone() model.Model {
	clone := *m
	clone.managementActions.SetRandomNumberGenerator(rand.NewTimeSeeded())
	return &clone
}

func (m *Model) TearDown() {
	m.sourceDataSet.Teardown()
}
