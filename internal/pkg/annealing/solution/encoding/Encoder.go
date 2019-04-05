// Copyright (c) 2019 Australian Rivers Institute.

package encoding

import "github.com/LindsayBradford/crem/internal/pkg/annealing/solution"

type Encoder interface {
	Encode(solution *solution.Solution) error
}

var NullEncoder = new(nullEncoder)

type nullEncoder struct{}

func (ne *nullEncoder) Encode(solution *solution.Solution) error { return nil }
