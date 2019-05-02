// Copyright (c) 2018 Australian Rivers Institute.

package annealers

import "github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"

const (
	_                        = iota
	MaximumIterations string = "MaximumIterations"
)

type Parameters struct {
	parameters.Parameters
}

func (p *Parameters) Initialise() *Parameters {
	p.Parameters.Initialise()
	p.buildMetaData()
	p.CreateDefaults()
	return p
}

func (p *Parameters) buildMetaData() {
	p.AddMetaData(
		parameters.MetaData{
			Key:          MaximumIterations,
			Validator:    p.IsNonNegativeInteger,
			DefaultValue: int64(0),
		},
	)
}
