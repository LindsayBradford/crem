// Copyright (c) 2019 Australian Rivers Institute.

package encoding

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/encoding/csv"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/encoding/json"
)

type OutputType string

const (
	NoOutput   = ""
	CsvOutput  = "CSV"
	JsonOutput = "JSON"
	// ExcelType
)

type Builder struct {
	outputType OutputType
	outputPath string
}

func (b *Builder) ForOutputType(encoderType OutputType) *Builder {
	b.outputType = encoderType
	return b
}

func (b *Builder) WithOutputPath(outputPath string) *Builder {
	b.outputPath = outputPath
	return b
}

func (b *Builder) Build() Encoder {
	switch b.outputType {
	case NoOutput:
		return NullEncoder
	case CsvOutput:
		return new(csv.CsvEncoder).WithOutputPath(b.outputPath)
	case JsonOutput:
		return new(json.JsonEncoder).WithOutputPath(b.outputPath)
	default:
		return NullEncoder
	}
}
