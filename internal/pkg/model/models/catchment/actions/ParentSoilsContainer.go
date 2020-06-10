// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"fmt"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	assert "github.com/LindsayBradford/crem/pkg/assert/debug"
)

type ParentSoilSource string

const (
	ActionedTotalCarbon action.ModelVariableName = "ActionedTotalCarbon"
	OriginalTotalCarbon action.ModelVariableName = "OriginalTotalCarbon"
	TotalNitrogen       action.ModelVariableName = "TotalNitrogen"

	parentSoilsPlanningUnitIndex = 0
	sourceFilterIndex            = 1
	nitrogenValueIndex           = 2
	carbonValueIndex             = 3
	deltaCarbonValueIndex        = 4

	RiparianSource  ParentSoilSource = "Riparian"
	HillSlopeSource ParentSoilSource = "Hillslope"
	GullySource     ParentSoilSource = "Gully"
	UndefinedSource ParentSoilSource = ""

	NitrogenAttribute    = "Nitrogen"
	CarbonAttribute      = "Carbon"
	CarbonDeltaAttribute = "CarbonDelta"
)

type ParentSoilsContainer struct {
	sourceFilter   ParentSoilSource
	parentSoilsMap map[string]float64
}

func (c *ParentSoilsContainer) WithSourceFilter(sourceFilter ParentSoilSource) *ParentSoilsContainer {
	c.sourceFilter = sourceFilter
	return c
}

func (c *ParentSoilsContainer) WithParentSoilsTable(parentSoilsTable tables.CsvTable) *ParentSoilsContainer {
	_, rowCount := parentSoilsTable.ColumnAndRowSize()
	c.parentSoilsMap = make(map[string]float64, 0)

	for rowNumber := uint(0); rowNumber < rowCount; rowNumber++ {

		sourceType := ParentSoilSource(parentSoilsTable.CellString(sourceFilterIndex, rowNumber))
		if c.sourceFilter != UndefinedSource && sourceType != c.sourceFilter {
			continue
		}

		var mapKey string
		planningUnit := planningunit.Id(parentSoilsTable.CellFloat64(parentSoilsPlanningUnitIndex, rowNumber))

		nitrogenValue := parentSoilsTable.CellFloat64(nitrogenValueIndex, rowNumber)
		mapKey = c.DeriveMapKey(planningUnit, sourceType, NitrogenAttribute)
		c.parentSoilsMap[mapKey] = nitrogenValue

		carbonValue := parentSoilsTable.CellFloat64(carbonValueIndex, rowNumber)
		mapKey = c.DeriveMapKey(planningUnit, sourceType, CarbonAttribute)
		c.parentSoilsMap[mapKey] = carbonValue

		carbonDeltaValue := parentSoilsTable.CellFloat64(deltaCarbonValueIndex, rowNumber)
		mapKey = c.DeriveMapKey(planningUnit, sourceType, CarbonDeltaAttribute)
		c.parentSoilsMap[mapKey] = carbonDeltaValue
	}
	return c
}

func (c *ParentSoilsContainer) MapValue(key string) float64 {
	mappedValue := c.parentSoilsMap[key]
	failureMsg := fmt.Sprintf("ParentSoilsContainer doesn't have value mapped to key [%s]", key)
	assert.That(mappedValue > 0).WithFailureMessage(failureMsg).Holds()
	return mappedValue
}

func (c *ParentSoilsContainer) DeriveMapKey(planningUnit planningunit.Id, sourceType ParentSoilSource, elementType string) string {
	if c.sourceFilter == UndefinedSource {
		return fmt.Sprintf("%d,%s,%s", planningUnit, sourceType, elementType)
	}
	return fmt.Sprintf("%d,%s", planningUnit, elementType)
}

func (c *ParentSoilsContainer) nitrogenAttributeValue(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.sourceFilter, NitrogenAttribute)
	return c.parentSoilsMap[key]
}

func (c *ParentSoilsContainer) carbonAttributeValue(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.sourceFilter, CarbonAttribute)
	return c.parentSoilsMap[key]
}

func (c *ParentSoilsContainer) deltaCarbonAttributeValue(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.sourceFilter, CarbonDeltaAttribute)
	return c.parentSoilsMap[key]
}
