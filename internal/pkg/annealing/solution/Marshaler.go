// Copyright (c) 2019 Australian Rivers Institute.

package solution

type Marshaler interface {
	Marshal(solution *Solution) ([]byte, error)
}
