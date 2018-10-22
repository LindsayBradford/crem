// Copyright (c) 2018 Australian Rivers Institute.

package api

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/LindsayBradford/crem/config"
	"github.com/LindsayBradford/crem/internal/app/cremserver/components/scenario"
	"github.com/LindsayBradford/crem/server"
	"github.com/pkg/errors"
)

const scenarioConfigTextKey = "ScenarioConfigText"
const scenarioConfigStructKey = "ScenarioConfigStruct"

const scenarioPath = "scenario"

func (cam *CremApiMux) V1HandleJobs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		cam.v1PostJob(w, r)
	case http.MethodGet:
		cam.v1GetJobs(w, r)
	case http.MethodDelete:
		cam.v1DeleteJobs(w, r)
	default:
		cam.MethodNotAllowedError(w, r)
	}
}

func (cam *CremApiMux) v1PostJob(w http.ResponseWriter, r *http.Request) {
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

func (cam *CremApiMux) createJobFromRequest(r *http.Request) *server.Job {
	newJob := new(server.Job).Initialise()
	newJob.HiddenAttributes[scenarioConfigTextKey] = requestBodyToString(r)
	return newJob
}

func (cam *CremApiMux) requestContentTypeWasNotToml(r *http.Request, w http.ResponseWriter) bool {
	if r.Header.Get(server.ContentTypeHeaderKey) != server.TomlMimeType {
		cam.MethodNotAllowedError(w, r)
		return true
	}
	return false
}

func (cam *CremApiMux) scenarioTextSuppliedIsInvalid(job *server.Job, w http.ResponseWriter, r *http.Request) bool {
	scenarioText := job.HiddenAttributes[scenarioConfigTextKey].(string)
	scenarioConfig, retrieveError := config.RetrieveCremFromString(scenarioText)

	if retrieveError != nil {
		cam.handleScenarioConfigRetrieveError(retrieveError, job, w, r)
		return true
	}

	cam.Logger().Info("Received new job [" + string(job.Id) + "] with scenario name [" + scenarioConfig.ScenarioName + "].")
	job.HiddenAttributes[scenarioConfigStructKey] = scenarioConfig
	return false
}

func (cam *CremApiMux) handleScenarioConfigRetrieveError(retrieveError error, job *server.Job, w http.ResponseWriter, r *http.Request) {
	wrappingError := errors.Wrap(retrieveError, "retrieving scenario configuration")
	cam.Logger().Warn(wrappingError)
	job.SetStatus(server.JobInvalid)
	cam.AddToHistory(job)
	cam.InternalServerError(w, r, errors.New("Invalid scenario configuration supplied"))
	cam.allowJobQueries(job)
}

func (cam *CremApiMux) allowJobQueries(job *server.Job) {
	cam.AddHandler(BuildApiPath(jobsPath, string(job.Id)), cam.v1GetJob)
	cam.AddHandler(BuildApiPath(jobsPath, string(job.Id), scenarioPath), cam.v1GetJobScenario)
}

func (cam *CremApiMux) jobEnqueueFailed(job *server.Job, w http.ResponseWriter, r *http.Request) server.JobEnqueuedStatus {
	cam.AddToHistory(job)
	enqueuedStatus := cam.jobs.Enqueue(job)
	if enqueuedStatus == server.JobEnqueueFailed {
		enqueueFailedError := errors.New("job queue full")
		cam.ServiceUnavailableError(w, r, enqueueFailedError)
	}
	return !enqueuedStatus
}

func (cam *CremApiMux) writeJobAsResponse(newJob *server.Job, w http.ResponseWriter) {
	restResponse := new(server.RestResponse).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusCreated).
		WithCacheControlMaxAge(cam.CacheMaxAge()).
		WithJsonContent(newJob)
	cam.writeResponse(restResponse, "create job")
}

func (cam *CremApiMux) DoJob(job *server.Job) {
	scenarioConfig, ok := job.HiddenAttributes[scenarioConfigStructKey].(*config.CREMConfig)
	if !ok {
		cam.Logger().Error("Unexpected error in retrieving CREMConfig from job hidden attributes")
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

func (cam *CremApiMux) v1GetJobs(w http.ResponseWriter, r *http.Request) {
	restResponse := new(server.RestResponse).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(cam.CacheMaxAge()).
		WithJsonContent(cam.JobHistory)

	cam.writeResponse(restResponse, "get jobs")
}

func (cam *CremApiMux) v1GetJob(w http.ResponseWriter, r *http.Request) {
	if cam.expectedMethodNotSupplied(http.MethodGet, w, r) {
		return
	}

	job := cam.deriveJobFromGetJobRequest(r)
	cam.writeJobAsResponse(job, w)
}

func (cam *CremApiMux) deriveJobFromGetJobRequest(r *http.Request) *server.Job {
	desiredId := getJobIdFromRequestUrl(r)
	return cam.JobWithId(desiredId)
}

func getJobIdFromRequestUrl(r *http.Request) string {
	splitURL := strings.Split(r.URL.Path, server.UrlPathSeparator)
	desiredId := splitURL[len(splitURL)-1]
	return desiredId
}

func (cam *CremApiMux) v1GetJobScenario(w http.ResponseWriter, r *http.Request) {
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

func (cam *CremApiMux) expectedMethodNotSupplied(expectedHttpMethod string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != expectedHttpMethod {
		cam.MethodNotAllowedError(w, r)
		return true
	}
	return false
}

func (cam *CremApiMux) writeResponse(response *server.RestResponse, context string) {
	writeError := response.Write()
	if writeError != nil {
		wrappingError := errors.Wrap(writeError, context)
		cam.Logger().Error(wrappingError)
	}
}

func (cam *CremApiMux) deriveJobFromGetJobScenarioRequest(r *http.Request) *server.Job {
	desiredId := getJobIdFromScenarioRequestUrl(r)
	return cam.JobWithId(desiredId)
}

func getJobIdFromScenarioRequestUrl(r *http.Request) string {
	splitURL := strings.Split(r.URL.Path, server.UrlPathSeparator)
	desiredId := splitURL[len(splitURL)-2]
	return desiredId
}

func (cam *CremApiMux) v1DeleteJobs(w http.ResponseWriter, r *http.Request) {
	deletedJobs := cam.deleteProcessedJobs()

	cam.Logger().Info("Deleting all processed jobs.")

	restResponse := new(server.RestResponse).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(cam.CacheMaxAge()).
		WithJsonContent(deletedJobs)

	cam.writeResponse(restResponse, "delete jobs")
}

func (cam *CremApiMux) deleteProcessedJobs() []*server.Job {
	processedJobs := make([]*server.Job, 0)
	unprocessedJobs := make([]*server.Job, 0)

	for _, job := range cam.JobHistory {
		if isProcessed(job) {
			processedJobs = append(processedJobs, job)
		} else {
			unprocessedJobs = append(unprocessedJobs, job)
		}
	}

	cam.JobHistory = unprocessedJobs
	return processedJobs
}

func isProcessed(job *server.Job) bool {
	if job.Status() == server.JobCompleted || job.Status() == server.JobInvalid {
		return true
	}
	return false
}
