// Copyright (c) 2019 Australian Rivers Institute.

package assert

import "github.com/pkg/errors"

type Assertion struct {
	condition      bool
	failureMessage string
}

const defaultFailureMessage = "assertion failed"

func That(condition bool) *Assertion {
	assertion := new(Assertion)

	assertion.condition = condition
	assertion.failureMessage = defaultFailureMessage

	return assertion
}

func (a *Assertion) WithFailureMessage(failureMessage string) *Assertion {
	a.failureMessage = failureMessage
	return a
}

func (a *Assertion) Holds() {
	if !a.condition {
		panic(errors.New(a.failureMessage))
	}
}
