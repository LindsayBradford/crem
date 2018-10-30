// Copyright (c) 2018 Australian Rivers Institute.

package uuid

import "github.com/nu7hatch/gouuid"

func New() string {
	newUuid, _ := uuid.NewV4()
	return newUuid.String()
}
