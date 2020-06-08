// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"fmt"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
)

const (
	ActionedTotalCarbon action.ModelVariableName = "ActionedTotalCarbon"
	OriginalTotalCarbon action.ModelVariableName = "OriginalTotalCarbon"
	TotalNitrogen       action.ModelVariableName = "TotalNitrogen"

	parentSoilsPlanningUnitIndex = 0
	sourceFilterIndex            = 1
	nitrogenValueIndex           = 2
	carbonValueIndex             = 3
	deltaCarbonValueIndex        = 4

	nitrogenAttribute    = "Nitrogen"
	carbonAttribute      = "Carbon"
	carbonDeltaAttribute = "CarbonDelta"
)

type parentSoilsContainer struct {
	parentSoilsMap map[string]float64
}

func nitrogenKey(planningUnit planningunit.Id) string {
	return parentSoilsMapKey(planningUnit, nitrogenAttribute)
}

func carbonKey(planningUnit planningunit.Id) string {
	return parentSoilsMapKey(planningUnit, carbonAttribute)
}

func deltaCarbonKey(planningUnit planningunit.Id) string {
	return parentSoilsMapKey(planningUnit, carbonDeltaAttribute)
}

func parentSoilsMapKey(planningUnit planningunit.Id, elementType string) string {
	return fmt.Sprintf("%d,%s", planningUnit, elementType)
}

func (c *parentSoilsContainer) WithParentSoilsTable(parentSoilsTable tables.CsvTable, sourceFilter string) *parentSoilsContainer {
	_, rowCount := parentSoilsTable.ColumnAndRowSize()
	c.parentSoilsMap = make(map[string]float64, 0)

	for rowNumber := uint(0); rowNumber < rowCount; rowNumber++ {

		sourceType := parentSoilsTable.CellString(sourceFilterIndex, rowNumber)
		if sourceType != sourceFilter {
			continue
		}

		planningUnit := planningunit.Id(parentSoilsTable.CellFloat64(parentSoilsPlanningUnitIndex, rowNumber))

		nitrogenValue := parentSoilsTable.CellFloat64(nitrogenValueIndex, rowNumber)
		c.parentSoilsMap[nitrogenKey(planningUnit)] = nitrogenValue

		carbonValue := parentSoilsTable.CellFloat64(carbonValueIndex, rowNumber)
		c.parentSoilsMap[carbonKey(planningUnit)] = carbonValue

		carbonDeltaValue := parentSoilsTable.CellFloat64(deltaCarbonValueIndex, rowNumber)
		c.parentSoilsMap[deltaCarbonKey(planningUnit)] = carbonDeltaValue
	}
	return c
}

func (c *parentSoilsContainer) nitrogenAttributeValue(planningUnit planningunit.Id) float64 {
	key := nitrogenKey(planningUnit)
	return c.parentSoilsMap[key]
}

func (c *parentSoilsContainer) carbonAttributeValue(planningUnit planningunit.Id) float64 {
	key := carbonKey(planningUnit)
	return c.parentSoilsMap[key]
}

func (c *parentSoilsContainer) deltaCarbonAttributeValue(planningUnit planningunit.Id) float64 {
	key := deltaCarbonKey(planningUnit)
	return c.parentSoilsMap[key]
}
