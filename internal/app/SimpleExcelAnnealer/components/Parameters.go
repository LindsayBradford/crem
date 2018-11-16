// Copyright (c) 2018 Australian Rivers Institute.

package components

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters"
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
			Key:          Penalty,
			Validator:    p.IsDecimal,
			DefaultValue: 1.0,
		},
	)

	p.AddMetaData(
		parameters.MetaData{
			Key:          DataSourcePath,
			Validator:    p.IsReadableFile,
			DefaultValue: "",
		},
	)
}
