// Copyright (c) 2019 Australian Rivers Institute.

package encoding

import (
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/set"
)

type Marshaler interface {
	Marshal(summary *set.Summary) ([]byte, error)
}
