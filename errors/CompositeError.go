// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

// Package errors offers an extension of functionality to the default golang errors package.
package errors

import (
	"encoding/json"
	"errors"
	"strings"

	cremstrings "github.com/LindsayBradford/crem/strings"
)

// CompositeError offers a convenience wrapper to a number of related error instances.
// It allows a number of errors to be collected together and delivered  as if they were one error, along with the
// ability to learn more about individual errors if needed.
type CompositeError struct {
	compositeText    string  `json:"Summary"`
	individualErrors []error `json:"Errors"`
}

// NewComposite returns a CompositeError that formats as the given text prefixing a list of error texts for those
// errors that it is composed of.
func NewComposite(text string) *CompositeError {
	newError := new(CompositeError)
	newError.compositeText = text
	return newError
}

func (ce *CompositeError) Error() string {
	if len(ce.individualErrors) == 1 {
		return ce.individualErrors[0].Error()
	}
	return ce.buildCompositeErrorString()
}

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

// Size returns the number of sub-errors that have been composed together to form the given CompositeError
func (ce *CompositeError) Size() int {
	return len(ce.individualErrors)
}

// Add combines newError to the array of sub-errors that have been composed together to form the given CompositeError
func (ce *CompositeError) Add(newError error) {
	ce.individualErrors = append(ce.individualErrors, newError)
}

// Add combines message as a new error  to the array of sub-errors that have been composed together to form the given CompositeError
func (ce *CompositeError) AddMessage(message string) {
	newError := errors.New(message)
	ce.Add(newError)
}

// SubError returns the sub-error at the index specified by position for the given CompositeError
func (ce *CompositeError) SubError(position int) error {
	return ce.individualErrors[position]
}
