// Copyright (c) 2019 Australian Rivers Institute.

package excel

import (
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"os"
	"path"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/excel"
	"github.com/LindsayBradford/crem/pkg/threading"
	"github.com/pkg/errors"
)

const fileType = "xlsx"
const fileTypeExtension = "." + fileType

type Encoder struct {
	loggers.ContainedLogger
	marshaler  Marshaler
	outputPath string
}

func (e *Encoder) WithOutputPath(outputPath string) *Encoder {
	e.outputPath = outputPath
	return e
}

func (e *Encoder) WithLogHandler(logHandler logging.Logger) *Encoder {
	e.SetLogHandler(logHandler)
	return e
}

func (e Encoder) Encode(solution *solution.Solution) error {
	e.LogHandler().Info("Saving [" + solution.Id + "] as [Excel]")

	dataSet := excel.NewDataSet(solution.FileNameSafeId(), threading.GetMainThreadChannel().Call)
	defer dataSet.Teardown()

	if marshalError := e.marshaler.Marshal(solution, dataSet); marshalError != nil {
		return errors.Wrap(marshalError, fileType+" marshaling of solution")
	}

	outputPath := e.deriveOutputPath(solution)
	e.LogHandler().Debug("Encoding [" + solution.Id + "] to [" + outputPath + "]")
	return e.encodeMarshaled(dataSet, outputPath)
}

func (e Encoder) encodeMarshaled(dataSet *excel.DataSet, outputPath string) error {
	currentDir, _ := os.Getwd()
	absolutePath := path.Join(currentDir, outputPath)
	dataSet.SaveAs(absolutePath)
	return nil
}

func (e Encoder) deriveOutputPath(solution *solution.Solution) (outputPath string) {
	safeIdBasedFileName := solution.FileNameSafeId() + fileTypeExtension
	return path.Join(outputPath, e.outputPath, safeIdBasedFileName)
}
