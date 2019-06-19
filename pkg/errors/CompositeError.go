// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

// Package errors offers an extension of functionality to the default golang errors package.
package errors

import (
	"encoding/json"
	"errors"
	"strings"

	cremstrings "github.com/LindsayBradford/crem/pkg/strings"
)

// CompositeError allows a number of errors to be collected together and treated as if they were a single error.
type CompositeError struct {
	compositeText    string
	individualErrors []error
}

// New returns a new CompositeError. The text parameter supplied is used as a prefix when reporting a list of errors.
func New(text string) *CompositeError {
	newError := new(CompositeError).Initialise(text)
	return newError
}

func (ce *CompositeError) Initialise(text string) *CompositeError {
	ce.compositeText = text
	ce.individualErrors = make([]error, 0)
	return ce
}

// Error conforms to the built-in interface type for representing an error condition over a composed set of errors.
func (ce *CompositeError) Error() string {
	if len(ce.individualErrors) == 1 {
		return ce.individualErrors[0].Error()
	}
	return ce.buildCompositeErrorString()
}

// MarshalJSON conforms to the built-in interface type for encoding a composed set of errors as json string array.
func (ce *CompositeError) MarshalJSON() ([]byte, error) {
	errorMessages := deriveMessagesFromError(ce, 0)
	return json.Marshal(errorMessages)
}

type errorWithCause interface {
	Cause() error
}

func deriveMessagesFromError(error error, indentLevel int) []string {
	messages := make([]string, 0)

	switch error.(type) {
	case errorWithCause:
		causingError := deriveRootCauseOf(error)
		return deriveMessagesFromError(causingError, indentLevel)
	case *CompositeError:
		if errorAsComposite, errorIsComposite := error.(*CompositeError); errorIsComposite {
			newMessages := deriveMessagesFromCompositeError(errorAsComposite, indentLevel)
			messages = append(messages, newMessages...)
		}
	default:
		messages = append(messages, indent(indentLevel)+error.Error())
	}

	return messages
}

func deriveMessagesFromCompositeError(error *CompositeError, indentLevel int) []string {
	messages := make([]string, 0)
	messages = append(messages, indent(indentLevel)+error.compositeText)
	for _, individualError := range error.individualErrors {
		newMessages := deriveMessagesFromError(individualError, indentLevel+1)
		messages = append(messages, newMessages...)
	}
	return messages
}

func deriveRootCauseOf(error error) error {
	var rootError = error
	// loop through errors with causes... last error in chain (without a cause) IS the root cause.
	for typedError, hasCause := rootError.(errorWithCause); hasCause; typedError, hasCause = rootError.(errorWithCause) {
		rootError = typedError.Cause()
	}
	return rootError
}

func indent(indentLevel int) string {
	const indentString = " "
	return strings.Repeat(indentString, indentLevel)
}

func (ce *CompositeError) buildCompositeErrorString() string {
	builder := cremstrings.FluentBuilder{}

	builder.Add(ce.compositeText, ", composed of: [\n")

	for _, currError := range ce.individualErrors {
		builder.Add("\t", currError.Error(), "\n")
	}

	builder.Add("]")

	return builder.String()
}

// ColumnAndRowSize returns the number of sub-errors in a CompositeError. It counts only directly accessible sub-errors (no nesting).
func (ce *CompositeError) Size() int {
	return len(ce.individualErrors)
}

// Add includes newError as one of the sub-errors of a CompositeError
func (ce *CompositeError) Add(newError error) {
	switch typedError := newError.(type) {
	case nil:
	case *CompositeError:
		if typedError.Size() > 0 {
			ce.individualErrors = append(ce.individualErrors, newError)
		}
	default:
		ce.individualErrors = append(ce.individualErrors, newError)
	}
}

// Add combines the supplied message as a new built-in error to the array of sub-errors of a CompositeError
func (ce *CompositeError) AddMessage(message string) {
	newError := errors.New(message)
	ce.Add(newError)
}

// SubError returns the sub-error at the given array index of a CompositeError
func (ce *CompositeError) SubError(index int) error {
	return ce.individualErrors[index]
}
