// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package components

import (
	"math/rand"
	"time"

	. "github.com/LindsayBradford/crm/annealing/objectives"
)

type KnapsackObjectiveManager struct {
	BaseObjectiveManager

	dataSourcePath string

	annealingData *annealingTable
	trackingData  *trackingTable

	temperature float64
}

var (
	randomNumberGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type annealingTable struct {
	rows []annealingData
}

type annealingData struct {
	Cost    float64
	Feature float64
	X       uint64
	Y       uint64
	InOut   uint64
	Reached float64
	InCost  float64
}

func (this *annealingTable) ToggleRandomInOutValue() (indexToggled uint64, newInOutValue uint64) {
	indexToggled = this.GenerateRandomInOutIndex()
	this.ToggleInOutValueAtIndex(indexToggled)
	newInOutValue = this.rows[indexToggled].InOut
	return
}

func (this *annealingTable) GenerateRandomInOutIndex() uint64 {
	tableSize := len(this.rows)
	return (uint64)(randomNumberGenerator.Intn(tableSize))
}

func (this *annealingTable) ToggleInOutValueAtIndex(index uint64) {
	newInOutValue := (this.rows[index].InOut + 1) % 2
	this.setInOutValueAtIndex(newInOutValue, index)
}

func (this *annealingTable) setInOutValueAtIndex(inOutValue uint64, index uint64) {
	this.rows[index].InOut = inOutValue
	this.rows[index].Reached = (float64)(this.rows[index].InOut) * this.rows[index].Cost
	this.rows[index].InCost = (float64)(this.rows[index].InOut) * this.rows[index].Feature
}

type trackingTable struct {
	rows []trackingData
}

type trackingData struct {
	ObjectiveFunctionValue float64
	Temperature            float64
	ChangeIsDesirable      bool
	AcceptanceProbability  float64
	ChangeAccepted         bool
	InFirst50              uint64
	InSecond50             uint64
	TotalCost              float64
}

func (this *KnapsackObjectiveManager) Initialise() {
	this.BaseObjectiveManager.Initialise()

	this.dataSourcePath = initialiseDataSource()
	this.LogHandler().Info("Opening Excel workbook [" + this.dataSourcePath + "] as data source")

	this.LogHandler().Debug("Retrieving annealing data from workbook")
	this.annealingData = retrieveAnnealingTableFromWorkbook()

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
	this.temperature = temperature  // TODO: Smelly
	this.makeRandomChange(temperature)
	DecideOnWhetherToAcceptChange(this, temperature)
}

func (this *KnapsackObjectiveManager) makeRandomChange(temperature float64) {
	_, newInOutValue := this.annealingData.ToggleRandomInOutValue()

	var changeInObjectiveValue float64
	switch newInOutValue {
	case 0:
		changeInObjectiveValue = -1
	case 1:
		changeInObjectiveValue = 1
	}

	if this.ObjectiveValue()+changeInObjectiveValue >= 0 {
		this.SetChangeInObjectiveValue(changeInObjectiveValue)
	} else {
		this.SetChangeInObjectiveValue(0)
	}
	this.SetObjectiveValue(this.ObjectiveValue() + this.ChangeInObjectiveValue())
}

func (this *KnapsackObjectiveManager) AcceptLastChange() {
	this.BaseObjectiveManager.AcceptLastChange()
	this.addTrackerData()
}

func (this *KnapsackObjectiveManager) addTrackerData() {
	newRow := new(trackingData)
	newRow.ObjectiveFunctionValue = this.ObjectiveValue()
	newRow.Temperature = this.temperature
	newRow.ChangeIsDesirable = this.ChangeIsDesirable()
	newRow.ChangeAccepted = this.ChangeAccepted()
	newRow.AcceptanceProbability = this.AcceptanceProbability()
	this.trackingData.rows = append(this.trackingData.rows, *newRow)
}

func (this *KnapsackObjectiveManager) RevertLastChange() {
	this.SetObjectiveValue(this.ObjectiveValue() - this.ChangeInObjectiveValue())
	this.BaseObjectiveManager.RevertLastChange()
}
