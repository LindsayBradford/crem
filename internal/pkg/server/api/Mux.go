// Copyright (c) 2018 Australian Rivers Institute.

package api

import (
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
)

const muxType = "API"

const BasePath = "api"

type Mux struct {
	rest.MuxImpl
}

func (m *Mux) Initialise() *Mux {
	m.MuxImpl.Initialise().WithType(muxType)
	return m
}
