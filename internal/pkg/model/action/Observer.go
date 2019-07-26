// Copyright (c) 2019 Australian Rivers Institute.

package action

// Reporting defines an interface allowing the observation of management actions via callback methods.
type Observer interface {
	// ObserveActionInitialising is a callback method, invoked when a management action is activated as part of any
	// necessary model initialisation.
	ObserveActionInitialising(action ManagementAction)

	// ObserverAction is a callback method, invoked when a management action's activation changes as a model runs
	// (only after the model has been initialised).
	ObserveAction(action ManagementAction)
}
