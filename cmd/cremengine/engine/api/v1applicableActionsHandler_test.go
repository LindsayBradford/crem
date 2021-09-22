package api

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	httptest "github.com/LindsayBradford/crem/internal/pkg/server/test"
	. "github.com/onsi/gomega"
)

const applicableActionsPath = "applicableActions"

func TestFirstApplicableActionsGetRequest_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	subcatchmentUrlUnderTest := baseSubcatchmentUrl + rest.UrlPathSeparator + validSubcatchment + rest.UrlPathSeparator + applicableActionsPath

	// when
	context := TestContext{
		Name: http.MethodGet + " " + subcatchmentUrlUnderTest + " request returns 404 (not found) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:    http.MethodGet,
			TargetUrl: subcatchmentUrlUnderTest,
		},
		ExpectedResponseStatus: http.StatusNotFound,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)
	muxUnderTest.Shutdown()
}

func TestMissingApplicableActionsGetRequest_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	buildValidScenario(t, muxUnderTest)

	// when
	subcatchmentUrlUnderTest := baseSubcatchmentUrl + "/1" + rest.UrlPathSeparator + applicableActionsPath
	getContext := TestContext{
		Name: http.MethodGet + " " + subcatchmentUrlUnderTest + " request returns 404 (not found) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:    http.MethodGet,
			TargetUrl: subcatchmentUrlUnderTest,
		},
		ExpectedResponseStatus: http.StatusNotFound,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestInvalidApplicableActionsGetRequest_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	buildValidScenario(t, muxUnderTest)

	// when
	subcatchmentUrlUnderTest := baseSubcatchmentUrl + "/nope" + rest.UrlPathSeparator + applicableActionsPath
	getContext := TestContext{
		Name: http.MethodGet + " " + subcatchmentUrlUnderTest + " request returns 404 (not found) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:    http.MethodGet,
			TargetUrl: subcatchmentUrlUnderTest,
		},
		ExpectedResponseStatus: http.StatusNotFound,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestGetApplicableActionsResource_OkResponse(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	muxUnderTest := buildMuxUnderTest()
	buildValidScenario(t, muxUnderTest)

	// when
	getContext := TestContext{
		Name: http.MethodGet + " " + validSubcatchmentUrl + " request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:    http.MethodGet,
			TargetUrl: validSubcatchmentUrl + rest.UrlPathSeparator + applicableActionsPath,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	response := verifyResponseStatusCode(muxUnderTest, getContext)
	t.Log(response.RawResponse)

	jsonResponse := response.JsonMap
	g.Expect(len(jsonResponse)).To(BeNumerically("==", 1))
	if _, keyFound := jsonResponse["ApplicableActions"]; !keyFound {
		g.Expect(false).To(BeTrue(), "Missing expected [ApplicableActions] json name")
	}
	g.Expect(jsonResponse["ApplicableActions"]).To(ConsistOf("GullyRestoration", "RiverBankRestoration"))

	muxUnderTest.Shutdown()
}
