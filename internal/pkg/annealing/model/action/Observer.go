// Copyright (c) 2019 Australian Rivers Institute.

package action

type Observer interface {
	ObserveAction(action ManagementAction)
	ObserveActionInitialising(action ManagementAction)
}
