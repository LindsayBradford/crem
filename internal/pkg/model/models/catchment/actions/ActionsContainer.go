// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"fmt"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/planningunit"
	assert "github.com/LindsayBradford/crem/pkg/assert/debug"
	"strconv"
	"strings"
)

type ActionType string

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

	dissolvedNitrogenOriginalIndex            = 10
	dissolvedNitrogenActionedIndex            = 11
	dissolvedNitrogenRemovalEfficiencyIndex   = 12
	particulateNitrogenRemovalEfficiencyIndex = 13
	sedimentRemovalEfficiencyIndex            = 14

	RiparianType  ActionType = "Riparian"
	HillSlopeType ActionType = "Hillslope"
	GullyType     ActionType = "Gully"
	WetlandType   ActionType = "Wetland"

	UndefinedType ActionType = ""

	OpportunityCostAttribute             = "OpportunityCot"
	ImplementationCostAttribute          = "ImplementationCost"
	ParticulateNitrogenOriginalAttribute = "ParticulateNitrogenOriginal"
	ParticulateNitrogenActionedAttribute = "ParticulateNitrogenActioned"
	HillSlopeErosionOriginalAttribute    = "HillSlopeErosionOriginal"
	HillSlopeErosionActionedAttribute    = "HillSlopeErosionActioned"
	FineSedimentOriginalAttribute        = "FineSedimentOriginal"
	FineSedimentActionedAttribute        = "FineSedimentActioned"

	DissolvedNitrogenOriginalAttribute   = "DissolvedNitrogenOriginal"
	DissolvedNitrogenActionedAttribute   = "DissolvedNitrogenActioned"
	DissolvedNitrogenRemovalEfficiency   = "DissolvedNitrogenRemovalEfficiency"
	ParticulateNitrogenRemovalEfficiency = "ParticulateNitrogenRemovalEfficiency"
	SedimentRemovalEfficiency            = "SedimentRemovalEfficiency"
)

type Container struct {
	filter     ActionType
	actionsMap map[string]float64
}

func (c *Container) WithFilter(filter ActionType) *Container {
	c.filter = filter
	return c
}

func (c *Container) WithActionsTable(actionsTable tables.CsvTable) *Container {
	_, rowCount := actionsTable.ColumnAndRowSize()
	c.actionsMap = make(map[string]float64, 0)

	for rowNumber := uint(0); rowNumber < rowCount; rowNumber++ {

		sourceType := ActionType(actionsTable.CellString(filterIndex, rowNumber))
		if c.filter != UndefinedType && sourceType != c.filter {
			continue
		}

		subCatchment := planningunit.Id(actionsTable.CellFloat64(subCatchmentIndex, rowNumber))

		mapAttribute := func(index uint, attribute string) {
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

		mapAttribute(dissolvedNitrogenOriginalIndex, DissolvedNitrogenOriginalAttribute)
		mapAttribute(dissolvedNitrogenActionedIndex, DissolvedNitrogenActionedAttribute)

		mapAttribute(dissolvedNitrogenRemovalEfficiencyIndex, DissolvedNitrogenRemovalEfficiency)
		mapAttribute(particulateNitrogenRemovalEfficiencyIndex, ParticulateNitrogenRemovalEfficiency)
		mapAttribute(sedimentRemovalEfficiencyIndex, SedimentRemovalEfficiency)
	}
	return c
}

func (c *Container) MapValue(key string) float64 {
	mappedValue := c.actionsMap[key]
	failureMsg := fmt.Sprintf("Container doesn't have value mapped to key [%s]", key)
	assert.That(mappedValue > 0).WithFailureMessage(failureMsg).Holds()
	return mappedValue
}

type KeyComponents struct {
	SubCatchment planningunit.Id
	Action       ActionType
	ElementType  string
}

func (c *Container) DeriveMapKey(subCatchment planningunit.Id, actionType ActionType, elementType string) string {
	if c.filter == UndefinedType {
		return fmt.Sprintf("%d,%s,%s", subCatchment, actionType, elementType)
	}
	return fmt.Sprintf("%d,%s", subCatchment, elementType)
}

func (c *Container) DeriveMapKeyComponents(mapKey string) *KeyComponents {
	rawComponents := strings.Split(mapKey, ",")

	rawPlanningUnitId, rawPlanningUnitIdError := strconv.ParseInt(rawComponents[0], 10, 64)
	if rawPlanningUnitIdError != nil {
		return nil
	}

	if c.filter == UndefinedType {
		actionType := ActionType(rawComponents[1])
		rawElementType := rawComponents[2]
		return &KeyComponents{
			SubCatchment: planningunit.Id(rawPlanningUnitId),
			Action:       actionType,
			ElementType:  rawElementType,
		}
	} else {
		rawElementType := rawComponents[1]
		return &KeyComponents{
			SubCatchment: planningunit.Id(rawPlanningUnitId),
			ElementType:  rawElementType,
		}
	}
}

func (c *Container) opportunityCost(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.filter, OpportunityCostAttribute)
	return c.actionsMap[key]
}

func (c *Container) implementationCost(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.filter, ImplementationCostAttribute)
	return c.actionsMap[key]
}

func (c *Container) originalParticulateNitrogen(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.filter, ParticulateNitrogenOriginalAttribute)
	return c.actionsMap[key]
}

func (c *Container) actionedParticulateNitrogen(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.filter, ParticulateNitrogenActionedAttribute)
	return c.actionsMap[key]
}

func (c *Container) originalHillSlopeErosion(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.filter, HillSlopeErosionOriginalAttribute)
	return c.actionsMap[key]
}

func (c *Container) actionedHillSlopeErosion(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.filter, HillSlopeErosionActionedAttribute)
	return c.actionsMap[key]
}

func (c *Container) originalFineSediment(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.filter, FineSedimentOriginalAttribute)
	return c.actionsMap[key]
}

func (c *Container) actionedFineSediment(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.filter, FineSedimentActionedAttribute)
	return c.actionsMap[key]
}

func (c *Container) originalDissolvedNitrogen(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.filter, DissolvedNitrogenOriginalAttribute)
	return c.actionsMap[key]
}

func (c *Container) actionedDissolvedNitrogen(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.filter, DissolvedNitrogenActionedAttribute)
	return c.actionsMap[key]
}

func (c *Container) dissolvedNitrogenRemovalEfficiency(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.filter, DissolvedNitrogenRemovalEfficiency)
	return c.actionsMap[key]
}

func (c *Container) particulateNitrogenRemovalEfficiency(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.filter, ParticulateNitrogenRemovalEfficiency)
	return c.actionsMap[key]
}

func (c *Container) sedimentNitrogenRemovalEfficiency(planningUnit planningunit.Id) float64 {
	key := c.DeriveMapKey(planningUnit, c.filter, SedimentRemovalEfficiency)
	return c.actionsMap[key]
}

func (c *Container) Map() map[string]float64 {
	return c.actionsMap
}
