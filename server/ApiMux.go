// Copyright (c) 2018 Australian Rivers Institute.

package server

const ApiMuxType = "API"

const ApiPath = "/api"
const V1Path = "/v1"

type ApiMux struct {
	BaseMux
}

func (am *ApiMux) Initialise() *ApiMux {
	am.BaseMux.Initialise().WithType(ApiMuxType)
	return am
}
