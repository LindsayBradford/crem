// Copyright (c) 2018 Australian Rivers Institute.

package parameters

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
)

type Parameters struct {
	parameters.Parameters
}

func (p *Parameters) Initialise() *Parameters {
	p.Parameters.CreateEmpty().
		WithSpecifications(
			DefineSpecifications(),
		).AssigningDefaults()
	return p
}
