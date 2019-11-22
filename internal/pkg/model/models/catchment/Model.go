// Copyright (c) 2019 Australian Rivers Institute.

package catchment

import (
	"os"
	"path/filepath"

	"github.com/LindsayBradford/crem/internal/pkg/dataset/excel"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	baseParameters "github.com/LindsayBradford/crem/internal/pkg/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
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
