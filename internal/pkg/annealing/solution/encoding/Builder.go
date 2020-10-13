// Copyright (c) 2019 Australian Rivers Institute.

package encoding

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/encoding/csv"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/encoding/excel"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/encoding/json"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
)

type OutputType string

const (
	UndefinedOutput = ""
	CsvOutput       = "CSV"
	JsonOutput      = "JSON"
	ExcelOutput     = "EXCEL"
)

type Builder struct {
	loggers.ContainedLogger
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

func (b *Builder) WithLogHandler(logHandler logging.Logger) *Builder {
	b.SetLogHandler(logHandler)
	return b
}

func (b *Builder) Build() Encoder {
	switch b.outputType {
	case UndefinedOutput, CsvOutput:
		return new(csv.Encoder).WithOutputPath(b.outputPath).WithLogHandler(b.LogHandler())
	case JsonOutput:
		return new(json.Encoder).WithOutputPath(b.outputPath).WithLogHandler(b.LogHandler())
	case ExcelOutput:
		return new(excel.Encoder).WithOutputPath(b.outputPath).WithLogHandler(b.LogHandler())
	default:
		return NullEncoder
	}
}
