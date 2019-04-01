// Copyright (c) 2019 Australian Rivers Institute.

package solution

import (
	"bufio"
	"os"
	"path"

	"github.com/pkg/errors"
)

type JsonEncoder struct {
	marshaler  JsonMarshaler
	outputPath string
}

func (je *JsonEncoder) WithOutputPath(outputPath string) *JsonEncoder {
	je.outputPath = outputPath
	return je
}

func (je JsonEncoder) Encode(solution *Solution) error {
	marshaledSolution, marshalError := je.marshaler.Marshal(solution)
	if marshalError != nil {
		return errors.Wrap(marshalError, "json marshaling of solution")
	}

	outputPath := je.deriveOutputPath(solution)
	return je.encodeMarshaled(marshaledSolution, outputPath)
}

func (je JsonEncoder) encodeMarshaled(marshaledSolution []byte, outputPath string) error {

	file, openError := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE, 0666)
	if openError != nil {
		return errors.Wrap(openError, "opening file for json encoding of solution")
	}
	defer file.Close()

	bufferedWriter := bufio.NewWriter(file)
	if _, writeError := bufferedWriter.Write(marshaledSolution); writeError != nil {
		return errors.Wrap(writeError, "writing marshaled json of solution")
	}

	bufferedWriter.Flush()
	return nil
}

func (je JsonEncoder) deriveOutputPath(solution *Solution) (outputPath string) {
	safeIdBasedFileName := solution.FileNameSafeId() + ".json"
	return path.Join(je.outputPath, safeIdBasedFileName)
}
