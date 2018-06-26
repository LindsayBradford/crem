// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

// Package errors offers an extension of functionality to the default golang errors package.
package errors

import "github.com/LindsayBradford/crm/strings"

// CompositeError offers a convenience wrapper to a number of related error instances.
// It allows a number of errors to be collected together and delivered  as if they were one error, along with the
// ability to learn more about individual errors if needed.
type CompositeError struct {
	compositeText    string
	individualErrors [] error
}

// NewComposite returns a CompositeError that formats as the given text prefixing a list of error texts for those
// errors that it is composed of.
func NewComposite(text string) *CompositeError {
	newError := new(CompositeError)
	newError.compositeText = text
	return newError
}

func (this *CompositeError) Error() string {
	if (len(this.individualErrors) == 1) {
		return this.individualErrors[0].Error()
	}
	return this.buildCompositeErrorString()
}

func (this *CompositeError) buildCompositeErrorString() string {
	builder := strings.FluentBuilder{}

	builder.Add(this.compositeText, ", composed of: [\n")

	for _, currError := range this.individualErrors {
		builder.Add("\t", currError.Error(), "\n")
	}

	builder.Add("]")

	return builder.String()
}

// Size returns the number of sub-errors that have been composed together to form the given CompositeError
func (this* CompositeError) Size() int {
	return len(this.individualErrors)
}

// Add combines newError to the array of sub-errors that have been composed together to form the given CompositeError
func (this *CompositeError) Add(newError error) {
	this.individualErrors = append(this.individualErrors, newError)
}

//SubError returns the sub-error at the index specified by position for the given CompositeError
func (this *CompositeError) SubError(position int) error {
	return this.individualErrors[position]
}