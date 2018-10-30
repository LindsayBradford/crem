// Copyright (c) 2018 Australian Rivers Institute.

package api

import (
	"github.com/LindsayBradford/crem/server/rest"
)

const muxType = "API"

const BasePath = "/api"
const V1Path = "/v1"

type Mux struct {
	rest.MuxImpl
}

func (m *Mux) Initialise() *Mux {
	m.MuxImpl.Initialise().WithType(muxType)
	return m
}
