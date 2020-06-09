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

	undefined = ""
)

type ParentSoilsContainer struct {
	sourceFilter   string
	parentSoilsMap map[string]float64
}

func (c *ParentSoilsContainer) WithSourceFilter(sourceFilter string) *ParentSoilsContainer {
	c.sourceFilter = sourceFilter
	return c
}

func (c *ParentSoilsContainer) WithParentSoilsTable(parentSoilsTable tables.CsvTable) *ParentSoilsContainer {
	_, rowCount := parentSoilsTable.ColumnAndRowSize()
	c.parentSoilsMap = make(map[string]float64, 0)

	for rowNumber := uint(0); rowNumber < rowCount; rowNumber++ {

		sourceType := parentSoilsTable.CellString(sourceFilterIndex, rowNumber)
		if c.sourceFilter != undefined && sourceType != c.sourceFilter {
			continue
		}

		var mapKey string
		planningUnit := planningunit.Id(parentSoilsTable.CellFloat64(parentSoilsPlanningUnitIndex, rowNumber))

		nitrogenValue := parentSoilsTable.CellFloat64(nitrogenValueIndex, rowNumber)
		mapKey = c.deriveMapKey(planningUnit, sourceType, nitrogenAttribute)
		c.parentSoilsMap[mapKey] = nitrogenValue

		carbonValue := parentSoilsTable.CellFloat64(carbonValueIndex, rowNumber)
		mapKey = c.deriveMapKey(planningUnit, sourceType, carbonAttribute)
		c.parentSoilsMap[mapKey] = carbonValue

		carbonDeltaValue := parentSoilsTable.CellFloat64(deltaCarbonValueIndex, rowNumber)
		mapKey = c.deriveMapKey(planningUnit, sourceType, carbonDeltaAttribute)
		c.parentSoilsMap[mapKey] = carbonDeltaValue
	}
	return c
}

func (c *ParentSoilsContainer) deriveMapKey(planningUnit planningunit.Id, sourceType string, elementType string) string {
	if c.sourceFilter == undefined {
		return fmt.Sprintf("%d,%s,%s", planningUnit, sourceType, elementType)
	}
	return fmt.Sprintf("%d,%s", planningUnit, elementType)
}

func (c *ParentSoilsContainer) nitrogenAttributeValue(planningUnit planningunit.Id) float64 {
	key := c.deriveMapKey(planningUnit, c.sourceFilter, nitrogenAttribute)
	return c.parentSoilsMap[key]
}

func (c *ParentSoilsContainer) carbonAttributeValue(planningUnit planningunit.Id) float64 {
	key := c.deriveMapKey(planningUnit, c.sourceFilter, carbonAttribute)
	return c.parentSoilsMap[key]
}

func (c *ParentSoilsContainer) deltaCarbonAttributeValue(planningUnit planningunit.Id) float64 {
	key := c.deriveMapKey(planningUnit, c.sourceFilter, carbonDeltaAttribute)
	return c.parentSoilsMap[key]
}
