// Copyright (c) 2019 Australian Rivers Institute.

package encoding

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/encoding"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/set/encoding/csv"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/set/encoding/excel"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/set/encoding/json"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	//"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/set/encoding/excel"
	//"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/set/encoding/json"
)

type Builder struct {
	loggers.ContainedLogger
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

func (b *Builder) WithLogHandler(logHandler logging.Logger) *Builder {
	b.SetLogHandler(logHandler)
	return b
}

func (b *Builder) Build() Encoder {
	switch b.outputType {
	case encoding.UndefinedOutput, encoding.CsvOutput:
		return new(csv.Encoder).WithOutputPath(b.outputPath).WithLogHandler(b.LogHandler())
	case encoding.JsonOutput:
		return new(json.Encoder).WithOutputPath(b.outputPath).WithLogHandler(b.LogHandler())
	case encoding.ExcelOutput:
		return new(excel.Encoder).WithOutputPath(b.outputPath).WithLogHandler(b.LogHandler())
	default:
		return NullEncoder
	}
}
