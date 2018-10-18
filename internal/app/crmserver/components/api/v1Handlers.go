// Copyright (c) 2018 Australian Rivers Institute.

package api

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/LindsayBradford/crm/config"
	"github.com/LindsayBradford/crm/internal/app/crmserver/components/scenario"
	"github.com/LindsayBradford/crm/server"
	"github.com/pkg/errors"
)

const scenarioConfigKey = "ScenarioConfig"
const completedTimeKey = "CompletedTime"

const scenarioPath = "scenario"

func (cam *CrmApiMux) V1HandleJobs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		cam.v1PostJob(w, r)
	case http.MethodGet:
		cam.v1GetJobs(w, r)
	default:
		cam.MethodNotAllowedError(w, r)
	}
}

func (cam *CrmApiMux) v1PostJob(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(server.ContentTypeHeaderKey) != server.TomlMimeType {
		cam.MethodNotAllowedError(w, r)
		return
	}

	scenarioText := requestBodyToString(r)

	newJob := new(Job).Initialise()
	newJob.HiddenAttributes[scenarioConfigKey] = scenarioText

	cam.jobs.Enqueue(*newJob)
	cam.AddHandler(BuildApiPath(jobsPath, newJob.JobId), cam.v1GetJob)
	cam.AddHandler(BuildApiPath(jobsPath, newJob.JobId, scenarioPath), cam.v1GetJobScenario)

	scenarioConfig, retrieveError := config.RetrieveCrmFromString(scenarioText)

	if retrieveError != nil {
		wrappingError := errors.Wrap(retrieveError, "retrieving scenario configuration")
		cam.Logger().Warn(wrappingError)
		cam.InternalServerError(w, r, errors.New("Invalid scenario configuration supplied"))
		return
	}

	cam.Logger().Info("Received new job [" + newJob.JobId + "] with scenario name [" + scenarioConfig.ScenarioName + "].")

	restResponse := new(server.RestResponse).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusCreated).
		WithCacheControlMaxAge(cam.CacheMaxAge()).
		WithJsonContent(newJob)

	cam.writeResponse(restResponse, "create job")

	// TODO: Push the actual running of the scenario to a concurrent channel.

	cam.Logger().Info("Running Job [" + newJob.JobId + "].")

	scenario.RunScenarioFromConfig(scenarioConfig)

	newJob.Attributes[statusKey] = "COMPLETED"
	newJob.Attributes[completedTimeKey] = server.FormattedTimestamp()
}

func requestBodyToString(r *http.Request) string {
	responseBodyBytes, _ := ioutil.ReadAll(r.Body)
	return string(responseBodyBytes)
}

func (cam *CrmApiMux) v1GetJobs(w http.ResponseWriter, r *http.Request) {
	restResponse := new(server.RestResponse).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(cam.CacheMaxAge()).
		WithJsonContent(cam.jobs)

	cam.writeResponse(restResponse, "get jobs")
}

func (cam *CrmApiMux) v1GetJob(w http.ResponseWriter, r *http.Request) {
	if cam.expectedMethodNotSupplied(http.MethodGet, w, r) {
		return
	}

	job := cam.deriveJobFromGetJobRequest(r)

	restResponse := new(server.RestResponse).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(cam.CacheMaxAge()).
		WithJsonContent(job)

	cam.writeResponse(restResponse, "get job")
}

func (cam *CrmApiMux) deriveJobFromGetJobRequest(r *http.Request) Job {
	desiredId := getJobIdFromRequestUrl(r)
	return cam.jobs.JobWithId(desiredId)
}

func getJobIdFromRequestUrl(r *http.Request) string {
	splitURL := strings.Split(r.URL.Path, server.UrlPathSeparator)
	desiredId := splitURL[len(splitURL)-1]
	return desiredId
}

func (cam *CrmApiMux) v1GetJobScenario(w http.ResponseWriter, r *http.Request) {
	if cam.expectedMethodNotSupplied(http.MethodGet, w, r) {
		return
	}

	job := cam.deriveJobFromGetJobScenarioRequest(r)

	restResponse := new(server.RestResponse).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlPublic().
		WithTomlContent(scenarioOf(job))

	cam.writeResponse(restResponse, "get job scenario")
}

func scenarioOf(job Job) interface{} {
	return job.HiddenAttributes[scenarioConfigKey]
}

func (cam *CrmApiMux) expectedMethodNotSupplied(expectedHttpMethod string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != expectedHttpMethod {
		cam.MethodNotAllowedError(w, r)
		return true
	}
	return false
}

func (cam *CrmApiMux) writeResponse(response *server.RestResponse, context string) {
	writeError := response.Write()
	if writeError != nil {
		wrappingError := errors.Wrap(writeError, context)
		cam.Logger().Error(wrappingError)
	}
}

func (cam *CrmApiMux) deriveJobFromGetJobScenarioRequest(r *http.Request) Job {
	desiredId := getJobIdFromScenarioRequestUrl(r)
	return cam.jobs.JobWithId(desiredId)
}

func getJobIdFromScenarioRequestUrl(r *http.Request) string {
	splitURL := strings.Split(r.URL.Path, server.UrlPathSeparator)
	desiredId := splitURL[len(splitURL)-2]
	return desiredId
}
