// Copyright (c) 2018 Australian Rivers Institute.

package api

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/LindsayBradford/crm/config"
	"github.com/LindsayBradford/crm/internal/app/crmserver/components/scenario"
	"github.com/LindsayBradford/crm/server"
	"github.com/pkg/errors"
)

type Job struct {
	ScenarioConfig *config.CRMConfig
	status         string
}

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

	scenarioConfig, retrieveError := config.RetrieveCrmFromString(scenarioText)

	if retrieveError != nil {
		wrappingError := errors.Wrap(retrieveError, "retrieving scenario configuration")
		cam.Logger().Warn(wrappingError)
		cam.InternalServerError(w, r, errors.New("Invalid scenario configuration supplied"))
		return
	}

	sendTextOnResponseBody(scenarioText, w)

	scenario.RunScenarioFromConfig(scenarioConfig)
}

func sendTextOnResponseBody(text string, w http.ResponseWriter) {
	fmt.Fprintf(w, text)
}

func requestBodyToString(r *http.Request) string {
	responseBodyBytes, _ := ioutil.ReadAll(r.Body)
	return string(responseBodyBytes)
}

func (cam *CrmApiMux) v1GetJobs(w http.ResponseWriter, r *http.Request) {
	scenarioText := requestBodyToString(r)
	sendTextOnResponseBody(scenarioText, w)
}
