// Copyright (c) 2019 Australian Rivers Institute.

package excel

import (
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
	marshaler  Marshaler
	outputPath string
}

func (e *Encoder) WithOutputPath(outputPath string) *Encoder {
	e.outputPath = outputPath
	return e
}

func (e Encoder) Encode(solution *solution.Solution) error {
	dataSet := excel.NewDataSet(solution.FileNameSafeId(), threading.GetMainThreadChannel().Call)
	defer dataSet.Teardown()

	if marshalError := e.marshaler.Marshal(solution, dataSet); marshalError != nil {
		return errors.Wrap(marshalError, fileType+" marshaling of solution")
	}

	outputPath := e.deriveOutputPath(solution)
	return e.encodeMarshaled(dataSet, outputPath)
}

func (e Encoder) encodeMarshaled(dataSet *excel.DataSet, outputPath string) error {
	dataSet.SaveAs(outputPath)
	return nil
}

func (e Encoder) deriveOutputPath(solution *solution.Solution) (outputPath string) {
	currentDir, _ := os.Getwd()
	safeIdBasedFileName := solution.FileNameSafeId() + fileTypeExtension
	return path.Join(currentDir, e.outputPath, safeIdBasedFileName)
}
