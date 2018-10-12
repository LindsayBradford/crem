// Copyright (c) 2018 Australian Rivers Institute.

package api

import (
	"github.com/LindsayBradford/crm/config"
	"github.com/LindsayBradford/crm/server"
)

const jobsPath = "/jobs"

type CrmApiMux struct {
	server.ApiMux
}

type ScenarioJob struct {
	jobId          string
	status         string
	ScenarioConfig *config.CRMConfig
}

type ScenarioJobQueue struct {
	jobs []ScenarioJob
}

func (cam *CrmApiMux) Initialise() *CrmApiMux {
	cam.ApiMux.Initialise()
	cam.AddHandler(baseApiPath()+jobsPath, cam.V1HandleJobs)
	return cam
}

func baseApiPath() string {
	return server.ApiPath + server.V1Path
}

func (cam *CrmApiMux) WithType(muxType string) *CrmApiMux {
	cam.ApiMux.WithType(muxType)
	return cam
}
