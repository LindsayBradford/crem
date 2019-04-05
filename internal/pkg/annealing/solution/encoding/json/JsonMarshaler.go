// Copyright (c) 2019 Australian Rivers Institute.

// Copyright (c) 2019 Australian Rivers Institute.

package json

import (
	"encoding/json"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
)

type JsonMarshaler struct{}

func (jm *JsonMarshaler) Marshal(solution *solution.Solution) ([]byte, error) {
	const (
		newLinePrefix = ""
		indent        = "  "
	)

	return json.MarshalIndent(solution, newLinePrefix, indent)
}
