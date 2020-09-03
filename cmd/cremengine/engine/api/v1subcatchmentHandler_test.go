package api

import (
	"encoding/json"
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	httptest "github.com/LindsayBradford/crem/internal/pkg/server/test"
	"github.com/LindsayBradford/crem/pkg/attributes"
	"net/http"
	"testing"
)

const (
	validScenarioFile = "testdata/ValidTestScenario.toml"

	scenarioUrl = baseUrl + "api/v1/scenario"

	baseSubcatchmentUrl  = baseUrl + "api/v1/model/subcatchment"
	validSubcatchment    = "18"
	validSubcatchmentUrl = baseSubcatchmentUrl + rest.UrlPathSeparator + validSubcatchment
)

func TestFirstSubcatchmentGetRequest_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	subcatchmentUrlUnderTest := baseSubcatchmentUrl + rest.UrlPathSeparator + validSubcatchment

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

func TestMissingSubcatchmentGetRequest_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	buildValidScenario(t, muxUnderTest)

	// when
	subcatchmentUrlUnderTest := baseSubcatchmentUrl + "/1"
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

func TestInvalidSubcatchmentGetRequest_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	buildValidScenario(t, muxUnderTest)

	// when
	subcatchmentUrlUnderTest := baseSubcatchmentUrl + "/nope"
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

func TestDeleteValidSubcathmentResource_BadMethodResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	buildValidScenario(t, muxUnderTest)

	// when
	getContext := TestContext{
		Name: http.MethodDelete + " " + validSubcatchmentUrl + " request returns 405 (method not allowed) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:    http.MethodDelete,
			TargetUrl: validSubcatchmentUrl,
		},
		ExpectedResponseStatus: http.StatusMethodNotAllowed,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestGetValidSubcathmentResource_OkResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	buildValidScenario(t, muxUnderTest)

	// when
	getContext := TestContext{
		Name: http.MethodGet + " " + validSubcatchmentUrl + " request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:    http.MethodGet,
			TargetUrl: validSubcatchmentUrl,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()

	// TODO: Limitation of current test framework is that I can't easily get to response content.
	// TODO: Consider retrieval of response body and interrogating its Json payload for expected action state.
}

func TestFirstSubcatchmentPostRequest_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	subcatchmentUrlUnderTest := baseSubcatchmentUrl + rest.UrlPathSeparator + validSubcatchment

	// when
	context := TestContext{
		Name: http.MethodPost + " " + subcatchmentUrlUnderTest + " request returns 404 (not found) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:    http.MethodPost,
			TargetUrl: subcatchmentUrlUnderTest,
		},
		ExpectedResponseStatus: http.StatusNotFound,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)
	muxUnderTest.Shutdown()
}

func TestMissingSubcatchmentPostRequest_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	buildValidScenario(t, muxUnderTest)

	// when
	subcatchmentUrlUnderTest := baseSubcatchmentUrl + "/1"
	getContext := TestContext{
		Name: http.MethodPost + " " + subcatchmentUrlUnderTest + " request returns 404 (not found) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:    http.MethodPost,
			TargetUrl: subcatchmentUrlUnderTest,
		},
		ExpectedResponseStatus: http.StatusNotFound,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestInvalidSubcatchmentPostRequest_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	buildValidScenario(t, muxUnderTest)

	// when
	subcatchmentUrlUnderTest := baseSubcatchmentUrl + "/nope"
	getContext := TestContext{
		Name: http.MethodPost + " " + subcatchmentUrlUnderTest + " request returns 404 (not found) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:    http.MethodPost,
			TargetUrl: subcatchmentUrlUnderTest,
		},
		ExpectedResponseStatus: http.StatusNotFound,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestPostValidSubcathmentResource_OkResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	buildValidScenario(t, muxUnderTest)

	actionAttributes := attributes.Attributes{}.
		Add("GullyRestoration", ActiveAction).
		Add("RiverBankRestoration", ActiveAction).
		Add("HillSlopeRestoration", InactiveAction)

	actionStatusBytes, _ := json.Marshal(actionAttributes)

	// when
	getContext := TestContext{
		Name: http.MethodPost + " " + validSubcatchmentUrl + " request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      http.MethodPost,
			TargetUrl:   validSubcatchmentUrl,
			RequestBody: string(actionStatusBytes),
			ContentType: rest.JsonMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestPostInvalidSubcatchmentJson_OkResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	buildValidScenario(t, muxUnderTest)

	invalidJson := "This is not valid\" Json}"

	// when
	getContext := TestContext{
		Name: http.MethodPost + " " + validSubcatchmentUrl + " request returns 400 (bad request) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      http.MethodPost,
			TargetUrl:   validSubcatchmentUrl,
			RequestBody: string(invalidJson),
			ContentType: rest.JsonMimeType,
		},
		ExpectedResponseStatus: http.StatusBadRequest,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestPostInvalidSubcatchmentActionResource_OkResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	buildValidScenario(t, muxUnderTest)

	actionAttributes := attributes.Attributes{}.
		Add("NonExistentAction", ActiveAction)

	actionStatusBytes, _ := json.Marshal(actionAttributes)

	// when
	getContext := TestContext{
		Name: http.MethodPost + " " + validSubcatchmentUrl + " request returns 400 (bad request) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      http.MethodPost,
			TargetUrl:   validSubcatchmentUrl,
			RequestBody: string(actionStatusBytes),
			ContentType: rest.JsonMimeType,
		},
		ExpectedResponseStatus: http.StatusBadRequest,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestPostInvalidSubcatchmentActionStateResource_OkResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	buildValidScenario(t, muxUnderTest)

	actionAttributes := attributes.Attributes{}.
		Add("GullyRestoration", "ThisIsNotAValidState")

	actionStatusBytes, _ := json.Marshal(actionAttributes)

	// when
	getContext := TestContext{
		Name: http.MethodPost + " " + validSubcatchmentUrl + " request returns 400 (bad request) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      http.MethodPost,
			TargetUrl:   validSubcatchmentUrl,
			RequestBody: string(actionStatusBytes),
			ContentType: rest.JsonMimeType,
		},
		ExpectedResponseStatus: http.StatusBadRequest,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func buildValidScenario(t *testing.T, muxUnderTest *Mux) {
	scenarioTomlText := readTestFileAsText(validScenarioFile)

	// when
	postContext := TestContext{
		Name: http.MethodPost + scenarioUrl + " request returns 200 (accepted) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      http.MethodPost,
			TargetUrl:   scenarioUrl,
			RequestBody: scenarioTomlText,
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)
}