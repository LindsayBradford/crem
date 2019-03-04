// Copyright (c) 2019 Australian Rivers Institute.

package variable

// BaseInductiveDecisionVariable is a DecisionVariable that allows an 'inductive' value to be temporarily stored
// and retrieved for the decision variable (typically based based on some management action).
// The induced value does not become the actual value for the decision variable without being explicitly accepted.
// The induced value can also be rejected, which sees it revert to the actual value of the variable.
type InductiveDecisionVariable interface {
	DecisionVariable

	InductiveValue() float64
	SetInductiveValue(value float64)
	DifferenceInValues() float64

	AcceptInductiveValue()
	RejectInductiveValue()
}
