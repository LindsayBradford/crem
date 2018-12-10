// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package components

import (
	"errors"
	"math"
	"path/filepath"
	"strings"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/rand"
	"github.com/LindsayBradford/crem/internal/pkg/scenario"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/name"
	"github.com/LindsayBradford/crem/pkg/threading"
)

type SimpleExcelModel struct {
	name.ContainedName
	scenario.ContainedScenarioId
	logging.ContainedLogger

	parameters            Parameters
	decisionVariables     map[string]model.DecisionVariable
	tempDecisionVariables map[string]model.DecisionVariable

	annealingData *annealingTable
	trackingData  *trackingTable

	explorerData                ExplorerData
	previousPlanningUnitChanged uint64
	excelDataAdapter            ExcelDataAdapter
	oleWrapper                  threading.MainThreadFunctionWrapper
}

func NewSimpleExcelModel() *SimpleExcelModel {
	newModel := new(SimpleExcelModel)
	newModel.parameters = *new(Parameters).Initialise()
	newModel.buildDecisionVariables()
	return newModel
}

func (sem *SimpleExcelModel) buildDecisionVariables() {
	sem.decisionVariables = make(map[string]model.DecisionVariable, 1)
	sem.tempDecisionVariables = make(map[string]model.DecisionVariable, 1)

	objectiveValueVar := new(model.DecisionVariableImpl)
	objectiveValueVar.SetName(model.ObjectiveValue)
	sem.decisionVariables[model.ObjectiveValue] = objectiveValueVar

	tempObjectiveValueVar := new(model.DecisionVariableImpl)
	tempObjectiveValueVar.SetName(model.ObjectiveValue)
	sem.tempDecisionVariables[model.ObjectiveValue] = tempObjectiveValueVar
}

func (sem *SimpleExcelModel) WithParameters(params parameters.Map) *SimpleExcelModel {
	sem.parameters.Merge(params)
	return sem
}

func (sem *SimpleExcelModel) ParameterErrors() error {
	return sem.parameters.ValidationErrors()
}

func (sem *SimpleExcelModel) WithOleFunctionWrapper(wrapper threading.MainThreadFunctionWrapper) *SimpleExcelModel {
	sem.oleWrapper = wrapper
	return sem
}

func (sem *SimpleExcelModel) WithName(name string) *SimpleExcelModel {
	sem.SetName(name)
	return sem
}

func (sem *SimpleExcelModel) SetExplorerData(data ExplorerData) {
	sem.explorerData = data
	sem.addTrackerData()
}

func (sem *SimpleExcelModel) Initialise() {
	sem.excelDataAdapter.Initialise().WithOleFunctionWrapper(sem.oleWrapper)

	dataSourcePath := sem.parameters.GetString(DataSourcePath)
	sem.excelDataAdapter.initialiseDataSource(dataSourcePath)

	sem.LogHandler().Info(sem.ScenarioId() + ": Opening Excel workbook [" + dataSourcePath + "] as data source")

	sem.LogHandler().Debug(sem.ScenarioId() + ": Retrieving annealing data from workbook")
	sem.annealingData = sem.excelDataAdapter.retrieveAnnealingTableFromWorkbook()

	currentPenalty := sem.deriveTotalPenalty()
	currentCost := sem.deriveFeatureCost()

	sem.decisionVariables[model.ObjectiveValue].SetValue(currentCost*0.8 + currentPenalty)

	sem.LogHandler().Debug(sem.ScenarioId() + ": Clearing tracking data from workbook")
	sem.trackingData = sem.excelDataAdapter.initialiseTrackingTable()

	sem.LogHandler().Info(sem.ScenarioId() + ": Data retrieved from workbook [" + dataSourcePath + "]")
}

func (sem *SimpleExcelModel) TearDown() {
	sem.saveDataToWorkbookAndClose()
	sem.excelDataAdapter.destroyExcelHandler()
}

