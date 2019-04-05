// Copyright (c) 2019 Australian Rivers Institute.

// Copyright (c) 2019 Australian Rivers Institute.

package csv

import (
	"bufio"
	"os"
	"path"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/pkg/errors"
)

type CsvEncoder struct {
	decisionVariableMarshaler CsvDecisionVariableMarshaler
	managementActionMarshaler CsvManagementActionMarshaler
	outputPath                string
}

func (ce *CsvEncoder) WithOutputPath(outputPath string) *CsvEncoder {
	ce.outputPath = outputPath
	return ce
}

func (ce CsvEncoder) Encode(solution *solution.Solution) error {
	if decisionVariableError := ce.encodeDecisionVariables(solution); decisionVariableError != nil {
		return errors.Wrap(decisionVariableError, "csv encoding of solution decision variables")
	}
	if managementActionError := ce.encodeManagementActions(solution); managementActionError != nil {
		return errors.Wrap(managementActionError, "csv encoding of solution decision variables")
	}
	return nil
}

func (ce CsvEncoder) encodeDecisionVariables(solution *solution.Solution) error {
	marshaledSolution, marshalError := ce.decisionVariableMarshaler.Marshal(solution)
	if marshalError != nil {
		return errors.Wrap(marshalError, "csv marshaling of solution decision variables")
	}

	outputPath := ce.deriveDecisionVariableOutputPath(solution)
	return ce.encodeMarshaled(marshaledSolution, outputPath)
}

func (ce CsvEncoder) encodeManagementActions(solution *solution.Solution) error {
	marshaledSolution, marshalError := ce.managementActionMarshaler.Marshal(solution)
	if marshalError != nil {
		return errors.Wrap(marshalError, "csv marshaling of solution management actions")
	}

	outputPath := ce.deriveManagementActionOutputPath(solution)
	return ce.encodeMarshaled(marshaledSolution, outputPath)
}

func (ce CsvEncoder) encodeMarshaled(marshaledSolution []byte, outputPath string) error {
	os.Remove(outputPath)

	file, openError := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE, 0666)
	if openError != nil {
		return errors.Wrap(openError, "opening file for csv encoding of solution")
	}
	defer file.Close()

	bufferedWriter := bufio.NewWriter(file)
	if _, writeError := bufferedWriter.Write(marshaledSolution); writeError != nil {
		return errors.Wrap(writeError, "writing marshaled csv of solution")
	}

	bufferedWriter.Flush()
	return nil
}

func (ce CsvEncoder) deriveDecisionVariableOutputPath(solution *solution.Solution) (outputPath string) {
	return ce.deriveOutputPath(solution, "DecisionVariables")
}

func (ce CsvEncoder) deriveManagementActionOutputPath(solution *solution.Solution) (outputPath string) {
	return ce.deriveOutputPath(solution, "ManagementActions")
}

func (ce CsvEncoder) deriveOutputPath(solution *solution.Solution, contentType string) (outputPath string) {
	safeIdBasedFileName := solution.FileNameSafeId() + "-" + contentType + ".csv"
	return path.Join(ce.outputPath, safeIdBasedFileName)
}
