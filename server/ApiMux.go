// Copyright (c) 2018 Australian Rivers Institute.

package server

type ApiMux struct {
	BaseMux
}

func (am *ApiMux) Initialise() *ApiMux {
	am.BaseMux.Initialise()
	return am
}

func (am *ApiMux) WithType(muxType string) *ApiMux {
	am.BaseMux.WithType(muxType)
	return am
}
