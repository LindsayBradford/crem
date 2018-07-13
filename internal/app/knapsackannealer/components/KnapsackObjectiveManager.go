// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package components

import (
	"math/rand"
	"os"
	"path/filepath"
	"time"

	. "github.com/LindsayBradford/crm/annealing/objectives"
	"github.com/LindsayBradford/crm/excel"
)

type KnapsackObjectiveManager struct {
	BaseObjectiveManager
	excelHandler *excel.ExcelHandler
	dataSource *excel.Workbook
	annealingData *annealingTable
	trackingData *trackingTable
}

var (
 randomNumberGenerator = rand.New(rand.NewSource(time.Now().UnixNano()))

)

type annealingTable struct {
	name string
	rows []annealingData
}

type annealingData struct {
	Cost            float64
	Feature         float64
	X               uint64
	Y               uint64
	InOut           uint64
	Reached         float64
	InCost          float64
}

func (this *annealingTable) Initialise(dataSource *excel.Workbook) {
	dataWorksheet := dataSource.WorksheetNamed("Data")
	this.name = dataWorksheet.Name()

	dataWorksheetRows := dataWorksheet.UsedRange().Rows().Count()
	this.rows = make([]annealingData, dataWorksheetRows - 1)

	for index := 0; index < len(this.rows); index++ {
		this.rows[index].Cost = dataWorksheet.Cells(2+uint(index), 1).Value().(float64)
		this.rows[index].Feature = dataWorksheet.Cells(2+uint(index), 2).Value().(float64)
		this.rows[index].X = (uint64)(dataWorksheet.Cells(2+uint(index), 3).Value().(float64))
		this.rows[index].Y = (uint64)(dataWorksheet.Cells(2+uint(index), 4).Value().(float64))
		this.rows[index].InOut = (uint64)(dataWorksheet.Cells(2+uint(index), 5).Value().(float64))
	}

	this.RandomiseInitialSolutionSet()
}

func (this *annealingTable) RandomiseInitialSolutionSet() {
	for index := 0; index < len(this.rows); index++ {
		randomInOutValue := generateRandomInOutValue()
		this.setInOutValueAtIndex(randomInOutValue, uint64(index))
	}
}

func generateRandomInOutValue() uint64 {
	return (uint64)(randomNumberGenerator.Intn(2))
}

func (this *annealingTable) ToggleRandomInOutValue() (indexToggled uint64, newInOutValue uint64){
		indexToggled = this.GenerateRandomInOutIndex()
		this.ToggleInOutValueAtIndex(indexToggled)
		newInOutValue = this.rows[indexToggled].InOut
		return
}

func (this *annealingTable) GenerateRandomInOutIndex() uint64 {
	tableSize := len(this.rows)
	return (uint64)(1 + randomNumberGenerator.Intn(tableSize))
}

func (this *annealingTable) ToggleInOutValueAtIndex(index uint64) {
	newInOutValue := this.rows[index].InOut + 1 % 2
	this.setInOutValueAtIndex(newInOutValue, index)
}

func (this *annealingTable) setInOutValueAtIndex(inOutValue uint64, index uint64) {
	this.rows[index].InOut = inOutValue
	this.rows[index].Reached = (float64)(this.rows[index].InOut) * this.rows[index].Cost
	this.rows[index].InCost = (float64)(this.rows[index].InOut) * this.rows[index].Feature
}

func (this *annealingTable) TearDown(dataSource *excel.Workbook) {
	dataWorksheet := dataSource.WorksheetNamed("Data")
	for index := 0; index < len(this.rows); index++ {
		dataWorksheet.Cells(2+uint(index), 5).SetValue(this.rows[index].InOut)
	}
}

type trackingTable []struct {
	name string
	rows []trackingData
}

type trackingData struct {
   ObjectiveFunctionValue float64
   Temperature float64
   TestValue   float64
   RandomValue float64
   GoodChoice  uint64
   InFirst50   uint64
   InSecond50 uint64
   TotalCost  float64
}

