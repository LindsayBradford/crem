// Copyright (c) 2019 Australian Rivers Institute.

package components

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/parameters/specification"
)

const (
	Penalty        string = "Penalty"
	DataSourcePath string = "DataSourcePath"
)

func DefineSpecifications() *specification.Specifications {
	specs := specification.New()
	specs.Add(
		specification.Specification{
			Key:          Penalty,
			Validator:    specification.IsDecimal,
			DefaultValue: 1.0,
		},
	).Add(
		specification.Specification{
			Key:          DataSourcePath,
			Validator:    specification.IsReadableFile,
			DefaultValue: "",
		},
	)
	return specs
}
