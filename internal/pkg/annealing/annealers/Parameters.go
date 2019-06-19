// Copyright (c) 2018 Australian Rivers Institute.

package annealers

import "github.com/LindsayBradford/crem/internal/pkg/parameters"

type Parameters struct {
	parameters.Parameters
}

func (p *Parameters) Initialise() *Parameters {
	p.Parameters.
		Initialise("Annealer Parameter Validation").
		Enforcing(DefineSpecifications())
	return p
}