func (this *trackingTable) Initialise(dataSource *excel.Workbook) {
	trackerWorksheet := dataSource.WorksheetNamed("Tracker")
	trackerWorksheet.UsedRange().Clear()
}

func (this *trackingTable) TearDown(dataSource *excel.Workbook) {
	dataWorksheet := dataSource.WorksheetNamed("Tracker")
	dataWorksheet.Cells(1, 1).SetValue("ObjFuncValue")
	dataWorksheet.Cells(1, 2).SetValue("Temperature")
	dataWorksheet.Cells(1, 3).SetValue("TestValue")
	dataWorksheet.Cells(1, 4).SetValue("RandomValue")
	dataWorksheet.Cells(1, 5).SetValue("GoodChoice")
	dataWorksheet.Cells(1, 6).SetValue("InFirst50")
	dataWorksheet.Cells(1, 7).SetValue("InSecond50")
	dataWorksheet.Cells(1, 8).SetValue("TotalCost")
}

func (this *KnapsackObjectiveManager) Initialise() {
	this.BaseObjectiveManager.Initialise()

	defer func() {
		if r := recover(); r != nil {
			this.excelHandler.Destroy()
			panic("Failed initialising via Excel data-source")
		}
	}()

	this.excelHandler, _ = excel.InitialiseHandler()
	this.InitialiseDataSource()

	this.LogHandler().Debug("Initialising annealing data")
	this.annealingData = new(annealingTable)
	this.annealingData.Initialise(this.dataSource)

	this.LogHandler().Debug("Initialising tracking data")
	this.trackingData = new(trackingTable)
	this.trackingData.Initialise(this.dataSource)
}

func (this *KnapsackObjectiveManager) TearDown() {
	this.BaseObjectiveManager.TearDown()

	this.annealingData.TearDown(this.dataSource)
	this.trackingData.TearDown(this.dataSource)

	this.dataSource.Save(); this.dataSource.Close()
	this.excelHandler.Destroy()
}

func (this *KnapsackObjectiveManager) InitialiseDataSource() {

	workingDirectory, _ := os.Getwd()
	testFixtureAbsolutePath := filepath.Join(workingDirectory, "testdata", "KnapsackAnnealerTestFixture.xls", "fart")

	this.LogHandler().Debug("Opening data-source [" + testFixtureAbsolutePath + "]")

	dataSource, dataSourceErr := this.excelHandler.Workbooks().Open(testFixtureAbsolutePath)

	if dataSourceErr  != nil {
		panic("Datasource [" + testFixtureAbsolutePath + "] could not be opened.")
	}
	this.dataSource = dataSource
}

func (this *KnapsackObjectiveManager) TryRandomChange(temperature float64) {
	this.makeRandomChange()
	DecideOnWhetherToAcceptChange(this, temperature)
}

func (this *KnapsackObjectiveManager) makeRandomChange() {
	_, newInOutValue := this.annealingData.ToggleRandomInOutValue()


	var changeInObjectiveValue float64
	switch newInOutValue {
	case 0:
		changeInObjectiveValue = -1
	case 1:
		changeInObjectiveValue = 1
	}

	if this.ObjectiveValue() + changeInObjectiveValue >= 0 {
		this.SetChangeInObjectiveValue(changeInObjectiveValue)
	} else {
		this.SetChangeInObjectiveValue(0)
	}
	this.SetObjectiveValue(this.ObjectiveValue() + this.ChangeInObjectiveValue())
}

func (this *KnapsackObjectiveManager) AcceptLastChange() {
	this.BaseObjectiveManager.AcceptLastChange()
}

func (this *KnapsackObjectiveManager) RevertLastChange() {
	this.SetObjectiveValue(this.ObjectiveValue() - this.ChangeInObjectiveValue())
	this.BaseObjectiveManager.RevertLastChange()
}
