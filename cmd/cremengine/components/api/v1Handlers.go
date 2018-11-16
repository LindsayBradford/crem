// Copyright (c) 2018 Australian Rivers Institute.

package api

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/LindsayBradford/crem/cmd/cremengine/components/scenario"
	"github.com/LindsayBradford/crem/internal/pkg/config"
	"github.com/LindsayBradford/crem/internal/pkg/server/job"
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	"github.com/pkg/errors"
)

const scenarioConfigTextKey = "ScenarioConfigText"
const scenarioConfigStructKey = "ScenarioConfigStruct"

const scenarioPath = "scenario"

func (m *Mux) V1HandleJobs(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		m.v1PostJob(w, r)
	case http.MethodGet:
		m.v1GetJobs(w, r)
	case http.MethodDelete:
		m.v1DeleteJobs(w, r)
	default:
		m.MethodNotAllowedError(w, r)
	}
}

func (m *Mux) v1PostJob(w http.ResponseWriter, r *http.Request) {
	if m.requestContentTypeWasNotToml(r, w) {
		return
	}

	newJob := m.createJobFromRequest(r)

	if m.scenarioTextSuppliedIsInvalid(newJob, w, r) || bool(m.jobEnqueueFailed(newJob, w, r)) {
		return
	}

	m.writeJobAsResponse(newJob, w)
	m.allowJobQueries(newJob)
}

func (m *Mux) createJobFromRequest(r *http.Request) *job.Job {
	newJob := new(job.Job).Initialise()
	newJob.HiddenAttributes[scenarioConfigTextKey] = requestBodyToString(r)
	return newJob
}

func (m *Mux) requestContentTypeWasNotToml(r *http.Request, w http.ResponseWriter) bool {
	if r.Header.Get(rest.ContentTypeHeaderKey) != rest.TomlMimeType {
		m.MethodNotAllowedError(w, r)
		return true
	}
	return false
}

func (m *Mux) scenarioTextSuppliedIsInvalid(job *job.Job, w http.ResponseWriter, r *http.Request) bool {
	scenarioText := job.HiddenAttributes[scenarioConfigTextKey].(string)
	scenarioConfig, retrieveError := config.RetrieveCremFromString(scenarioText)

	if retrieveError != nil {
		m.handleScenarioConfigRetrieveError(retrieveError, job, w, r)
		return true
	}

	m.Logger().Info("Received new job [" + string(job.Id) + "] with scenario name [" + scenarioConfig.ScenarioName + "].")
	job.HiddenAttributes[scenarioConfigStructKey] = scenarioConfig
	return false
}

func (m *Mux) handleScenarioConfigRetrieveError(retrieveError error, jobInError *job.Job, w http.ResponseWriter, r *http.Request) {
	wrappingError := errors.Wrap(retrieveError, "retrieving scenario configuration")
	m.Logger().Warn(wrappingError)
	jobInError.SetStatus(job.Invalid)
	jobInError.Attributes["InvalidationTime"] = rest.FormattedTimestamp()
	jobInError.Attributes["InvalidationDetail"] = wrappingError.Error()
	m.AddToHistory(jobInError)
	m.InternalServerError(w, r, errors.New("Invalid scenario configuration supplied"))
	m.allowJobQueries(jobInError)
}

func (m *Mux) allowJobQueries(job *job.Job) {
	m.AddHandler(BuildApiPath(jobsPath, string(job.Id)), m.v1GetJob)
	m.AddHandler(BuildApiPath(jobsPath, string(job.Id), scenarioPath), m.v1GetJobScenario)
}

func (m *Mux) jobEnqueueFailed(jobToEnqueue *job.Job, w http.ResponseWriter, r *http.Request) job.EnqueuedStatus {
	m.AddToHistory(jobToEnqueue)
	enqueuedStatus := m.jobs.Enqueue(jobToEnqueue)
	if enqueuedStatus == job.EnqueueFailed {
		enqueueFailedError := errors.New("jobToEnqueue queue full")
		m.ServiceUnavailableError(w, r, enqueueFailedError)
	}
	return !enqueuedStatus
}

func (m *Mux) writeJobAsResponse(newJob *job.Job, w http.ResponseWriter) {
	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusCreated).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithJsonContent(newJob)
	m.writeResponse(restResponse, "create job")
}

