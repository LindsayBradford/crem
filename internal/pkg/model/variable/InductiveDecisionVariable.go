// Copyright (c) 2019 Australian Rivers Institute.

package variable

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/variableNew"
)

// InductiveDecisionVariable is a DecisionVariable that allows an 'inductive' Value to be temporarily stored
// and retrieved for the decision variable (typically based based on some management action change).
// The induced Value does not become the actual Value for the decision variable without being explicitly accepted.
// The induced Value can also be rejected, which sees it revert to the actual Value of the variable.
type InductiveDecisionVariable interface {
	variableNew.DecisionVariable

	// InductiveValue returns an "inductive" Value for the variable.  This Value cannot become the
	// actual (induced) Value for the variable without a call to AcceptInductiveValue.
	InductiveValue() float64

	// SetInductiveValue allows an "inductive" Value for the variable to be set.  This inductive Value is not the
	// variable's actual. It is a Value that would result (or be induced) from, some management action change.
	// This Value  is expected to be a temporary, lasting only as long as it takes for decision-making logic to decide
	// whether to induce the variable to this Value, or reject it.
	SetInductiveValue(value float64)

	// DifferenceInValues report the difference in taking the variable's actual Value for its inductive Value.
	// This difference is often used in decision making around whether to accept the inductive Value.
	DifferenceInValues() float64

	// Accepts the inductive Value of the variable as the variable's actual Value.
	AcceptInductiveValue()

	// Rejects the inductive Value of the variable, resetting the inductive Value to the variable's actual Value.
	RejectInductiveValue()

	action.Observer
}
