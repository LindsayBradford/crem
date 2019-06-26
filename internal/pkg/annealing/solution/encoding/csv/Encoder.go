// Copyright (c) 2019 Australian Rivers Institute.

package csv

import (
	"bufio"
	"os"
	"path"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/pkg/errors"
)

const fileType = "csv"
const fileTypeExtension = "." + fileType

type Encoder struct {
	decisionVariableMarshaler DecisionVariableMarshaler
	managementActionMarshaler ManagementActionMarshaler
	outputPath                string
}

func (e *Encoder) WithOutputPath(outputPath string) *Encoder {
	e.outputPath = outputPath
	return e
}

func (e Encoder) Encode(solution *solution.Solution) error {
	if decisionVariableError := e.encodeDecisionVariables(solution); decisionVariableError != nil {
		return errors.Wrap(decisionVariableError, fileType+" encoding of solution decision variables")
	}
	if managementActionError := e.encodeManagementActions(solution); managementActionError != nil {
		return errors.Wrap(managementActionError, fileType+" encoding of solution decision variables")
	}
	return nil
}

func (e Encoder) encodeDecisionVariables(solution *solution.Solution) error {
	marshaledSolution, marshalError := e.decisionVariableMarshaler.Marshal(solution)
	if marshalError != nil {
		return errors.Wrap(marshalError, fileType+" marshaling of solution decision variables")
	}

	outputPath := e.deriveDecisionVariableOutputPath(solution)
	return e.encodeMarshaled(marshaledSolution, outputPath)
}

func (e Encoder) encodeManagementActions(solution *solution.Solution) error {
	marshaledSolution, marshalError := e.managementActionMarshaler.Marshal(solution)
	if marshalError != nil {
		return errors.Wrap(marshalError, fileType+" marshaling of solution management actions")
	}

	outputPath := e.deriveManagementActionOutputPath(solution)
	return e.encodeMarshaled(marshaledSolution, outputPath)
}

func (e Encoder) encodeMarshaled(marshaledSolution []byte, outputPath string) error {
	os.Remove(outputPath)

	file, openError := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if openError != nil {
		return errors.Wrap(openError, "opening file for "+fileType+" encoding of solution")
	}
	defer file.Close()

	bufferedWriter := bufio.NewWriter(file)
	if _, writeError := bufferedWriter.Write(marshaledSolution); writeError != nil {
		return errors.Wrap(writeError, "writing marshaled "+fileType+" of solution")
	}

	bufferedWriter.Flush()
	return nil
}

func (e Encoder) deriveDecisionVariableOutputPath(solution *solution.Solution) (outputPath string) {
	return e.deriveOutputPath(solution, "DecisionVariables")
}

func (e Encoder) deriveManagementActionOutputPath(solution *solution.Solution) (outputPath string) {
	return e.deriveOutputPath(solution, "ManagementActions")
}

func (e Encoder) deriveOutputPath(solution *solution.Solution, contentType string) (outputPath string) {
	safeIdBasedFileName := solution.FileNameSafeId() + "-" + contentType + fileTypeExtension
	return path.Join(e.outputPath, safeIdBasedFileName)
}
