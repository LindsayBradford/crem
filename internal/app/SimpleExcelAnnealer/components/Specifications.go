// Copyright (c) 2019 Australian Rivers Institute.

package components

import (
	. "github.com/LindsayBradford/crem/internal/pkg/parameters/specification"
)

const (
	Penalty        string = "Penalty"
	DataSourcePath string = "DataSourcePath"
)

func DefineSpecifications() *Specifications {
	specs := NewSpecifications()
	specs.Add(
		Specification{
			Key:          Penalty,
			Validator:    IsDecimal,
			DefaultValue: 1.0,
		},
	).Add(
		Specification{
			Key:          DataSourcePath,
			Validator:    IsReadableFile,
			DefaultValue: "",
		},
	)
	return specs
}
