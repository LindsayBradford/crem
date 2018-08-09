// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package components

import (
	"math"
	"math/rand"
	"time"

	. "github.com/LindsayBradford/crm/annealing/solution"
)

var (
	randomNumberGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type SimpleExcelSolutionExplorer struct {
	BaseSolutionExplorer

	dataSourcePath string
	penalty        float64

	annealingData *annealingTable
	trackingData  *trackingTable

	temperature float64

	previousPlanningUnitChanged uint64
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

func (explorer *SimpleExcelSolutionExplorer) Initialise() {
	explorer.BaseSolutionExplorer.Initialise()

	initialiseDataSource(explorer.dataSourcePath)
	explorer.LogHandler().Info("Opening Excel workbook [" + explorer.dataSourcePath + "] as data source")

	explorer.LogHandler().Debug("Retrieving annealing data from workbook")
	explorer.annealingData = retrieveAnnealingTableFromWorkbook()

	currentPenalty := explorer.deriveTotalPenalty()
	currentCost := explorer.deriveFeatureCost()

	explorer.SetObjectiveValue(currentCost*0.8 + currentPenalty)

	explorer.LogHandler().Debug("Clearing tracking data from workbook")
	explorer.trackingData = initialiseTrackingTable()

	explorer.LogHandler().Info("Data retrieved from workbook [" + explorer.dataSourcePath + "]")
}

func (explorer *SimpleExcelSolutionExplorer) WithPenalty(penalty float64) *SimpleExcelSolutionExplorer {
	explorer.penalty = penalty
	return explorer
}

func (explorer *SimpleExcelSolutionExplorer) WithInputFile(inputFilePath string) *SimpleExcelSolutionExplorer {
	explorer.dataSourcePath = inputFilePath
	return explorer
}

func (explorer *SimpleExcelSolutionExplorer) TearDown() {
	explorer.BaseSolutionExplorer.TearDown()
	explorer.saveDataToWorkbookAndClose()
	destroyExcelHandler()
}

func (explorer *SimpleExcelSolutionExplorer) saveDataToWorkbookAndClose() {
	explorer.LogHandler().Info("Storing data to workbook [" + explorer.dataSourcePath + "]")
	storeAnnealingTableToWorkbook(explorer.annealingData)
	storeTrackingTableToWorkbook(explorer.trackingData)

	explorer.LogHandler().Debug("Saving workbook [" + explorer.dataSourcePath + "]")
	saveAndCloseWorkbook()

	explorer.LogHandler().Debug("Workbook [" + explorer.dataSourcePath + "] closed")
}

func (explorer *SimpleExcelSolutionExplorer) TryRandomChange(temperature float64) {
	explorer.temperature = temperature
	explorer.makeRandomChange(temperature)
	DecideOnWhetherToAcceptChange(explorer, temperature)
}

func (explorer *SimpleExcelSolutionExplorer) makeRandomChange(temperature float64) {
	previousPenalty := explorer.deriveTotalPenalty()
	previousCost := explorer.deriveFeatureCost()

	explorer.previousPlanningUnitChanged = explorer.annealingData.ToggleRandomPlanningUnit()

	currentPenalty := explorer.deriveTotalPenalty()
	currentCost := explorer.deriveFeatureCost()

	changeInObjectiveValue := (currentCost-previousCost)*0.8 + (currentPenalty - previousPenalty)

	explorer.SetChangeInObjectiveValue(changeInObjectiveValue)
	explorer.SetObjectiveValue(explorer.ObjectiveValue() + explorer.ChangeInObjectiveValue())
}

func (explorer *SimpleExcelSolutionExplorer) deriveTotalPenalty() float64 {
	totalPenalty := float64(0)
	for index := 0; index < len(explorer.annealingData.rows); index++ {
		totalPenalty += float64(explorer.annealingData.rows[index].PlanningUnitStatus) * explorer.annealingData.rows[index].Cost
	}
	return math.Max(0, explorer.penalty-totalPenalty)
}

func (explorer *SimpleExcelSolutionExplorer) deriveFeatureCost() float64 {
	dataToCost := explorer.annealingData
	totalFeatureCost := float64(0)
	for index := 0; index < len(dataToCost.rows); index++ {
		totalFeatureCost +=
			float64(dataToCost.rows[index].PlanningUnitStatus) * dataToCost.rows[index].Feature
	}
	return totalFeatureCost
}

func (explorer *SimpleExcelSolutionExplorer) AcceptLastChange() {
	explorer.BaseSolutionExplorer.AcceptLastChange()
	explorer.addTrackerData()
}

func (explorer *SimpleExcelSolutionExplorer) addTrackerData() {
	newRow := new(trackingData)

	newRow.ObjectiveFunctionChange = explorer.ChangeInObjectiveValue()
	newRow.Temperature = explorer.temperature
	newRow.ChangeIsDesirable = explorer.ChangeIsDesirable()
	newRow.ChangeAccepted = explorer.ChangeAccepted()
	newRow.AcceptanceProbability = explorer.AcceptanceProbability()
	newRow.InFirst50 = explorer.deriveSmallPUs()
	newRow.InSecond50 = explorer.deriveLargePUs()
	newRow.TotalCost = explorer.deriveFeatureCost() / 3

	explorer.trackingData.rows = append(explorer.trackingData.rows, *newRow)
}

func (explorer *SimpleExcelSolutionExplorer) deriveSmallPUs() uint64 {
	dataToTraverse := explorer.annealingData
	totalSmallPUs := uint64(0)
	var index = 0
	for index = 0; index < len(dataToTraverse.rows)/2; index++ {
		totalSmallPUs += uint64(dataToTraverse.rows[index].PlanningUnitStatus)
	}
	return totalSmallPUs
}

func (explorer *SimpleExcelSolutionExplorer) deriveLargePUs() uint64 {
	dataToTraverse := explorer.annealingData
	totalLargePUs := uint64(0)
	var index = 0
	for index = len(dataToTraverse.rows) / 2; index < len(dataToTraverse.rows); index++ {
		totalLargePUs += uint64(dataToTraverse.rows[index].PlanningUnitStatus)
	}
	return totalLargePUs
}

func (explorer *SimpleExcelSolutionExplorer) RevertLastChange() {
	explorer.annealingData.TogglePlanningUnitStatusAtIndex(explorer.previousPlanningUnitChanged)
	explorer.SetObjectiveValue(explorer.ObjectiveValue() - explorer.ChangeInObjectiveValue())
	explorer.SetChangeInObjectiveValue(0)
	explorer.addTrackerData()
	explorer.BaseSolutionExplorer.RevertLastChange()
}

func (explorer *SimpleExcelSolutionExplorer) WithName(name string) *SimpleExcelSolutionExplorer {
	explorer.BaseSolutionExplorer.SetName(name)
	return explorer
}
