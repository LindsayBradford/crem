// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"fmt"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	assert "github.com/LindsayBradford/crem/pkg/assert/debug"
)

type ActionSource string

const (
	subCatchmentIndex                = 0
	filterIndex                      = 1
	opportunityCostIndex             = 2
	implementationCostIndex          = 3
	particulateNitrogenOriginalIndex = 4
	particulateNitrogenActionedIndex = 5
	hillSlopeErosionOriginalIndex    = 6
	hillSlopeErosionActionedIndex    = 7
	fineSedimentOriginalIndex        = 8
	fineSedimentActionedIndex        = 9

	RiparianSource  ActionSource = "Riparian"
	HillSlopeSource ActionSource = "Hillslope"
	GullySource     ActionSource = "Gully"
	UndefinedSource ActionSource = ""

	OpportunityCostAttribute             = "OpportunityCot"
	ImplementationCostAttribute          = "ImplementationCost"
	ParticulateNitrogenOriginalAttribute = "ParticulateNitrogenOriginal"
	ParticulateNitrogenActionedAttribute = "ParticulateNitrogenActioned"
	HillSlopeErosionOriginalAttribute    = "HillSlopeErosionOriginal"
	HillSlopeErosionActionedAttribute    = "HillSlopeErosionActioned"
	FineSedimentOriginalAttribute        = "FineSedimentOriginal"
	FineSedimentActionedAttribute        = "FineSedimentActioned"
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

		sourceType := ActionSource(actionsTable.CellString(filterIndex, rowNumber))
		if c.sourceFilter != UndefinedSource && sourceType != c.sourceFilter {
			continue
		}

		subCatchment := planningunit.Id(actionsTable.CellFloat64(subCatchmentIndex, rowNumber))

		var mapAttribute = func(index uint, attribute string) {
			value := actionsTable.CellFloat64(index, rowNumber)
			mapKey := c.DeriveMapKey(subCatchment, sourceType, attribute)
			c.actionsMap[mapKey] = value
		}

		mapAttribute(opportunityCostIndex, OpportunityCostAttribute)
		mapAttribute(implementationCostIndex, ImplementationCostAttribute)

		mapAttribute(particulateNitrogenOriginalIndex, ParticulateNitrogenOriginalAttribute)
		mapAttribute(particulateNitrogenActionedIndex, ParticulateNitrogenActionedAttribute)

		mapAttribute(hillSlopeErosionOriginalIndex, HillSlopeErosionOriginalAttribute)
		mapAttribute(hillSlopeErosionActionedIndex, HillSlopeErosionActionedAttribute)

		mapAttribute(fineSedimentOriginalIndex, FineSedimentOriginalAttribute)
		mapAttribute(fineSedimentActionedIndex, FineSedimentActionedAttribute)
	}
	return c
}

func (c *Container) MapValue(key string) float64 {
	mappedValue := c.actionsMap[key]
	failureMsg := fmt.Sprintf("Container doesn't have value mapped to key [%s]", key)
	assert.That(mappedValue > 0).WithFailureMessage(failureMsg).Holds()
	return mappedValue
}

func (c *Container) DeriveMapKey(subCatchment planningunit.Id, sourceType ActionSource, elementType string) string {
	if c.sourceFilter == UndefinedSource {
		return fmt.Sprintf("%d,%s,%s", subCatchment, sourceType, elementType)
	}
	return fmt.Sprintf("%d,%s", subCatchment, elementType)
}

func (c *Container) opportunityCost(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.sourceFilter, OpportunityCostAttribute)
	return c.actionsMap[key]
}

func (c *Container) implementationCost(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.sourceFilter, ImplementationCostAttribute)
	return c.actionsMap[key]
}

func (c *Container) originalParticulateNitrogen(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.sourceFilter, ParticulateNitrogenOriginalAttribute)
	return c.actionsMap[key]
}

func (c *Container) actionedParticulateNitrogen(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.sourceFilter, ParticulateNitrogenActionedAttribute)
	return c.actionsMap[key]
}

func (c *Container) originalHillSlopeErosion(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.sourceFilter, HillSlopeErosionOriginalAttribute)
	return c.actionsMap[key]
}

func (c *Container) actionedHillSlopeErosion(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.sourceFilter, HillSlopeErosionOriginalAttribute)
	return c.actionsMap[key]
}

func (c *Container) originalFineSediment(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.sourceFilter, FineSedimentOriginalAttribute)
	return c.actionsMap[key]
}

func (c *Container) actionedFineSediment(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.sourceFilter, FineSedimentActionedAttribute)
	return c.actionsMap[key]
}
