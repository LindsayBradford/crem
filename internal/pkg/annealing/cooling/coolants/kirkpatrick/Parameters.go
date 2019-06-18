// Copyright (c) 2018 Australian Rivers Institute.

package kirkpatrick

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
)

const (
	_                          = iota
	CoolingFactor       string = "CoolingFactor"
	StartingTemperature string = "StartingTemperature"
)

type Parameters struct {
	parameters.Parameters
}

func (kp *Parameters) Initialise() *Parameters {
	kp.Parameters.CreateEmpty().
		WithSpecifications(
			DefineSpecifications(),
		).AssigningDefaults()
	return kp
}
