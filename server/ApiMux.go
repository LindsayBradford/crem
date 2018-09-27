// Copyright (c) 2018 Australian Rivers Institute.

package server

type ApiMux struct {
	RestMux
}

func (am *ApiMux) Initialise() *ApiMux {
	am.RestMux.Initialise()
	return am
}

func (am *ApiMux) WithType(muxType string) *ApiMux {
	am.RestMux.WithType(muxType)
	return am
}
