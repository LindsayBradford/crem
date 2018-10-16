// Copyright (c) 2018 Australian Rivers Institute.

package server

import (
	"net/http"
	"testing"
	"time"

	"github.com/LindsayBradford/crm/logging/handlers"
	. "github.com/onsi/gomega"
)

type testContext struct {
	name       string
	t          *testing.T
	configFile string
}

func TestValidStatusRequest_OkResponse(t *testing.T) {
	context := testContext{
		name:       "GET /status request returns 200 response",
		t:          t,
		configFile: "testdata/server.toml",
	}

	verifyResponseToValidStatusRequest(context)
}

func TestInvalidStatusRequest_MethodNotAllowedResponse(t *testing.T) {
	context := testContext{
		name:       "POST /status request returns 405 response",
		t:          t,
		configFile: "testdata/server.toml",
	}

	verifyResponseToInvalidStatusRequest(context)
}

func verifyResponseToValidStatusRequest(context testContext) {
	g := NewGomegaWithT(context.t)

	muxUnderTest := buildMuxUnderTest()

	expectedMessage := "Some bogus status message"
	expectedName := "Some bogus name"
	expectedVersion := "some bogus version"
	bogusTime := "some bogus time"
	muxUnderTest.Status = ServiceStatus{ServiceName: expectedName, Version: expectedVersion, Status: expectedMessage, Time: bogusTime}

	requestContext := HttpTestRequestContext{
		Method:    "GET",
		TargetUrl: "http://dummyUrl/status",
		Handler:   muxUnderTest.statusHandler,
	}

	responseContainer := requestContext.BuildJsonResponse()

	g.Expect(responseContainer.StatusCode).To(BeNumerically("==", http.StatusOK), context.name+" should return OK status")
	g.Expect(responseContainer.JsonMap["ServiceName"]).To(Equal(expectedName), context.name+" should return expected status name")
	g.Expect(responseContainer.JsonMap["Version"]).To(Equal(expectedVersion), context.name+" should return expected status version")
	g.Expect(responseContainer.JsonMap["Status"]).To(Equal(expectedMessage), context.name+" should return expected status message")

	verifyResponseTimeIsAboutNow(g, responseContainer)
}

func verifyResponseToInvalidStatusRequest(context testContext) {
	g := NewGomegaWithT(context.t)

	muxUnderTest := buildMuxUnderTest()

	requestContext := HttpTestRequestContext{
		Method:    "POST",
		TargetUrl: "http://dummyUrl/status",
		Handler:   muxUnderTest.statusHandler,
	}

	responseContainer := requestContext.BuildJsonResponse()

	expectedResponseCode := http.StatusMethodNotAllowed
	g.Expect(responseContainer.StatusCode).To(BeNumerically("==", expectedResponseCode), context.name+" should return Method not Allowed status")
	expectedMessage := "HTTP Method not allowed"
	g.Expect(responseContainer.JsonMap["ErrorMessage"]).To(Equal(expectedMessage), context.name+" should return expected status message")

	verifyResponseTimeIsAboutNow(g, responseContainer)
}

func buildMuxUnderTest() *AdminMux {
	muxUnderTest := new(AdminMux).Initialise()
	muxUnderTest.SetLogger(handlers.DefaultTestingLogHandler)
	return muxUnderTest
}

func verifyResponseTimeIsAboutNow(g *GomegaWithT, responseContainer JsonResponseContainer) {
	responseTimeString, ok := responseContainer.JsonMap["Time"].(string)
	g.Expect(ok).To(Equal(true), " should return a string encoding of time")

	responseTime, parseErr := time.Parse(time.RFC3339Nano, responseTimeString)
	g.Expect(parseErr).To(BeNil(), " should return a RFC3339Nano encoded string of time")

	g.Expect(responseTime).To(BeTemporally("~", time.Now(), time.Millisecond*5), " should return status time of about now")
}

func TestValidShutdownRequest_OkResponse(t *testing.T) {
	context := testContext{
		name:       "POST /shutdown request returns 200 response",
		t:          t,
		configFile: "testdata/server.toml",
	}

	verifyResponseToValidShutdownRequest(context)
}

func verifyResponseToValidShutdownRequest(context testContext) {
	g := NewGomegaWithT(context.t)

	muxUnderTest := buildMuxUnderTest()

	requestContext := HttpTestRequestContext{
		Method:    "POST",
		TargetUrl: "http://dummyUrl/shutdown",
		Handler:   muxUnderTest.shutdownHandler,
	}

	var responseContainer JsonResponseContainer
	go func() {
		responseContainer = requestContext.BuildJsonResponse()
	}()

	waitFinished := true
	waitFunc := func() bool {
		muxUnderTest.WaitForShutdownSignal()
		return true
	}()

	g.Eventually(waitFunc).Should(Equal(waitFinished), context.name+" should finish Wait function")
	g.Expect(responseContainer.StatusCode).To(BeNumerically("==", http.StatusOK), context.name+" should return OK status")

	verifyResponseTimeIsAboutNow(g, responseContainer)
}

func TestInvalidShutdownRequest_MethodNotAllowedResponse(t *testing.T) {
	context := testContext{
		name:       "GET /shutdown request returns 405 response",
		t:          t,
		configFile: "testdata/server.toml",
	}

	verifyResponseToInvalidShutdownRequest(context)
}

func verifyResponseToInvalidShutdownRequest(context testContext) {
	g := NewGomegaWithT(context.t)

	muxUnderTest := buildMuxUnderTest()

	requestContext := HttpTestRequestContext{
		Method:    "GET",
		TargetUrl: "http://dummyUrl/shutdown",
		Handler:   muxUnderTest.shutdownHandler,
	}

	responseContainer := requestContext.BuildJsonResponse()

	expectedResponseCode := http.StatusMethodNotAllowed
	g.Expect(responseContainer.StatusCode).To(BeNumerically("==", expectedResponseCode), context.name+" should return Method not Allowed status")
	expectedMessage := "HTTP Method not allowed"
	g.Expect(responseContainer.JsonMap["ErrorMessage"]).To(Equal(expectedMessage), context.name+" should return expected status message")

	verifyResponseTimeIsAboutNow(g, responseContainer)
}
