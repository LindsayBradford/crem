// Copyright (c) 2019 Australian Rivers Institute.

package variable

// InductiveDecisionVariable is a DecisionVariable that allows an 'inductive' value to be temporarily stored
// and retrieved for the decision variable (typically based based on some management action change).
// The induced value does not become the actual value for the decision variable without being explicitly accepted.
// The induced value can also be rejected, which sees it revert to the actual value of the variable.
type InductiveDecisionVariable interface {
	DecisionVariable

	// InductiveValue returns an "inductive" value for the variable.  This value cannot become the
	// actual (induced) value for the variable without a call to AcceptInductiveValue.
	InductiveValue() float64

	// SetInductiveValue allows an "inductive" value for the variable to be set.  This inductive value is not the
	// variable's actual. It is a value that would result (or be induced) from, some management action change.
	// This value  is expected to be a temporary, lasting only as long as it takes for decision-making logic to decide
	// whether to induce the variable to this value, or reject it.
	SetInductiveValue(value float64)

	// DifferenceInValues report the difference in taking the variable's actual value for its inductive value.
	// This difference is often used in decision making around whether to accept the inductive value.
	DifferenceInValues() float64

	// Accepts the inductive value of the variable as the variable's actual value.
	AcceptInductiveValue()

	// Rejects the inductive value of the variable, resetting the inductive value to the variable's actual value.
	RejectInductiveValue()
}
