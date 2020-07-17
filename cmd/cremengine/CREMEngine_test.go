// Copyright (c) 2019 Australian Rivers Institute.

package main

import (
	"github.com/LindsayBradford/crem/internal/pkg/server/admin"
	"net/http"
	"testing"
	"time"

	webTesting "github.com/LindsayBradford/crem/internal/pkg/config/testing"
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	"github.com/LindsayBradford/crem/internal/pkg/server/test"
	"github.com/LindsayBradford/crem/pkg/logging/loggers"
	"github.com/onsi/gomega"
)

func TestValidStatusGetRequest_OkResponse(t *testing.T) {
	context := webTesting.WhiteboxTestingContext{
		Name:           "GET /status request returns 200 response",
		T:              t,
		ConfigFilePath: "testdata/TestEngine.toml",
	}

	verifyResponseToValidStatusGetRequest(context)
}

func verifyResponseToValidStatusGetRequest(context webTesting.WhiteboxTestingContext) {

	g := gomega.NewGomegaWithT(context.T)

	muxUnderTest := buildAdminMuxUnderTest()

	requestContext := test.HttpTestRequestContext{
		Method:    "GET",
		TargetUrl: "http://dummyUrl/status",
		Handler:   muxUnderTest.ServeHTTP,
	}

	responseContainer := requestContext.BuildJsonResponse()

	g.Expect(responseContainer.StatusCode).To(gomega.BeNumerically("==", http.StatusOK), context.Name+" should return OK status")
}

func TestInvalidStatusPostRequest_InternalServerErrorResponse(t *testing.T) {
	context := webTesting.WhiteboxTestingContext{
		Name:           "POST /admin/status request of invalid scenario returns 500 response",
		T:              t,
		ConfigFilePath: "testdata/TestEngine.toml",
	}

	verifyInternalServerErrorResponseToInvalidJobsPostRequest(context)
}

func verifyInternalServerErrorResponseToInvalidJobsPostRequest(context webTesting.WhiteboxTestingContext) {

	g := gomega.NewGomegaWithT(context.T)

	muxUnderTest := buildAdminMuxUnderTest()

	requestContext := test.HttpTestRequestContext{
		Method:      "POST",
		TargetUrl:   "http://dummyUrl/status",
		ContentType: rest.TomlMimeType,
		RequestBody: "invalidScenarioText: isInvalid",
		Handler:     muxUnderTest.ServeHTTP,
	}

	responseContainer := requestContext.BuildJsonResponse()

	g.Expect(responseContainer.StatusCode).To(gomega.BeNumerically("==", http.StatusMethodNotAllowed), context.Name+" should return 405 status")

	verifyResponseTimeIsAboutNow(g, responseContainer)
}

func testValidShutdownPostRequest_OkResponse(t *testing.T) {
	// TODO: the guts of this triggers a goroutine completion with successful processing, which as-is, stops subsequent tests from running.
	// turning off for now until I can better consider how to wrap the goroutine without lots of work.
	context := webTesting.WhiteboxTestingContext{
		Name:           "POST /shutdown request returns 200 response",
		T:              t,
		ConfigFilePath: "testdata/TestEngine.toml",
	}

	verifyResponseToValidShutdownPostRequest(context)
}

func verifyResponseToValidShutdownPostRequest(context webTesting.WhiteboxTestingContext) {

	g := gomega.NewGomegaWithT(context.T)

	muxUnderTest := buildAdminMuxUnderTest()

	requestContext := test.HttpTestRequestContext{
		Method:    "POST",
		TargetUrl: "http://dummyUrl/shutdown",
		Handler:   muxUnderTest.ServeHTTP,
	}

	responseContainer := requestContext.BuildJsonResponse()

	g.Expect(responseContainer.StatusCode).To(gomega.BeNumerically("==", http.StatusOK), context.Name+" should return OK status")
}

func TestInvalidShutdownGetRequest_InternalServerErrorResponse(t *testing.T) {
	context := webTesting.WhiteboxTestingContext{
		Name:           "POST /admin/status request of invalid scenario returns 500 response",
		T:              t,
		ConfigFilePath: "testdata/TestEngine.toml",
	}

	verifyInternalServerErrorResponseToInvalidJobsGetRequest(context)
}

func verifyInternalServerErrorResponseToInvalidJobsGetRequest(context webTesting.WhiteboxTestingContext) {

	g := gomega.NewGomegaWithT(context.T)

	muxUnderTest := buildAdminMuxUnderTest()

	requestContext := test.HttpTestRequestContext{
		Method:      "GET",
		TargetUrl:   "http://dummyUrl/shutdown",
		ContentType: rest.TomlMimeType,
		RequestBody: "invalidScenarioText: isInvalid",
		Handler:     muxUnderTest.ServeHTTP,
	}

	responseContainer := requestContext.BuildJsonResponse()

	g.Expect(responseContainer.StatusCode).To(gomega.BeNumerically("==", http.StatusMethodNotAllowed), context.Name+" should return 405 status")

	verifyResponseTimeIsAboutNow(g, responseContainer)
}

func buildAdminMuxUnderTest() *admin.Mux {
	muxUnderTest := new(admin.Mux).Initialise()
	muxUnderTest.Status = admin.ServiceStatus{
		ServiceName: "test admin service",
		Version:     "0.1",
		Status:      "TESTING"}
	muxUnderTest.SetLogger(loggers.DefaultTestingLogger)
	return muxUnderTest
}

func verifyResponseTimeIsAboutNow(g *gomega.GomegaWithT, responseContainer test.JsonResponseContainer) {
	responseTimeString, ok := responseContainer.JsonMap["Time"].(string)
	g.Expect(ok).To(gomega.Equal(true), " should return a string encoding of time")

	responseTime, parseErr := time.Parse(time.RFC3339Nano, responseTimeString)
	g.Expect(parseErr).To(gomega.BeNil(), " should return a RFC3339Nano encoded string of time")

	g.Expect(responseTime).To(gomega.BeTemporally("~", time.Now(), time.Millisecond*5), " should return status time of about now")
}
