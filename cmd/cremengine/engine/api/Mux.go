// Copyright (c) 2018 Australian Rivers Institute.

package api

import (
	"io/ioutil"
	"net/http"

	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution/encoding/json"
	"github.com/LindsayBradford/crem/internal/pkg/config/interpreter"
	"github.com/LindsayBradford/crem/internal/pkg/dataset"
	"github.com/LindsayBradford/crem/internal/pkg/model/models/catchment"
	serverApi "github.com/LindsayBradford/crem/internal/pkg/server/api"
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	"github.com/LindsayBradford/crem/pkg/attributes"
	"github.com/LindsayBradford/crem/pkg/threading"
)

const (
	v1Path = "v1"

	scenarioTextKey = "ScenarioText"
	scenarioNameKey = "ScenarioName"

	solutionsTextKey = "SolutionsText"
)

type Mux struct {
	serverApi.Mux
	mainThreadChannel *threading.MainThreadChannel

	modelConfigInterpreter *interpreter.ModelConfigInterpreter
	model                  *catchment.Model
	modelSolution          *solution.Solution

	solutionPool     SolutionPool
	solutionSetTable dataset.HeadingsTable

	jsonMarshaler json.Marshaler

	attributes.ContainedAttributes
}

func (m *Mux) Initialise() *Mux {
	const (
		scenarioPath         = "scenario"
		solutionsPath        = "solutions"
		modelPath            = "model"
		actionsPath          = "actions"
		activePath           = "active"
		applicablePath       = "applicable"
		subcatchmentPath     = "subcatchment"
		identityMatchingPath = "\\d+"
		solutionLabelPath    = "[\\w\\-]+"
	)

	m.Mux.Initialise()

	m.modelConfigInterpreter = interpreter.NewModelConfigInterpreter()

	m.AddHandler(buildV1ApiPath(scenarioPath), m.v1scenarioHandler)
	m.AddHandler(buildV1ApiPath(solutionsPath), m.v1solutionSetHandler)
	m.AddHandler(buildV1ApiPath(solutionsPath, solutionLabelPath), m.v1solutionHandler)
	m.AddHandler(buildV1ApiPath(modelPath), m.v1modelHandler)
	m.AddHandler(buildV1ApiPath(modelPath, actionsPath, applicablePath), m.v1ApplicableActionsHandler)
	m.AddHandler(buildV1ApiPath(modelPath, actionsPath, activePath), m.v1activeActionsHandler)
	m.AddHandler(buildV1ApiPath(modelPath, subcatchmentPath, identityMatchingPath), m.v1subcatchmentHandler)

	return m
}

func (m *Mux) WithMainThreadChannel(channel *threading.MainThreadChannel) *Mux {
	m.mainThreadChannel = channel
	return m
}

func (m *Mux) WithCacheMaxAge(maxAgeInSeconds uint64) *Mux {
	m.MuxImpl.WithCacheMaxAge(maxAgeInSeconds)
	return m
}

func buildV1ApiPath(pathElements ...string) string {
	const (
		startPathMatcher = "^"
		endPathMatcher   = "$"
	)

	builtPath := rest.UrlPathSeparator + serverApi.BasePath + rest.UrlPathSeparator + v1Path

	for _, element := range pathElements {
		builtPath = builtPath + rest.UrlPathSeparator + element
	}

	return startPathMatcher + builtPath + endPathMatcher
}

func (m *Mux) AddHandler(address string, handler rest.HandlerFunc) {
	m.HandlerMap.AddHandler(address, handler)
}

func (m *Mux) Shutdown() {
	if m.model != nil {
		m.model.TearDown()
	}
	m.MuxImpl.Shutdown()
}

func requestBodyToBytes(r *http.Request) []byte {
	responseBodyBytes, _ := ioutil.ReadAll(r.Body)
	return responseBodyBytes
}

func requestBodyToString(r *http.Request) string {
	return string(requestBodyToBytes(r))
}
