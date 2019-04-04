// Copyright (c) 2019 Australian Rivers Institute.

package solution

var NullEncoder = new(nullEncoder)

type nullEncoder struct{}

func (ne *nullEncoder) Encode(solution *Solution) error { return nil }