func (sem *SimpleExcelModel) saveDataToWorkbookAndClose() {
	newFileName := toSafeFileName(sem.ScenarioId())

	originalFilePath := filepath.Dir(sem.excelDataAdapter.absoluteFilePath)
	outputPath := filepath.Join(originalFilePath, newFileName)

	sem.LogHandler().Info(sem.ScenarioId() + ": Storing data to workbook [" + outputPath + "]")
	sem.excelDataAdapter.storeAnnealingTableToWorkbook(sem.annealingData)
	sem.excelDataAdapter.storeTrackingTableToWorkbook(sem.trackingData)

	sem.LogHandler().Debug(sem.ScenarioId() + ": Saving workbook [" + outputPath + "]")
	sem.excelDataAdapter.saveAndCloseWorkbookAs(outputPath)

	sem.LogHandler().Debug(sem.ScenarioId() + ": Workbook [" + outputPath + "] closed")
}

func toSafeFileName(possiblyUnsafeFilePath string) (response string) {
	response = strings.Replace(possiblyUnsafeFilePath, " ", "", -1)
	response = strings.Replace(response, "/", "_of_", -1)
	response = response + ".xls"
	return response
}

func (sem *SimpleExcelModel) TryRandomChange() {
	previousPenalty := sem.deriveTotalPenalty()
	previousCost := sem.deriveFeatureCost()

	sem.previousPlanningUnitChanged = sem.annealingData.ToggleRandomPlanningUnit()

	currentPenalty := sem.deriveTotalPenalty()
	currentCost := sem.deriveFeatureCost()

	objectiveValue := sem.decisionVariables[model.ObjectiveValue].Value()
	changeInObjectiveValue := (currentCost-previousCost)*0.8 + (currentPenalty - previousPenalty)

	sem.tempDecisionVariables[model.ObjectiveValue].SetValue(objectiveValue + changeInObjectiveValue)
}

func (sem *SimpleExcelModel) deriveTotalPenalty() float64 {
	penalty := sem.parameters.GetFloat64(Penalty)

	totalPenalty := float64(0)
	for index := 0; index < len(sem.annealingData.rows); index++ {
		totalPenalty += float64(sem.annealingData.rows[index].PlanningUnitStatus) * sem.annealingData.rows[index].Cost
	}
	return math.Max(0, penalty-totalPenalty)
}

func (sem *SimpleExcelModel) deriveFeatureCost() float64 {
	dataToCost := sem.annealingData
	totalFeatureCost := float64(0)
	for index := 0; index < len(dataToCost.rows); index++ {
		totalFeatureCost +=
			float64(dataToCost.rows[index].PlanningUnitStatus) * dataToCost.rows[index].Feature
	}
	return totalFeatureCost
}

func (sem *SimpleExcelModel) AcceptChange() {
	// DeliberatelyDoesNothing
}

func (sem *SimpleExcelModel) addTrackerData() {
	newRow := new(trackingData)

	if change, err := sem.DecisionVariableChange(model.ObjectiveValue); err != nil {
		newRow.ObjectiveFunctionChange = change
	}
	newRow.Temperature = sem.explorerData.Temperature
	newRow.ChangeIsDesirable = sem.explorerData.ChangeIsDesirable
	newRow.ChangeAccepted = sem.explorerData.ChangeAccepted
	newRow.AcceptanceProbability = sem.explorerData.AcceptanceProbability
	newRow.InFirst50 = sem.deriveSmallPUs()
	newRow.InSecond50 = sem.deriveLargePUs()
	newRow.TotalCost = sem.deriveFeatureCost() / 3

	sem.trackingData.rows = append(sem.trackingData.rows, *newRow)
}

func (sem *SimpleExcelModel) deriveSmallPUs() uint64 {
	dataToTraverse := sem.annealingData
	totalSmallPUs := uint64(0)
	for index := 0; index < len(dataToTraverse.rows)/2; index++ {
		totalSmallPUs += uint64(dataToTraverse.rows[index].PlanningUnitStatus)
	}
	return totalSmallPUs
}

func (sem *SimpleExcelModel) deriveLargePUs() uint64 {
	dataToTraverse := sem.annealingData
	totalLargePUs := uint64(0)
	for index := len(dataToTraverse.rows) / 2; index < len(dataToTraverse.rows); index++ {
		totalLargePUs += uint64(dataToTraverse.rows[index].PlanningUnitStatus)
	}
	return totalLargePUs
}

