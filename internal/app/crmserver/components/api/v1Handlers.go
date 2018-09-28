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

func (cam *CrmApiMux) v1PostScenario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		cam.ServeMethodNotAllowedError(w, r)
		return
	}

	if r.Header.Get("Content-type") != server.TomlMimeType {
		cam.ServeMethodNotAllowedError(w, r)
		return
	}

	scenarioText := requestBodyToString(r)

	scenarioConfig, retrieveError := config.RetrieveCrmFromString(scenarioText)

	if retrieveError != nil {
		wrappingError := errors.Wrap(retrieveError, "retrieving scenario configuration")
		cam.Logger().Warn(wrappingError)
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
