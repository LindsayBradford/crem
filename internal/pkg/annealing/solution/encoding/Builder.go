// Copyright (c) 2019 Australian Rivers Institute.

package encoding

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/encoding/csv"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/encoding/excel"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/encoding/json"
)

type OutputType string

const (
	UndefinedOutput = ""
	CsvOutput       = "CSV"
	JsonOutput      = "JSON"
	ExcelOutput     = "EXCEL"
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
	case UndefinedOutput, CsvOutput:
		return new(csv.Encoder).WithOutputPath(b.outputPath)
	case JsonOutput:
		return new(json.Encoder).WithOutputPath(b.outputPath)
	case ExcelOutput:
		return new(excel.Encoder).WithOutputPath(b.outputPath)
	default:
		return NullEncoder
	}
}
