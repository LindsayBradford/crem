// Copyright (c) 2018 Australian Rivers Institute.

package api

import "github.com/LindsayBradford/crem/server"

const muxType = "API"

const BasePath = "/api"
const V1Path = "/v1"

type Mux struct {
	server.BaseMux
}

func (m *Mux) Initialise() *Mux {
	m.BaseMux.Initialise().WithType(muxType)
	return m
}
