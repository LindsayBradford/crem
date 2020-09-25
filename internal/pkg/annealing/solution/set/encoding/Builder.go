// Copyright (c) 2019 Australian Rivers Institute.

package encoding

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/encoding"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/set/encoding/csv"
	//"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/set/encoding/excel"
	//"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/set/encoding/json"
)

type Builder struct {
	outputType encoding.OutputType
	outputPath string
}

func (b *Builder) ForOutputType(encoderType encoding.OutputType) *Builder {
	b.outputType = encoderType
	return b
}

func (b *Builder) WithOutputPath(outputPath string) *Builder {
	b.outputPath = outputPath
	return b
}

func (b *Builder) Build() Encoder {
	switch b.outputType {
	// TODO: Support for Json and Excel needed later.
	case encoding.UndefinedOutput, encoding.CsvOutput:
		return new(csv.Encoder).WithOutputPath(b.outputPath)
	case encoding.JsonOutput:
		return new(csv.Encoder).WithOutputPath(b.outputPath)
		//return new(json.Encoder).WithOutputPath(b.outputPath)
	case encoding.ExcelOutput:
		return new(csv.Encoder).WithOutputPath(b.outputPath)
		//return new(excel.Encoder).WithOutputPath(b.outputPath)
	default:
		return NullEncoder
	}
}
