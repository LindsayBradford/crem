// Copyright (c) 2019 Australian Rivers Institute.

package parameters

import "github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"

const (
	InitialObjectiveOneValue   = "InitialObjectiveOneValue"
	InitialObjectiveTwoValue   = "InitialObjectiveTwoValue"
	InitialObjectiveThreeValue = "InitialObjectiveThreeValue"

	NumberOfPlanningUnits = "NumberOfPlanningUnits"
)

type Parameters struct {
	parameters.Parameters
}

func (kp *Parameters) Initialise() *Parameters {
	kp.Parameters.Initialise()
	kp.buildMetaData()
	kp.CreateDefaults()
	return kp
}

func (kp *Parameters) buildMetaData() {
	kp.AddMetaData(
		parameters.MetaData{
			Key:          InitialObjectiveOneValue,
			Validator:    kp.Parameters.IsDecimal,
			DefaultValue: float64(1000),
		},
	)
	kp.AddMetaData(
		parameters.MetaData{
			Key:          InitialObjectiveTwoValue,
			Validator:    kp.Parameters.IsDecimal,
			DefaultValue: float64(1000),
		},
	)
	kp.AddMetaData(
		parameters.MetaData{
			Key:          InitialObjectiveThreeValue,
			Validator:    kp.Parameters.IsDecimal,
			DefaultValue: float64(1000),
		},
	)
	kp.AddMetaData(
		parameters.MetaData{
			Key:          NumberOfPlanningUnits,
			Validator:    kp.Parameters.IsNonNegativeInteger,
			DefaultValue: int64(100),
		},
	)
}
