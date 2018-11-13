// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package components

import (
	"math"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	. "github.com/LindsayBradford/crem/annealing/explorer"
)

var (
	randomNumberGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))
)

const (
	_              = iota
	Penalty string = "Penalty"
)

type SimpleExcelSolutionExplorer struct {
	SingleObjectiveAnnealableExplorer

	dataSourcePath string
	parameters     *Parameters

	annealingData *annealingTable
	trackingData  *trackingTable

	temperature float64

	previousPlanningUnitChanged uint64
	excelDataAdapter            ExcelDataAdapter
	oleWrapper                  func(f func())
}

type annealingTable struct {
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
	return (uint64)(randomNumberGenerator.Intn(tableSize))
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

func (e *SimpleExcelSolutionExplorer) Initialise() {
	e.SingleObjectiveAnnealableExplorer.Initialise()
	e.parameters.Initialise()
	e.excelDataAdapter.Initialise().WithOleFunctionWrapper(e.oleWrapper)

	e.excelDataAdapter.initialiseDataSource(e.dataSourcePath)
	e.LogHandler().Info(e.ScenarioId() + ": Opening Excel workbook [" + e.dataSourcePath + "] as data source")

	e.LogHandler().Debug(e.ScenarioId() + ": Retrieving annealing data from workbook")
	e.annealingData = e.excelDataAdapter.retrieveAnnealingTableFromWorkbook()

	currentPenalty := e.deriveTotalPenalty()
	currentCost := e.deriveFeatureCost()

	e.SetObjectiveValue(currentCost*0.8 + currentPenalty)

	e.LogHandler().Debug(e.ScenarioId() + ": Clearing tracking data from workbook")
	e.trackingData = e.excelDataAdapter.initialiseTrackingTable()

	e.LogHandler().Info(e.ScenarioId() + ": Data retrieved from workbook [" + e.dataSourcePath + "]")
}

func (e *SimpleExcelSolutionExplorer) WithParameters(params map[string]interface{}) *SimpleExcelSolutionExplorer {
	e.parameters = new(Parameters).Initialise()
	e.parameters.Merge(params)
	return e
}

func (e *SimpleExcelSolutionExplorer) ParameterErrors() error {
	return e.parameters.ValidationErrors()
}

func (e *SimpleExcelSolutionExplorer) WithInputFile(inputFilePath string) *SimpleExcelSolutionExplorer {
	e.dataSourcePath = inputFilePath
	return e
}

func (e *SimpleExcelSolutionExplorer) WithOleFunctionWrapper(wrapper func(f func())) *SimpleExcelSolutionExplorer {
	e.oleWrapper = wrapper
	return e
}

func (e *SimpleExcelSolutionExplorer) TearDown() {
	e.SingleObjectiveAnnealableExplorer.TearDown()
	e.saveDataToWorkbookAndClose()
	e.excelDataAdapter.destroyExcelHandler()
}

func (e *SimpleExcelSolutionExplorer) saveDataToWorkbookAndClose() {
	newFileName := toSafeFileName(e.ScenarioId())

	originalFilePath := filepath.Dir(e.excelDataAdapter.absoluteFilePath)
	outputPath := filepath.Join(originalFilePath, newFileName)

	e.LogHandler().Info(e.ScenarioId() + ": Storing data to workbook [" + outputPath + "]")
	e.excelDataAdapter.storeAnnealingTableToWorkbook(e.annealingData)
	e.excelDataAdapter.storeTrackingTableToWorkbook(e.trackingData)

	e.LogHandler().Debug(e.ScenarioId() + ": Saving workbook [" + outputPath + "]")
	e.excelDataAdapter.saveAndCloseWorkbookAs(outputPath)

	e.LogHandler().Debug(e.ScenarioId() + ": Workbook [" + outputPath + "] closed")
}

func toSafeFileName(possiblyUnsafeFilePath string) (response string) {
	response = strings.Replace(possiblyUnsafeFilePath, " ", "", -1)
	response = strings.Replace(response, "/", "_of_", -1)
	response = response + ".xls"
	return response
}

func (e *SimpleExcelSolutionExplorer) TryRandomChange(temperature float64) {
	e.temperature = temperature
	e.makeRandomChange(temperature)
	e.DecideOnWhetherToAcceptChange(temperature, e.AcceptLastChange, e.RevertLastChange)
}

func (e *SimpleExcelSolutionExplorer) makeRandomChange(temperature float64) {
	previousPenalty := e.deriveTotalPenalty()
	previousCost := e.deriveFeatureCost()

	e.previousPlanningUnitChanged = e.annealingData.ToggleRandomPlanningUnit()

	currentPenalty := e.deriveTotalPenalty()
	currentCost := e.deriveFeatureCost()

	changeInObjectiveValue := (currentCost-previousCost)*0.8 + (currentPenalty - previousPenalty)

	e.SetChangeInObjectiveValue(changeInObjectiveValue)
	e.SetObjectiveValue(e.ObjectiveValue() + e.ChangeInObjectiveValue())
}

func (e *SimpleExcelSolutionExplorer) deriveTotalPenalty() float64 {
	penalty := e.parameters.GetFloat64(Penalty)

	totalPenalty := float64(0)
	for index := 0; index < len(e.annealingData.rows); index++ {
		totalPenalty += float64(e.annealingData.rows[index].PlanningUnitStatus) * e.annealingData.rows[index].Cost
	}
	return math.Max(0, penalty-totalPenalty)
}

func (e *SimpleExcelSolutionExplorer) deriveFeatureCost() float64 {
	dataToCost := e.annealingData
	totalFeatureCost := float64(0)
	for index := 0; index < len(dataToCost.rows); index++ {
		totalFeatureCost +=
			float64(dataToCost.rows[index].PlanningUnitStatus) * dataToCost.rows[index].Feature
	}
	return totalFeatureCost
}

func (e *SimpleExcelSolutionExplorer) AcceptLastChange() {
	e.SingleObjectiveAnnealableExplorer.AcceptLastChange()
	e.addTrackerData()
}

func (e *SimpleExcelSolutionExplorer) addTrackerData() {
	newRow := new(trackingData)

	newRow.ObjectiveFunctionChange = e.ChangeInObjectiveValue()
	newRow.Temperature = e.temperature
	newRow.ChangeIsDesirable = e.ChangeIsDesirable()
	newRow.ChangeAccepted = e.ChangeAccepted()
	newRow.AcceptanceProbability = e.AcceptanceProbability()
	newRow.InFirst50 = e.deriveSmallPUs()
	newRow.InSecond50 = e.deriveLargePUs()
	newRow.TotalCost = e.deriveFeatureCost() / 3

	e.trackingData.rows = append(e.trackingData.rows, *newRow)
}

func (e *SimpleExcelSolutionExplorer) deriveSmallPUs() uint64 {
	dataToTraverse := e.annealingData
	totalSmallPUs := uint64(0)
	var index = 0
	for index = 0; index < len(dataToTraverse.rows)/2; index++ {
		totalSmallPUs += uint64(dataToTraverse.rows[index].PlanningUnitStatus)
	}
	return totalSmallPUs
}

func (e *SimpleExcelSolutionExplorer) deriveLargePUs() uint64 {
	dataToTraverse := e.annealingData
	totalLargePUs := uint64(0)
	var index = 0
	for index = len(dataToTraverse.rows) / 2; index < len(dataToTraverse.rows); index++ {
		totalLargePUs += uint64(dataToTraverse.rows[index].PlanningUnitStatus)
	}
	return totalLargePUs
}

func (e *SimpleExcelSolutionExplorer) RevertLastChange() {
	e.annealingData.TogglePlanningUnitStatusAtIndex(e.previousPlanningUnitChanged)
	e.SetObjectiveValue(e.ObjectiveValue() - e.ChangeInObjectiveValue())
	e.SetChangeInObjectiveValue(0)
	e.addTrackerData()
	e.SingleObjectiveAnnealableExplorer.RevertLastChange()
}

func (e *SimpleExcelSolutionExplorer) WithName(name string) *SimpleExcelSolutionExplorer {
	e.SingleObjectiveAnnealableExplorer.SetName(name)
	return e
}

func (e *SimpleExcelSolutionExplorer) Clone() Explorer {
	clone := *e
	return &clone
}
