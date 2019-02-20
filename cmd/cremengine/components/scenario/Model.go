// Copyright (c) 2018 Australian Rivers Institute.

package scenario

import (
	"math"
	"os"
	"path/filepath"
	"strconv"

	"github.com/LindsayBradford/crem/cmd/cremengine/components/scenario/actions"
	"github.com/LindsayBradford/crem/cmd/cremengine/components/scenario/parameters"
	"github.com/LindsayBradford/crem/cmd/cremengine/components/scenario/variables"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/model/action"
	baseParameters "github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/excel"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/observer"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/pkg/name"
	"github.com/LindsayBradford/crem/pkg/strings"
	"github.com/LindsayBradford/crem/pkg/threading"
	"github.com/pkg/errors"
)

const (
	PlanningUnitsTableName = "PlanningUnits"
)

func NewModel() *Model {
	newModel := new(Model)
	newModel.SetName("CatchmentModel")

	newModel.parameters.Initialise()
	newModel.managementActions.Initialise()
	newModel.DecisionVariables().Initialise()

	return newModel
}

type Model struct {
	name.ContainedName
	name.ContainedIdentifier
	observer.ContainedEventNotifier

	parameters parameters.Parameters

	managementActions action.ManagementActions
	variables.ContainedDecisionVariables

	sedimentLoad *variables.SedimentLoad

	dataSet *excel.DataSet
}

func (m *Model) WithName(name string) *Model {
	m.SetName(name)
	return m
}

func (m *Model) WithOleFunctionWrapper(wrapper threading.MainThreadFunctionWrapper) *Model {
	m.dataSet = excel.NewDataSet("CatchmentDataSet", wrapper)
	return m
}

func (m *Model) WithParameters(params baseParameters.Map) *Model {
	m.SetParameters(params)
	return m
}

func (m *Model) SetParameters(params baseParameters.Map) error {
	m.parameters.Merge(params)

	return m.parameters.ValidationErrors()
}

func (m *Model) ParameterErrors() error {
	return m.parameters.ValidationErrors()
}

func (m *Model) Initialise() {
	m.note("Initialising")

	dataSourcePath := m.deriveDataSourcePath()

	m.dataSet.Load(dataSourcePath)

	planningUnitTable, tableError := m.dataSet.Table(PlanningUnitsTableName)
	if tableError != nil {
		panic(errors.New("Expected data set supplied to have a [" + PlanningUnitsTableName + "] table"))
	}

	csvPlanningUnitTable, ok := planningUnitTable.(*tables.CsvTable)
	if !ok {
		panic(errors.New("Expected data set table [" + PlanningUnitsTableName + "] to be a CSV type"))
	}

	m.sedimentLoad = new(variables.SedimentLoad).Initialise(csvPlanningUnitTable, m.parameters)

	// TODO: Create other sediment management actions
	riverBankRestorations := new(actions.RiverBankRestorations).Initialise(csvPlanningUnitTable, m.parameters)
	for _, action := range riverBankRestorations.ManagementActions() {
		action.Subscribe(m)
		action.Subscribe(m.sedimentLoad)
		m.managementActions.Add(action)
	}

	m.managementActions.RandomlyToggleAllActivations()
	m.DecisionVariables().Add(&m.sedimentLoad.VolatileDecisionVariable)
}

func (m *Model) AcceptChange() {
	m.note("Accepting Change")
	m.sedimentLoad.Accept()
}

func (m *Model) RevertChange() {
	m.note("Reverting Change")
	m.sedimentLoad.Revert()
	m.managementActions.UndoLastActivationToggleUnobserved()
}

func (m *Model) deriveDataSourcePath() string {
	relativeFilePath := m.parameters.GetString(parameters.DataSourcePath)
	workingDirectory, _ := os.Getwd()
	return filepath.Join(workingDirectory, relativeFilePath)
}

func (m *Model) TearDown() {
	//  TODO: Do I need to do any special shutdown behaviour?
}

func (m *Model) TryRandomChange() {
	m.note("Trying Random Change")
	m.managementActions.RandomlyToggleOneActivation()
}

func (m *Model) Observe(action action.ManagementAction) {
	m.noteAppliedManagementAction(action)
}

func (m *Model) noteAppliedManagementAction(action action.ManagementAction) {
	var builder strings.FluentBuilder
	builder.Add("Type [", string(action.Type()),
		"], ", "PlanningUnit [", action.PlanningUnit(),
		"], ", "Active [", strconv.FormatBool(action.IsActive()),
		"]",
	)

	event := observer.Event{EventType: observer.ManagementAction, EventSource: m, Note: builder.String()}
	m.EventNotifier().NotifyObserversOfEvent(event)
}

func (m *Model) note(text string) {
	event := observer.Event{EventType: observer.Note, EventSource: m, Note: text}
	m.EventNotifier().NotifyObserversOfEvent(event)
}

func (m *Model) capChangeOverRange(value float64) float64 {
	return math.Max(0, value)
}

func (m *Model) objectiveValue() float64 {
	return m.sedimentLoad.Value()
}

func (m *Model) setObjectiveValue(value float64) {
	m.sedimentLoad.SetTemporaryValue(value)
}

func (m *Model) DeepClone() model.Model {
	clone := *m
	clone.managementActions.SetRandomNumberGenerator(rand.NewTimeSeeded())
	return &clone
}
