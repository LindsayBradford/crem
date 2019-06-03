// Copyright (c) 2019 Australian Rivers Institute.

package actions

import (
	"github.com/LindsayBradford/crem/internal/pkg/dataset/tables"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment/parameters"
)

const (
	gullyIdIndex             = 0
	gulliesPlanningUnitIndex = 1
	gullyErosionVolumeIndex  = 2
	gullyChannelLength       = 3
)

type gullySedimentTracker struct {
	GullyId            float64
	SedimentProduction float64
	ChannelLength      float64
}

type GullySedimentContribution struct {
	gulliesTable tables.CsvTable
	parameters   parameters.Parameters

	contributionMap map[string][]gullySedimentTracker
}

func (bsc *GullySedimentContribution) Initialise(gulliesTable tables.CsvTable, parameters parameters.Parameters) {
	bsc.gulliesTable = gulliesTable
	bsc.parameters = parameters
	bsc.populateContributionMap()
}

func (bsc *GullySedimentContribution) populateContributionMap() {
	_, rowCount := bsc.gulliesTable.ColumnAndRowSize()
	bsc.contributionMap = make(map[string][]gullySedimentTracker, 0)

	for row := uint(0); row < rowCount; row++ {
		bsc.populateContributionMapEntry(row)
	}
}

func (bsc *GullySedimentContribution) populateContributionMapEntry(rowNumber uint) {
	planningUnit := bsc.gulliesTable.CellFloat64(gulliesPlanningUnitIndex, rowNumber)
	mapKey := Float64ToPlanningUnitId(planningUnit)

	newGullyTracker := gullySedimentTracker{
		GullyId:            bsc.gulliesTable.CellFloat64(gullyIdIndex, rowNumber),
		SedimentProduction: bsc.gullySediment(rowNumber),
		ChannelLength:      bsc.channelLength(rowNumber),
	}

	if trackers, hasTrackers := bsc.contributionMap[mapKey]; hasTrackers {
		trackers = append(trackers, newGullyTracker)
		bsc.contributionMap[mapKey] = trackers
	} else {
		bsc.contributionMap[mapKey] = []gullySedimentTracker{newGullyTracker}
	}
}

func (bsc *GullySedimentContribution) gullySediment(rowNumber uint) float64 {
	gullyVolume := bsc.gulliesTable.CellFloat64(gullyErosionVolumeIndex, rowNumber)
	return bsc.SedimentFromVolume(gullyVolume)
}

func (bsc *GullySedimentContribution) SedimentFromVolume(gullyVolume float64) float64 {
	if gullyVolume == 0 {
		return 0
	}

	drySedimentDensity := bsc.parameters.GetFloat64(parameters.SedimentDensity)
	nonLinearErosionRateCompensationFactor := bsc.parameters.GetFloat64(parameters.GullyCompensationFactor)

	totalErosionTime := float64(bsc.parameters.GetInt64(parameters.YearsOfErosion))
	suspendedSedimentProportion := bsc.parameters.GetFloat64(parameters.SuspendedSedimentProportion)

	return ((gullyVolume * drySedimentDensity * nonLinearErosionRateCompensationFactor) / totalErosionTime) *
		suspendedSedimentProportion
}

func (bsc *GullySedimentContribution) channelLength(rowNumber uint) float64 {
	return bsc.gulliesTable.CellFloat64(gullyChannelLength, rowNumber)
}

func (bsc *GullySedimentContribution) OriginalSedimentContribution() float64 {
	sedimentContribution := float64(0)
	for _, trackers := range bsc.contributionMap {
		for _, tracker := range trackers {
			sedimentContribution += tracker.SedimentProduction
		}
	}
	return sedimentContribution
}

func (bsc *GullySedimentContribution) SedimentContribution(planningUnit string) float64 {
	sedimentContribution := float64(0)

	if trackers, hasTrackers := bsc.contributionMap[planningUnit]; hasTrackers {
		for _, tracker := range trackers {
			sedimentContribution += tracker.SedimentProduction
		}
	}

	return sedimentContribution
}

func (bsc *GullySedimentContribution) ChannelLength(planningUnit string) float64 {
	channelLength := float64(0)

	if trackers, hasTrackers := bsc.contributionMap[planningUnit]; hasTrackers {
		for _, tracker := range trackers {
			channelLength += tracker.ChannelLength
		}
	}

	return channelLength
}
