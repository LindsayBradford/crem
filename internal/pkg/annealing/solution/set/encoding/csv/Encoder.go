// Copyright (c) 2019 Australian Rivers Institute.

package csv

import (
	"bufio"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/set"
	"os"
	"path"

	"github.com/pkg/errors"
)

const fileType = "csv"
const fileTypeExtension = "." + fileType

type Encoder struct {
	summaryMarshaler SummaryMarshaler
	outputPath       string
}

func (e *Encoder) WithOutputPath(outputPath string) *Encoder {
	e.outputPath = outputPath
	return e
}

func (e Encoder) Encode(summary *set.Summary) error {
	if decisionVariableError := e.encodeDecisionVariables(summary); decisionVariableError != nil {
		return errors.Wrap(decisionVariableError, fileType+" encoding of solution decision variables")
	}
	return nil
}

func (e Encoder) encodeDecisionVariables(summary *set.Summary) error {
	marshaledSolution, marshalError := e.summaryMarshaler.Marshal(summary)
	if marshalError != nil {
		return errors.Wrap(marshalError, fileType+" marshaling of solution decision variables")
	}

	outputPath := e.deriveSummaryOutputPath(summary)
	return e.encodeMarshaled(marshaledSolution, outputPath)
}

func (e Encoder) encodeMarshaled(marshaledSummary []byte, outputPath string) error {
	os.Remove(outputPath)

	file, openError := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if openError != nil {
		return errors.Wrap(openError, "opening file for "+fileType+" encoding of solution summary")
	}
	defer file.Close()

	bufferedWriter := bufio.NewWriter(file)
	if _, writeError := bufferedWriter.Write(marshaledSummary); writeError != nil {
		return errors.Wrap(writeError, "writing marshaled "+fileType+" of solution summary")
	}

	bufferedWriter.Flush()
	return nil
}

func (e Encoder) deriveSummaryOutputPath(summary *set.Summary) (outputPath string) {
	return e.deriveOutputPath(summary, "Summary")
}

func (e Encoder) deriveOutputPath(summary *set.Summary, contentType string) (outputPath string) {
	safeIdBasedFileName := summary.FileNameSafeId() + "-" + contentType + fileTypeExtension
	return path.Join(e.outputPath, safeIdBasedFileName)
}
