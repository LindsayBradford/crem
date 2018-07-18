// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package components

import (
	"math"
	"math/rand"
	"time"

	. "github.com/LindsayBradford/crm/annealing/objectives"
)

var (
	randomNumberGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type KnapsackObjectiveManager struct {
	BaseObjectiveManager

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
	Cost     float64
	Feature  float64
	PUStatus InclusionStatus
}

type InclusionStatus uint64

const (
	OUT InclusionStatus = 0
	IN  InclusionStatus = 1
)

func (this *annealingTable) ToggleRandomPlanningUnit() (rowIndexToggled uint64) {
	rowIndexToggled = this.GenerateRandomInOutIndex()
	this.TogglePUStatusAtIndex(rowIndexToggled)
	return
}

func (this *annealingTable) GenerateRandomInOutIndex() uint64 {
	tableSize := len(this.rows)
	return (uint64)(randomNumberGenerator.Intn(tableSize))
}

func (this *annealingTable) TogglePUStatusAtIndex(index uint64) {
	newInOutValue := (InclusionStatus)((this.rows[index].PUStatus + 1) % 2)
	this.setPUStatusAtIndex(newInOutValue, index)
}

func (this *annealingTable) setPUStatusAtIndex(inOutValue InclusionStatus, index uint64) {
	this.rows[index].PUStatus = inOutValue
}

type trackingTable struct {
	rows []trackingData
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

func (this *KnapsackObjectiveManager) Initialise() {
	this.BaseObjectiveManager.Initialise()
	this.penalty = 100 // TODO: make configurable

	this.dataSourcePath = initialiseDataSource()
	this.LogHandler().Info("Opening Excel workbook [" + this.dataSourcePath + "] as data source")

	this.LogHandler().Debug("Retrieving annealing data from workbook")
	this.annealingData = retrieveAnnealingTableFromWorkbook()

	currentPenalty := this.deriveTotalPenalty()
	currentCost := this.deriveFeatureCost()

	this.SetObjectiveValue(currentCost*0.8 + currentPenalty)

	this.LogHandler().Debug("Clearing tracking data from workbook")
	this.trackingData = clearTrackingDataFromWorkbook()

	this.LogHandler().Info("Data retrieved from workbook [" + this.dataSourcePath + "]")
}

func (this *KnapsackObjectiveManager) TearDown() {
	this.BaseObjectiveManager.TearDown()
	this.saveDataToWorkbookAndClose()
	destroyExcelHandler()
}

func (this *KnapsackObjectiveManager) saveDataToWorkbookAndClose() {
	this.LogHandler().Info("Storing data to workbook [" + this.dataSourcePath + "]")
	storeAnnealingTableToWorkbook(this.annealingData)
	storeTrackingTableToWorkbook(this.trackingData)

	this.LogHandler().Debug("Saving workbook [" + this.dataSourcePath + "]")
	saveAndCloseWorkbook()

	this.LogHandler().Debug("Workbook [" + this.dataSourcePath + "] closed")
}

func (this *KnapsackObjectiveManager) TryRandomChange(temperature float64) {
	this.temperature = temperature
	this.makeRandomChange(temperature)
	DecideOnWhetherToAcceptChange(this, temperature)
}

func (this *KnapsackObjectiveManager) makeRandomChange(temperature float64) {
	previousPenalty := this.deriveTotalPenalty()
	previousCost := this.deriveFeatureCost()

	this.previousPlanningUnitChanged = this.annealingData.ToggleRandomPlanningUnit()

	currentPenalty := this.deriveTotalPenalty()
	currentCost := this.deriveFeatureCost()

	changeInObjectiveValue := ((currentCost-previousCost)*0.8 + (currentPenalty - previousPenalty))

	this.SetChangeInObjectiveValue(changeInObjectiveValue)
	this.SetObjectiveValue(this.ObjectiveValue() + this.ChangeInObjectiveValue())
}

func (this *KnapsackObjectiveManager) deriveTotalPenalty() float64 {
	totalPenalty := float64(0)
	for index := 0; index < len(this.annealingData.rows); index++ {
		totalPenalty += float64(this.annealingData.rows[index].PUStatus) * this.annealingData.rows[index].Cost
	}
	return math.Max(0, this.penalty-totalPenalty)
}

func (this *KnapsackObjectiveManager) deriveFeatureCost() float64 {
	totalFeatureCost := float64(0)
	for index := 0; index < len(this.annealingData.rows); index++ {
		totalFeatureCost +=
			float64(this.annealingData.rows[index].PUStatus) * this.annealingData.rows[index].Feature
	}
	return totalFeatureCost
}

func (this *KnapsackObjectiveManager) AcceptLastChange() {
	this.BaseObjectiveManager.AcceptLastChange()
	this.addTrackerData()
}

func (this *KnapsackObjectiveManager) addTrackerData() {
	newRow := new(trackingData)
	newRow.ObjectiveFunctionChange = this.ChangeInObjectiveValue()
	newRow.Temperature = this.temperature
	newRow.ChangeIsDesirable = this.ChangeIsDesirable()
	newRow.ChangeAccepted = this.ChangeAccepted()
	newRow.AcceptanceProbability = this.AcceptanceProbability()
	newRow.InFirst50 = this.deriveSmallPUs()
	newRow.InSecond50 = this.deriveLargePUs()
	newRow.TotalCost = this.deriveFeatureCost() / 3
	this.trackingData.rows = append(this.trackingData.rows, *newRow)
}

func (this *KnapsackObjectiveManager) deriveSmallPUs() uint64 {
	totalSmallPUs := uint64(0)
	var index = 0
	for index = 0; index < len(this.annealingData.rows) / 2; index++ {
		totalSmallPUs += uint64(this.annealingData.rows[index].PUStatus)
	}
	return totalSmallPUs
}

func (this *KnapsackObjectiveManager) deriveLargePUs() uint64 {
	totalLargePUs := uint64(0)
	var index = 0
	for index = len(this.annealingData.rows) / 2; index < len(this.annealingData.rows); index++ {
		totalLargePUs += uint64(this.annealingData.rows[index].PUStatus)
	}
	return totalLargePUs
}

func (this *KnapsackObjectiveManager) RevertLastChange() {
	this.annealingData.TogglePUStatusAtIndex(this.previousPlanningUnitChanged)
	this.SetObjectiveValue(this.ObjectiveValue() - this.ChangeInObjectiveValue())
	this.SetChangeInObjectiveValue(0)
	this.addTrackerData()
	this.BaseObjectiveManager.RevertLastChange()
}
