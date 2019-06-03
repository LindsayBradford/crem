// Copyright (c) 2018 Australian Rivers Institute.

package dumb

import "github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"

const (
	InitialObjectiveValue = "InitialObjectiveValue"
	MinimumObjectiveValue = "MinimumObjectiveValue"
	MaximumObjectiveValue = "MaximumObjectiveValue"
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
			Key:          InitialObjectiveValue,
			Validator:    kp.Parameters.IsDecimal,
			DefaultValue: float64(1000),
		},
	)
	kp.AddMetaData(
		parameters.MetaData{
			Key:          MinimumObjectiveValue,
			Validator:    kp.Parameters.IsDecimal,
			DefaultValue: float64(0),
		},
	)
	kp.AddMetaData(
		parameters.MetaData{
			Key:          MaximumObjectiveValue,
			Validator:    kp.Parameters.IsDecimal,
			DefaultValue: float64(2000),
		},
	)
}
