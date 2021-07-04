// Copyright (c) 2019 Australian Rivers Institute.

package catchment

import (
	errors2 "errors"
	"github.com/LindsayBradford/crem/internal/pkg/dataset"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/csv"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/LindsayBradford/crem/internal/pkg/dataset/excel"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
	baseParameters "github.com/LindsayBradford/crem/internal/pkg/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/pkg/threading"
)

var _ model.Model = NewModel()

func NewModel() *Model {
	newModel := new(Model)

	newModel.SetName("CatchmentModel")

	newModel.parameters.Initialise()
	newModel.managementActions.Initialise()
	newModel.ContainedDecisionVariables.Initialise()

	return newModel
}

type Model struct {
	sourceDataLoaded   bool
	sourceDataSet      dataset.DataSet
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

func (m *Model) Initialise(initialisationType model.InitialisationType) {
	m.note("Initialising")

	if !m.sourceDataLoaded {
		loadError := m.loadSourceDataSet()
		if loadError != nil {
			m.parameters.AddValidationErrorMessage(loadError.Error())
			return
		}
		m.sourceDataLoaded = true
	}
	m.CoreModel.Initialise(initialisationType)
}

func (m *Model) Randomize() {
	m.note("Randomizing")
	m.CoreModel.Randomize()
}

func (m *Model) RandomlyInitialiseActions() {
	m.initialising = true
	m.CoreModel.InitialiseActions(model.Random)
	m.CoreModel.Randomize()
	m.initialising = false
}

func (m *Model) loadSourceDataSet() error {
	dataSourcePath := m.deriveDataSourcePath()
	pathExtension := strings.ToLower(path.Ext(dataSourcePath))
	switch pathExtension {
	case ".csv":
		return m.loadCsvSourceDataSet(dataSourcePath)
	case ".xlsx":
		return m.loadExcelSourceDataSet(dataSourcePath)
	default:
		return errors2.New("Source data file not supported: Initialisation failed")
	}
}

func (m *Model) loadCsvSourceDataSet(dataSourcePath string) error {
	dataSet := csv.NewDataSet("DataSetImpl")

	loadError := dataSet.Load(dataSourcePath)
	if loadError != nil {
		return loadError
	}

	m.sourceDataSet = dataSet
	m.WithSourceDataSet(m.sourceDataSet)

	return nil
}

func (m *Model) loadExcelSourceDataSet(dataSourcePath string) error {
	dataSet := excel.NewDataSet("DataSetImpl", m.oleFunctionWrapper)

	loadError := dataSet.Load(dataSourcePath)
	if loadError != nil {
		return loadError
	}

	m.sourceDataSet = dataSet
	m.WithSourceDataSet(m.sourceDataSet)

	return nil
}

func (m *Model) deriveDataSourcePath() string {
	relativeFilePath := m.parameters.GetString(parameters.DataSourcePath)
	workingDirectory, _ := os.Getwd()
	return filepath.Join(workingDirectory, relativeFilePath)
}

func (m *Model) DeepClone() model.Model {
	clone := *m
	clone.managementActions.SetRandomNumberGenerator(rand.NewTimeSeeded())
	clone.Initialise(model.Unchanged)

	return &clone
}

func (m *Model) TearDown() {
	m.sourceDataSet.Teardown()
}
