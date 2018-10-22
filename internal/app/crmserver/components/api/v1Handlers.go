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

const scenarioConfigTextKey = "ScenarioConfigText"
const scenarioConfigStructKey = "ScenarioConfigStruct"

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
	if cam.requestContentTypeWasNotToml(r, w) {
		return
	}

	newJob := cam.createJobFromRequest(r)

	if cam.scenarioTextSuppliedIsInvalid(newJob, w, r) || bool(cam.jobEnqueueFailed(newJob, w, r)) {
		return
	}

	cam.writeJobAsResponse(newJob, w)
	cam.allowJobQueries(newJob)
}

func (cam *CrmApiMux) createJobFromRequest(r *http.Request) *server.Job {
	newJob := new(server.Job).Initialise()
	newJob.HiddenAttributes[scenarioConfigTextKey] = requestBodyToString(r)
	return newJob
}

func (cam *CrmApiMux) requestContentTypeWasNotToml(r *http.Request, w http.ResponseWriter) bool {
	if r.Header.Get(server.ContentTypeHeaderKey) != server.TomlMimeType {
		cam.MethodNotAllowedError(w, r)
		return true
	}
	return false
}

func (cam *CrmApiMux) scenarioTextSuppliedIsInvalid(job *server.Job, w http.ResponseWriter, r *http.Request) bool {
	scenarioText := job.HiddenAttributes[scenarioConfigTextKey].(string)
	scenarioConfig, retrieveError := config.RetrieveCrmFromString(scenarioText)

	if retrieveError != nil {
		cam.handleScenarioConfigRetrieveError(retrieveError, job, w, r)
		return true
	}

	cam.Logger().Info("Received new job [" + string(job.Id) + "] with scenario name [" + scenarioConfig.ScenarioName + "].")
	job.HiddenAttributes[scenarioConfigStructKey] = scenarioConfig
	return false
}

func (cam *CrmApiMux) handleScenarioConfigRetrieveError(retrieveError error, job *server.Job, w http.ResponseWriter, r *http.Request) {
	wrappingError := errors.Wrap(retrieveError, "retrieving scenario configuration")
	cam.Logger().Warn(wrappingError)
	job.SetStatus(server.JobInvalid)
	cam.AddToHistory(job)
	cam.InternalServerError(w, r, errors.New("Invalid scenario configuration supplied"))
	cam.allowJobQueries(job)
}

func (cam *CrmApiMux) allowJobQueries(job *server.Job) {
	cam.AddHandler(BuildApiPath(jobsPath, string(job.Id)), cam.v1GetJob)
	cam.AddHandler(BuildApiPath(jobsPath, string(job.Id), scenarioPath), cam.v1GetJobScenario)
}

func (cam *CrmApiMux) jobEnqueueFailed(job *server.Job, w http.ResponseWriter, r *http.Request) server.JobEnqueuedStatus {
	cam.AddToHistory(job)
	enqueuedStatus := cam.jobs.Enqueue(job)
	if enqueuedStatus == server.JobEnqueueFailed {
		enqueueFailedError := errors.New("job queue full")
		cam.ServiceUnavailableError(w, r, enqueueFailedError)
	}
	return !enqueuedStatus
}

func (cam *CrmApiMux) writeJobAsResponse(newJob *server.Job, w http.ResponseWriter) {
	restResponse := new(server.RestResponse).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusCreated).
		WithCacheControlMaxAge(cam.CacheMaxAge()).
		WithJsonContent(newJob)
	cam.writeResponse(restResponse, "create job")
}

func (cam *CrmApiMux) DoJob(job *server.Job) {
	scenarioConfig, ok := job.HiddenAttributes[scenarioConfigStructKey].(*config.CRMConfig)
	if !ok {
		cam.Logger().Error("Unexpected error in retrieving CRMConfig from job hidden attributes")
		return
	}

	cam.Logger().Info("Running Job [" + string(job.Id) + "].")

	scenario.RunScenarioFromConfig(scenarioConfig)

	job.SetStatus(server.JobCompleted)

	job.RecordCompletionTime()

	delete(job.HiddenAttributes, scenarioConfigStructKey)
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
		WithJsonContent(cam.JobHistory)

	cam.writeResponse(restResponse, "get jobs")
}

func (cam *CrmApiMux) v1GetJob(w http.ResponseWriter, r *http.Request) {
	if cam.expectedMethodNotSupplied(http.MethodGet, w, r) {
		return
	}

	job := cam.deriveJobFromGetJobRequest(r)
	cam.writeJobAsResponse(job, w)
}

func (cam *CrmApiMux) deriveJobFromGetJobRequest(r *http.Request) *server.Job {
	desiredId := getJobIdFromRequestUrl(r)
	return cam.JobWithId(desiredId)
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

func scenarioOf(job *server.Job) interface{} {
	return job.HiddenAttributes[scenarioConfigTextKey]
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

func (cam *CrmApiMux) deriveJobFromGetJobScenarioRequest(r *http.Request) *server.Job {
	desiredId := getJobIdFromScenarioRequestUrl(r)
	return cam.JobWithId(desiredId)
}

func getJobIdFromScenarioRequestUrl(r *http.Request) string {
	splitURL := strings.Split(r.URL.Path, server.UrlPathSeparator)
	desiredId := splitURL[len(splitURL)-2]
	return desiredId
}
