// Copyright (c) 2019 Australian Rivers Institute.

package variable

// UndoableDecisionVariable is a DecisionVariable that allows an 'inductive' Value to be temporarily stored
// and retrieved for the decision variableOld (typically based based on some management action change).
// The induced Value does not become the actual Value for the decision variableOld without being explicitly accepted.
// The induced Value can also be rejected, which sees it revert to the actual Value of the variableOld.
type UndoableDecisionVariable interface {
	DecisionVariable

	// UndoableValue returns an "inductive" Value for the variableOld.  This Value cannot become the
	// actual (induced) Value for the variableOld without a call to ApplyDoneValue.
	UndoableValue() float64

	// SetUndoableValue allows an "undoable"vValue for the variable to be set.  This undoable Value is not the
	// variable's actual value. It is a Value that would result (or be induced) from, some management action change.
	// This Value  is expected to be a temporary, lasting only as long as it takes for decision-making logic to decide
	// whether to accept the variable value, or reject/undo it.
	SetUndoableValue(value float64)

	// DifferenceInValues report the difference in taking the variableOld's actual Value for its inductive Value.
	// This difference is often used in decision making around whether to accept the inductive Value.
	DifferenceInValues() float64

	// Accepts the inductive Value of the variableOld as the variableOld's actual Value.
	ApplyDoneValue()

	// Rejects the inductive Value of the variableOld, resetting the inductive Value to the variableOld's actual Value.
	ApplyUndoneValue()
}
