// Copyright (c) 2018 Australian Rivers Institute. Author: Lindsay Bradford

package errors

import "github.com/LindsayBradford/crm/strings"

type CompositeError struct {
	compositeMessage string
	individualErrors[] error
}

func NewComposite(message string) *CompositeError {
	newError := new(CompositeError)
	newError.compositeMessage = message
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

	builder.Add(this.compositeMessage, ", composed of: [\n")

	for _, currError := range this.individualErrors {
		builder.Add("\t", currError.Error(), "\n")
	}

	builder.Add("]")

	return builder.String()
}

func (this* CompositeError) Size() int {
	return len(this.individualErrors)
}

func (this *CompositeError) Add(newError error) {
	this.individualErrors = append(this.individualErrors, newError)
}

func (this *CompositeError) SubError(position int) error {
	return this.individualErrors[position]
}