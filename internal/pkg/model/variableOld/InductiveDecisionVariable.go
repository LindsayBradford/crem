// Copyright (c) 2019 Australian Rivers Institute.

package variableOld

import (
	"github.com/LindsayBradford/crem/internal/pkg/model/action"
	"github.com/LindsayBradford/crem/internal/pkg/model/variable"
)

// InductiveDecisionVariable is a DecisionVariable that allows an 'inductive' Value to be temporarily stored
// and retrieved for the decision variableOld (typically based based on some management action change).
// The induced Value does not become the actual Value for the decision variableOld without being explicitly accepted.
// The induced Value can also be rejected, which sees it revert to the actual Value of the variableOld.
type InductiveDecisionVariable interface {
	variable.DecisionVariable

	// InductiveValue returns an "inductive" Value for the variableOld.  This Value cannot become the
	// actual (induced) Value for the variableOld without a call to AcceptInductiveValue.
	InductiveValue() float64

	// SetInductiveValue allows an "inductive" Value for the variableOld to be set.  This inductive Value is not the
	// variableOld's actual. It is a Value that would result (or be induced) from, some management action change.
	// This Value  is expected to be a temporary, lasting only as long as it takes for decision-making logic to decide
	// whether to induce the variableOld to this Value, or reject it.
	SetInductiveValue(value float64)

	// DifferenceInValues report the difference in taking the variableOld's actual Value for its inductive Value.
	// This difference is often used in decision making around whether to accept the inductive Value.
	DifferenceInValues() float64

	// Accepts the inductive Value of the variableOld as the variableOld's actual Value.
	AcceptInductiveValue()

	// Rejects the inductive Value of the variableOld, resetting the inductive Value to the variableOld's actual Value.
	RejectInductiveValue()

	action.Observer
}
