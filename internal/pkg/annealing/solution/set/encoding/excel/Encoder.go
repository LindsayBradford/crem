// Copyright (c) 2019 Australian Rivers Institute.

package excel

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/set"
	"github.com/LindsayBradford/crem/internal/pkg/dataset/excel"
	"github.com/LindsayBradford/crem/pkg/logging"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"github.com/LindsayBradford/crem/pkg/threading"
	"github.com/pkg/errors"
	"os"
	"path"
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

func (e Encoder) Encode(summary *set.Summary) error {
	e.LogHandler().Info("Saving [" + summary.Id() + "] as [Excel]")

	dataSet := excel.NewDataSet(summary.FileNameSafeId(), threading.GetMainThreadChannel().Call)
	defer dataSet.Teardown()

	if marshalError := e.marshaler.Marshal(summary, dataSet); marshalError != nil {
		return errors.Wrap(marshalError, fileType+" marshaling of solution")
	}

	outputPath := e.deriveSummaryOutputPath(summary)
	e.LogHandler().Debug("Encoding [" + summary.Id() + "] to [" + outputPath + "]")
	return e.encodeMarshaled(dataSet, outputPath)

}

func (e Encoder) encodeMarshaled(dataSet *excel.DataSet, outputPath string) error {
	currentDir, _ := os.Getwd()
	absolutePath := path.Join(currentDir, outputPath)
	dataSet.SaveAs(absolutePath)
	return nil
}

func (e Encoder) deriveSummaryOutputPath(summary *set.Summary) (outputPath string) {
	return e.deriveOutputPath(summary, "Summary")
}

func (e Encoder) deriveOutputPath(summary *set.Summary, contentType string) (outputPath string) {
	safeIdBasedFileName := summary.FileNameSafeId() + "-" + contentType + fileTypeExtension
	return path.Join(outputPath, e.outputPath, safeIdBasedFileName)
}
