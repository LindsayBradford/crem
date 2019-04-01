// Copyright (c) 2019 Australian Rivers Institute.

package solution

import "encoding/json"

type JsonMarshaler struct{}

func (jm *JsonMarshaler) Marshal(solution *Solution) ([]byte, error) {
	const (
		newLinePrefix = ""
		indent        = "  "
	)

	return json.MarshalIndent(solution, newLinePrefix, indent)
}
