// Copyright (c) 2019 Australian Rivers Institute.

package solution

type Encoder interface {
	Encode(solution *Solution) error
}
