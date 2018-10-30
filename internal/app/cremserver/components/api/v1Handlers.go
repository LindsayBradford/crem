// Copyright (c) 2018 Australian Rivers Institute.

package api

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/LindsayBradford/crem/config"
	"github.com/LindsayBradford/crem/internal/app/cremserver/components/scenario"
	"github.com/LindsayBradford/crem/server"
	"github.com/LindsayBradford/crem/server/job"
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

func (cam *CremApiMux) createJobFromRequest(r *http.Request) *job.Job {
	newJob := new(job.Job).Initialise()
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

func (cam *CremApiMux) scenarioTextSuppliedIsInvalid(job *job.Job, w http.ResponseWriter, r *http.Request) bool {
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

func (cam *CremApiMux) handleScenarioConfigRetrieveError(retrieveError error, jobInError *job.Job, w http.ResponseWriter, r *http.Request) {
	wrappingError := errors.Wrap(retrieveError, "retrieving scenario configuration")
	cam.Logger().Warn(wrappingError)
	jobInError.SetStatus(job.Invalid)
	cam.AddToHistory(jobInError)
	cam.InternalServerError(w, r, errors.New("Invalid scenario configuration supplied"))
	cam.allowJobQueries(jobInError)
}

func (cam *CremApiMux) allowJobQueries(job *job.Job) {
	cam.AddHandler(BuildApiPath(jobsPath, string(job.Id)), cam.v1GetJob)
	cam.AddHandler(BuildApiPath(jobsPath, string(job.Id), scenarioPath), cam.v1GetJobScenario)
}

func (cam *CremApiMux) jobEnqueueFailed(jobToEnqueue *job.Job, w http.ResponseWriter, r *http.Request) job.EnqueuedStatus {
	cam.AddToHistory(jobToEnqueue)
	enqueuedStatus := cam.jobs.Enqueue(jobToEnqueue)
	if enqueuedStatus == job.EnqueueFailed {
		enqueueFailedError := errors.New("jobToEnqueue queue full")
		cam.ServiceUnavailableError(w, r, enqueueFailedError)
	}
	return !enqueuedStatus
}

func (cam *CremApiMux) writeJobAsResponse(newJob *job.Job, w http.ResponseWriter) {
	restResponse := new(server.RestResponse).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusCreated).
		WithCacheControlMaxAge(cam.CacheMaxAge()).
		WithJsonContent(newJob)
	cam.writeResponse(restResponse, "create job")
}

func (cam *CremApiMux) DoJob(jobToDo *job.Job) {
	scenarioConfig, ok := jobToDo.HiddenAttributes[scenarioConfigStructKey].(*config.CREMConfig)
	if !ok {
		cam.Logger().Error("Unexpected error in retrieving CREMConfig from jobToDo hidden attributes")
		return
	}

	cam.Logger().Info("Running Job [" + string(jobToDo.Id) + "].")

	scenario.RunScenarioFromConfig(scenarioConfig)

	jobToDo.SetStatus(job.Completed)

	jobToDo.RecordCompletionTime()

	delete(jobToDo.HiddenAttributes, scenarioConfigStructKey)
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

func (cam *CremApiMux) deriveJobFromGetJobRequest(r *http.Request) *job.Job {
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

func scenarioOf(job *job.Job) interface{} {
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

func (cam *CremApiMux) deriveJobFromGetJobScenarioRequest(r *http.Request) *job.Job {
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

func (cam *CremApiMux) deleteProcessedJobs() []*job.Job {
	processedJobs := make([]*job.Job, 0)
	unprocessedJobs := make([]*job.Job, 0)

	for _, job := range cam.JobHistory {
		if job.IsProcessed() {
			processedJobs = append(processedJobs, job)
		} else {
			unprocessedJobs = append(unprocessedJobs, job)
		}
	}

	cam.JobHistory = unprocessedJobs
	return processedJobs
}
