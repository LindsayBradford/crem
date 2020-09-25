// Copyright (c) 2019 Australian Rivers Institute.

package encoding

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/set"
)

type Encoder interface {
	Encode(summary *set.Summary) error
}

var NullEncoder = new(nullEncoder)

type nullEncoder struct{}

func (ne *nullEncoder) Encode(summary *set.Summary) error { return nil }
