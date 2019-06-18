// Copyright (c) 2018 Australian Rivers Institute.

package dumb

import "github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"

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
