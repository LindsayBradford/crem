// Copyright (c) 2019 Australian Rivers Institute.

package assert

// GoLand 2019.1 isn't great at conditional compilation & running of build tagged code.
// Workaround is to have debug/ and release/ package variants of assert, and switch between them
// on a source-file by source-file basis.
// NOTE: Revisit Goland 2019.1 conditional test compilation WRT assertions.

type RuntimeAssertion struct{}

var nullAssertion = new(RuntimeAssertion)

func That(condition bool) *RuntimeAssertion {
	return nullAssertion
}

func (a *RuntimeAssertion) WithFailureMessage(messageOnFailure string) *RuntimeAssertion {
	return nullAssertion
}

func (a *RuntimeAssertion) Holds() {}
