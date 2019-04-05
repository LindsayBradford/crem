// Copyright (c) 2019 Australian Rivers Institute.

package encoding

import "github.com/LindsayBradford/crem/internal/pkg/annealing/solution"

type Marshaler interface {
	Marshal(solution *solution.Solution) ([]byte, error)
}
