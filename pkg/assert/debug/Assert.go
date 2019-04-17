// Copyright (c) 2019 Australian Rivers Institute.

package assert

import "github.com/pkg/errors"

type RuntimeAssertion struct {
	condition      bool
	failureMessage string
}

const defaultFailureMessage = "assertion failed"

func That(condition bool) *RuntimeAssertion {
	assertion := new(RuntimeAssertion)

	assertion.condition = condition
	assertion.failureMessage = defaultFailureMessage

	return assertion
}

func (a *RuntimeAssertion) WithFailureMessage(failureMessage string) *RuntimeAssertion {
	a.failureMessage = failureMessage
	return a
}

func (a *RuntimeAssertion) Holds() {
	if !a.condition {
		panic(errors.New(a.failureMessage))
	}
}