func (m *Mux) DoJob(jobToDo *job.Job) {
	scenarioConfig, ok := jobToDo.HiddenAttributes[scenarioConfigStructKey].(*config.CREMConfig)
	if !ok {
		m.Logger().Error("Unexpected error in retrieving CREMConfig from jobToDo hidden attributes")
		return
	}

	m.Logger().Info("Running Job [" + string(jobToDo.Id) + "].")

	runError := scenario.RunScenarioFromConfig(scenarioConfig)

	if runError != nil {
		jobToDo.Attributes["ErrorDetail"] = runError
		jobToDo.Attributes["ErrorTime"] = rest.FormattedTimestamp()
		jobToDo.SetStatus(job.Errored)
		return
	}

	jobToDo.SetStatus(job.Completed)

	jobToDo.RecordCompletionTime()

	delete(jobToDo.HiddenAttributes, scenarioConfigStructKey)
}

func requestBodyToString(r *http.Request) string {
	responseBodyBytes, _ := ioutil.ReadAll(r.Body)
	return string(responseBodyBytes)
}

func (m *Mux) v1GetJobs(w http.ResponseWriter, r *http.Request) {
	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithJsonContent(m.JobHistory)

	m.writeResponse(restResponse, "get jobs")
}

func (m *Mux) v1GetJob(w http.ResponseWriter, r *http.Request) {
	if m.expectedMethodNotSupplied(http.MethodGet, w, r) {
		return
	}

	job := m.deriveJobFromGetJobRequest(r)
	m.writeJobAsResponse(job, w)
}

func (m *Mux) deriveJobFromGetJobRequest(r *http.Request) *job.Job {
	desiredId := getJobIdFromRequestUrl(r)
	return m.JobWithId(desiredId)
}

func getJobIdFromRequestUrl(r *http.Request) string {
	splitURL := strings.Split(r.URL.Path, rest.UrlPathSeparator)
	desiredId := splitURL[len(splitURL)-1]
	return desiredId
}

func (m *Mux) v1GetJobScenario(w http.ResponseWriter, r *http.Request) {
	if m.expectedMethodNotSupplied(http.MethodGet, w, r) {
		return
	}

	job := m.deriveJobFromGetJobScenarioRequest(r)

	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlPublic().
		WithTomlContent(scenarioOf(job))

	m.writeResponse(restResponse, "get job scenario")
}

func scenarioOf(job *job.Job) interface{} {
	return job.HiddenAttributes[scenarioConfigTextKey]
}

func (m *Mux) expectedMethodNotSupplied(expectedHttpMethod string, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != expectedHttpMethod {
		m.MethodNotAllowedError(w, r)
		return true
	}
	return false
}

func (m *Mux) writeResponse(response *rest.Response, context string) {
	writeError := response.Write()
	if writeError != nil {
		wrappingError := errors.Wrap(writeError, context)
		m.Logger().Error(wrappingError)
	}
}

func (m *Mux) deriveJobFromGetJobScenarioRequest(r *http.Request) *job.Job {
	desiredId := getJobIdFromScenarioRequestUrl(r)
	return m.JobWithId(desiredId)
}

func getJobIdFromScenarioRequestUrl(r *http.Request) string {
	splitURL := strings.Split(r.URL.Path, rest.UrlPathSeparator)
	desiredId := splitURL[len(splitURL)-2]
	return desiredId
}

func (m *Mux) v1DeleteJobs(w http.ResponseWriter, r *http.Request) {
	deletedJobs := m.deleteProcessedJobs()

	m.Logger().Info("Deleting all processed jobs.")

	restResponse := new(rest.Response).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(m.CacheMaxAge()).
		WithJsonContent(deletedJobs)

	m.writeResponse(restResponse, "delete jobs")
}

func (m *Mux) deleteProcessedJobs() []*job.Job {
	processedJobs := make([]*job.Job, 0)
	unprocessedJobs := make([]*job.Job, 0)

	for _, job := range m.JobHistory {
		if job.IsProcessed() {
			processedJobs = append(processedJobs, job)
		} else {
			unprocessedJobs = append(unprocessedJobs, job)
		}
	}

	m.JobHistory = unprocessedJobs
	return processedJobs
}
