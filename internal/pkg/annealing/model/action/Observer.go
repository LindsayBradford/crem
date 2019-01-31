// Copyright (c) 2019 Australian Rivers Institute.

package action

type Observer interface {
	Observe(action ManagementAction)
}
