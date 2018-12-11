// Copyright (c) 2018 Australian Rivers Institute.

package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/LindsayBradford/crem/cmd/cremengine/components"
	"github.com/LindsayBradford/crem/cmd/cremengine/components/api"
	configTesting "github.com/LindsayBradford/crem/internal/pkg/config/testing"
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	"github.com/LindsayBradford/crem/internal/pkg/server/test"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	. "github.com/onsi/gomega"
)

func TestDumbAnnealerIntegrationOneRun(t *testing.T) {
	context := configTesting.Context{
		Name:           "Single run of Dumb annealer",
		T:              t,
		ConfigFilePath: "testdata/DumbAnnealerTestConfig-OneRun.toml",
		Runner:         components.RunScenarioFromConfigFile,
	}

	context.VerifyScenarioRunViaConfigFileDoesNotPanic()
}

func TestDumbAnnealerIntegrationThreeRunsSequentially(t *testing.T) {
	context := configTesting.Context{
		Name:           "Three sequential runs of Dumb annealer",
		T:              t,
		ConfigFilePath: "testdata/DumbAnnealerTestConfig-ThreeRunsSequentially.toml",
		Runner:         components.RunScenarioFromConfigFile,
	}

	context.VerifyScenarioRunViaConfigFileDoesNotPanic()
}

func TestDumbAnnealerIntegrationThreeRunsConcurrently(t *testing.T) {
	context := configTesting.Context{
		Name:           "Three concurrent runs of Dumb annealer",
		T:              t,
		ConfigFilePath: "testdata/DumbAnnealerTestConfig-ThreeRunsConcurrently.toml",
		Runner:         components.RunScenarioFromConfigFile,
	}

	context.VerifyScenarioRunViaConfigFileDoesNotPanic()
}

func TestValidJobsGetRequest_OkResponse(t *testing.T) {
	context := configTesting.Context{
		Name:           "GET /jobs request returns 200 response",
		T:              t,
		ConfigFilePath: "testdata/server.toml",
	}

	verifyResponseToValidJobsGetRequest(context)
}

func verifyResponseToValidJobsGetRequest(context configTesting.Context) {

	g := NewGomegaWithT(context.T)

	muxUnderTest := buildMuxUnderTest()

	requestContext := test.HttpTestRequestContext{
		Method:    "GET",
		TargetUrl: "http://dummyUrl/api/v1/jobs",
		Handler:   muxUnderTest.V1HandleJobs,
	}

	responseContainer := requestContext.BuildJsonResponse()

	g.Expect(responseContainer.StatusCode).To(BeNumerically("==", http.StatusOK), context.Name+" should return OK status")

	// verifyResponseTimeIsAboutNow(g, responseContainer)
}

func TestInvalidValidJobsPostRequest_InternalServerErrorResponse(t *testing.T) {
	context := configTesting.Context{
		Name:           "POST /jobs request of invalid scenario returns 500 response",
		T:              t,
		ConfigFilePath: "testdata/server.toml",
	}

	verifyInternalServerErrorResponseToInvalidJobsPostRequest(context)
}

func verifyInternalServerErrorResponseToInvalidJobsPostRequest(context configTesting.Context) {

	g := NewGomegaWithT(context.T)

	muxUnderTest := buildMuxUnderTest()

	requestContext := test.HttpTestRequestContext{
		Method:      "POST",
		TargetUrl:   "http://dummyUrl/api/v1/jobs",
		ContentType: rest.TomlMimeType,
		RequestBody: "invalidScenarioText: isInvalid",
		Handler:     muxUnderTest.V1HandleJobs,
	}

	responseContainer := requestContext.BuildJsonResponse()

	g.Expect(responseContainer.StatusCode).To(BeNumerically("==", http.StatusInternalServerError), context.Name+" should return 500 status")

	verifyResponseTimeIsAboutNow(g, responseContainer)
}

func buildMuxUnderTest() *api.Mux {
	muxUnderTest := new(api.Mux).Initialise()
	muxUnderTest.SetLogger(loggers.DefaultTestingLogger)
	return muxUnderTest
}

func verifyResponseTimeIsAboutNow(g *GomegaWithT, responseContainer test.JsonResponseContainer) {
	responseTimeString, ok := responseContainer.JsonMap["Time"].(string)
	g.Expect(ok).To(Equal(true), " should return a string encoding of time")

	responseTime, parseErr := time.Parse(time.RFC3339Nano, responseTimeString)
	g.Expect(parseErr).To(BeNil(), " should return a RFC3339Nano encoded string of time")

	g.Expect(responseTime).To(BeTemporally("~", time.Now(), time.Millisecond*5), " should return status time of about now")
}

func TestKirkpatrickDumbAnnealerIntegrationOneRun(t *testing.T) {
	context := configTesting.Context{
		Name:           "Single run of Kirkpatrick Dumb annealer",
		T:              t,
		ConfigFilePath: "testdata/KirkpatrickDumbAnnealerTestConfig-OneRun.toml",
		Runner:         components.RunScenarioFromConfigFile,
	}

	context.VerifyScenarioRunViaConfigFileDoesNotPanic()
}
