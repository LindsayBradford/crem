// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"fmt"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	assert "github.com/LindsayBradford/crem/pkg/assert/debug"
)

type ActionSource string

const (
	ActionedTotalCarbon action.ModelVariableName = "ActionedTotalCarbon"
	OriginalTotalCarbon action.ModelVariableName = "OriginalTotalCarbon"
	TotalNitrogen       action.ModelVariableName = "TotalNitrogen"
	OpportunityCost     action.ModelVariableName = "OpportunityCost"

	parentSoilsPlanningUnitIndex = 0
	sourceFilterIndex            = 1
	nitrogenValueIndex           = 2
	carbonValueIndex             = 3
	deltaCarbonValueIndex        = 4
	opportunityCostIndex         = 5

	RiparianSource  ActionSource = "Riparian"
	HillSlopeSource ActionSource = "Hillslope"
	GullySource     ActionSource = "Gully"
	UndefinedSource ActionSource = ""

	NitrogenAttribute        = "Nitrogen"
	CarbonAttribute          = "Carbon"
	CarbonDeltaAttribute     = "CarbonDelta"
	OpportunityCostAttribute = "OpportunityCot"
)

type Container struct {
	sourceFilter ActionSource
	actionsMap   map[string]float64
}

func (c *Container) WithSourceFilter(sourceFilter ActionSource) *Container {
	c.sourceFilter = sourceFilter
	return c
}

func (c *Container) WithActionsTable(actionsTable tables.CsvTable) *Container {
	_, rowCount := actionsTable.ColumnAndRowSize()
	c.actionsMap = make(map[string]float64, 0)

	for rowNumber := uint(0); rowNumber < rowCount; rowNumber++ {

		sourceType := ActionSource(actionsTable.CellString(sourceFilterIndex, rowNumber))
		if c.sourceFilter != UndefinedSource && sourceType != c.sourceFilter {
			continue
		}

		var mapKey string
		planningUnit := planningunit.Id(actionsTable.CellFloat64(parentSoilsPlanningUnitIndex, rowNumber))

		nitrogenValue := actionsTable.CellFloat64(nitrogenValueIndex, rowNumber)
		mapKey = c.DeriveMapKey(planningUnit, sourceType, NitrogenAttribute)
		c.actionsMap[mapKey] = nitrogenValue

		carbonValue := actionsTable.CellFloat64(carbonValueIndex, rowNumber)
		mapKey = c.DeriveMapKey(planningUnit, sourceType, CarbonAttribute)
		c.actionsMap[mapKey] = carbonValue

		carbonDeltaValue := actionsTable.CellFloat64(deltaCarbonValueIndex, rowNumber)
		mapKey = c.DeriveMapKey(planningUnit, sourceType, CarbonDeltaAttribute)
		c.actionsMap[mapKey] = carbonDeltaValue

		opportunityCostValue := actionsTable.CellFloat64(opportunityCostIndex, rowNumber)
		mapKey = c.DeriveMapKey(planningUnit, sourceType, OpportunityCostAttribute)
		c.actionsMap[mapKey] = opportunityCostValue
	}
	return c
}

func (c *Container) MapValue(key string) float64 {
	mappedValue := c.actionsMap[key]
	failureMsg := fmt.Sprintf("Container doesn't have value mapped to key [%s]", key)
	assert.That(mappedValue > 0).WithFailureMessage(failureMsg).Holds()
	return mappedValue
}

func (c *Container) DeriveMapKey(planningUnit planningunit.Id, sourceType ActionSource, elementType string) string {
	if c.sourceFilter == UndefinedSource {
		return fmt.Sprintf("%d,%s,%s", planningUnit, sourceType, elementType)
	}
	return fmt.Sprintf("%d,%s", planningUnit, elementType)
}

func (c *Container) nitrogenAttributeValue(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.sourceFilter, NitrogenAttribute)
	return c.actionsMap[key]
}

func (c *Container) carbonAttributeValue(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.sourceFilter, CarbonAttribute)
	return c.actionsMap[key]
}

func (c *Container) deltaCarbonAttributeValue(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.sourceFilter, CarbonDeltaAttribute)
	return c.actionsMap[key]
}

func (c *Container) opportunityCostAttributeValue(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.sourceFilter, OpportunityCostAttribute)
	return c.actionsMap[key]
}
