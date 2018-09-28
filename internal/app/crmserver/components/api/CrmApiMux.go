// Copyright (c) 2018 Australian Rivers Institute.

package api

import (
	"fmt"
	"net/http"

	"github.com/LindsayBradford/crm/config"
	"github.com/LindsayBradford/crm/server"
)

const apiPath = "/api"
const v1Path = "/v1"

type CrmApiMux struct {
	server.ApiMux
}

func (cam *CrmApiMux) Initialise() *CrmApiMux {
	cam.ApiMux.Initialise()
	cam.AddHandler(apiPath+v1Path+"/scenarios", cam.v1PostScenario)
	return cam
}

func (cam *CrmApiMux) WithType(muxType string) *CrmApiMux {
	cam.ApiMux.WithType(muxType)
	return cam
}

func (cam *CrmApiMux) rootPathHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, nameAndVersionString())
}

func nameAndVersionString() string {
	return fmt.Sprintf("%s, version %s", config.LongApplicationName, config.Version)
}