func (sem *SimpleExcelModel) RevertChange() {
	sem.annealingData.TogglePlanningUnitStatusAtIndex(sem.previousPlanningUnitChanged)
	sem.copyDecisionVarValueToTemp(model.ObjectiveValue)
}

func (sem *SimpleExcelModel) copyDecisionVarValueToTemp(varName string) {
	sem.tempDecisionVariables[varName].SetValue(
		sem.decisionVariables[varName].Value(),
	)
}

func (sem *SimpleExcelModel) copyTempDecisionVarValueToActual(varName string) {
	sem.decisionVariables[varName].SetValue(
		sem.tempDecisionVariables[varName].Value(),
	)
}

func (sem *SimpleExcelModel) DecisionVariable(name string) (model.DecisionVariable, error) {
	if variable, found := sem.decisionVariables[name]; found == true {
		return variable, nil
	}
	return model.NullDecisionVariable, errors.New("decision variable [" + name + "] not defined for model [" + sem.Name() + " ].")
}

func (sem *SimpleExcelModel) DecisionVariableChange(variableName string) (float64, error) {
	decisionVariable, foundActual := sem.decisionVariables[variableName]
	tmpDecisionVar, foundTemp := sem.tempDecisionVariables[decisionVariable.Name()]
	if !foundActual || !foundTemp {
		return 0, errors.New("no temporary decision variable of name [" + decisionVariable.Name() + "] in model [" + sem.Name() + "].")
	}

	difference := tmpDecisionVar.Value() - decisionVariable.Value()
	return difference, nil
}

func (sem *SimpleExcelModel) DeepClone() model.Model {
	clone := *sem
	return &clone
}

type ExplorerData struct {
	Temperature           float64
	ChangeIsDesirable     bool
	ChangeAccepted        bool
	AcceptanceProbability float64
}

type annealingTable struct {
	rand.ContainedRand
	rows []annealingData
}

type annealingData struct {
	Cost               float64
	Feature            float64
	PlanningUnitStatus InclusionStatus
}

type InclusionStatus uint64

const (
	OUT InclusionStatus = 0
	IN  InclusionStatus = 1
)

func (table *annealingTable) ToggleRandomPlanningUnit() (rowIndexToggled uint64) {
	rowIndexToggled = table.SelectRandomPlanningUnit()
	table.TogglePlanningUnitStatusAtIndex(rowIndexToggled)
	return
}

func (table *annealingTable) SelectRandomPlanningUnit() uint64 {
	tableSize := len(table.rows)
	return (uint64)(table.RandomNumberGenerator().Intn(tableSize))
}

func (table *annealingTable) TogglePlanningUnitStatusAtIndex(index uint64) {
	newStatusValue := (InclusionStatus)((table.rows[index].PlanningUnitStatus + 1) % 2)
	table.setPlanningUnitStatusAtIndex(newStatusValue, index)
}

func (table *annealingTable) setPlanningUnitStatusAtIndex(status InclusionStatus, index uint64) {
	table.rows[index].PlanningUnitStatus = status
}

type trackingTable struct {
	headings []trackingTableHeadings
	rows     []trackingData
}

type trackingTableHeadings int

const (
	UNKNOWN trackingTableHeadings = iota
	ObjFuncChange
	Temperature
	ChangeIsDesirable
	AcceptanceProbability
	ChangeAccepted
	InFirst50
	InSecond50
	TotalCost
)

func (heading trackingTableHeadings) String() string {
	columnNames := [...]string{
		"<UnknownHeader>",
		"ObjFuncChange",
		"Temperature",
		"ChangeIsDesirable",
		"AcceptanceProbability",
		"ChangeAccepted",
		"InFirst50",
		"InSecond50",
		"TotalCost",
	}
	if heading < ObjFuncChange || heading > TotalCost {
		return columnNames[UNKNOWN]
	}

	return columnNames[heading]
}

func (heading trackingTableHeadings) Index() uint {
	return uint(heading)
}

type trackingData struct {
	ObjectiveFunctionChange float64
	Temperature             float64
	ChangeIsDesirable       bool
	AcceptanceProbability   float64
	ChangeAccepted          bool
	InFirst50               uint64
	InSecond50              uint64
	TotalCost               float64
}
