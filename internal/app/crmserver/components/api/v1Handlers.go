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
	newJob.HiddenAttributes["ScenarioConfig"] = scenarioText

	cam.jobs.Enqueue(*newJob)
	cam.AddHandler(baseApiPath()+jobsPath+"/"+newJob.JobId, cam.v1GetJob)
	cam.AddHandler(baseApiPath()+jobsPath+"/"+newJob.JobId+"/scenario", cam.v1GetJobScenario)

	scenarioConfig, retrieveError := config.RetrieveCrmFromString(scenarioText)

	if retrieveError != nil {
		wrappingError := errors.Wrap(retrieveError, "retrieving scenario configuration")
		cam.Logger().Warn(wrappingError)
		cam.InternalServerError(w, r, errors.New("Invalid scenario configuration supplied"))
		return
	}

	cam.Logger().Info("New Job [" + newJob.JobId + "] received with scenario name [" + scenarioConfig.ScenarioName + "].")

	// scenarioJson, _ := json.Marshal(scenarioConfig)
	// cam.Logger().Debug("Job [" + newJob.JobId + "] scenario config: ["+ string(scenarioJson) + "].")

	restResponse := new(server.RestResponse).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusCreated).
		WithCacheControlMaxAge(cam.CacheMaxAge()).
		WithJsonContent(newJob)

	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "create job")
		cam.Logger().Error(wrappingError)
	}

	scenario.RunScenarioFromConfig(scenarioConfig)

	newJob.Attributes["Status"] = "COMPLETED"
	newJob.Attributes["CompletedTime"] = server.FormattedTimestamp()
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

	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "get jobs")
		cam.Logger().Error(wrappingError)
	}
}

func (cam *CrmApiMux) v1GetJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		cam.MethodNotAllowedError(w, r)
		return
	}

	desiredId := getJobIdFromRequestUrl(r)
	jobToReturn := cam.jobs.JobWithId(desiredId)

	restResponse := new(server.RestResponse).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlMaxAge(cam.CacheMaxAge()).
		WithJsonContent(jobToReturn)

	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "get job")
		cam.Logger().Error(wrappingError)
	}
}

func getJobIdFromRequestUrl(r *http.Request) string {
	splitURL := strings.Split(r.URL.Path, server.UrlPathSeparator)
	desiredId := splitURL[len(splitURL)-1]
	return desiredId
}

func (cam *CrmApiMux) v1GetJobScenario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		cam.MethodNotAllowedError(w, r)
		return
	}

	desiredId := getJobIdFromScenarioRequestUrl(r)
	desiredJob := cam.jobs.JobWithId(desiredId)

	restResponse := new(server.RestResponse).
		Initialise().
		WithWriter(w).
		WithResponseCode(http.StatusOK).
		WithCacheControlPublic().
		WithTomlContent(desiredJob.HiddenAttributes["ScenarioConfig"])

	writeError := restResponse.Write()

	if writeError != nil {
		wrappingError := errors.Wrap(writeError, "get job scenario")
		cam.Logger().Error(wrappingError)
	}
}

func getJobIdFromScenarioRequestUrl(r *http.Request) string {
	splitURL := strings.Split(r.URL.Path, server.UrlPathSeparator)
	desiredId := splitURL[len(splitURL)-2]
	return desiredId
}
