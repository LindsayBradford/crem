package api

import (
	"encoding/json"
	"github.com/LindsayBradford/crem/internal/pkg/annealing/solution"
	"github.com/LindsayBradford/crem/internal/pkg/model"
	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	httptest "github.com/LindsayBradford/crem/internal/pkg/server/test"
	"github.com/LindsayBradford/crem/pkg/attributes"
	. "github.com/onsi/gomega"
	"net/http"
	"testing"
)

func TestFirstModelGetRequest_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	context := TestContext{
		Name: "GET /api/v1/model request returns 404 (not found) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "GET",
			TargetUrl:   baseUrl + "api/v1/model",
			RequestBody: "here is some text",
		},
		ExpectedResponseStatus: http.StatusNotFound,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)
	muxUnderTest.Shutdown()
}

func TestGetValidModelResource_OkResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	postContext := TestContext{
		Name: "POST /api/v1/scenario request returns 202 (accepted) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: validScenarioTomlConfig,
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)

	// when
	getContext := TestContext{
		Name: "GET /api/v1/model request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "GET",
			TargetUrl:   baseUrl + "api/v1/model",
			RequestBody: "here is some text",
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, getContext)
	muxUnderTest.Shutdown()
}

func TestModelPostRequest_NotAllowedResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	context := TestContext{
		Name: "PUT /model request returns 405 (method not allowed) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/model",
			RequestBody: "here is some text",
		},
		ExpectedResponseStatus: http.StatusMethodNotAllowed,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)
	muxUnderTest.Shutdown()
}

func TestModelPutRequest_NotAllowedResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	context := TestContext{
		Name: "PUT /model request returns 405 (not allowed) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "PUT",
			TargetUrl:   baseUrl + "api/v1/model",
			RequestBody: "here is some text",
		},
		ExpectedResponseStatus: http.StatusMethodNotAllowed,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)
	muxUnderTest.Shutdown()
}

func TestModePatchRequest_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	context := TestContext{
		Name: "PUT /model request returns 405 (method not allowed) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "PATCH",
			TargetUrl:   baseUrl + "api/v1/model",
			RequestBody: "here is some text",
		},
		ExpectedResponseStatus: http.StatusNotFound,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)
	muxUnderTest.Shutdown()
}

func TestModePatchRequest_NotJsonContentType_UnsupportedMediaTypeResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	postContext := TestContext{
		Name: "POST /api/v1/scenario request returns 202 (accepted) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: validScenarioTomlConfig,
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)

	// when
	context := TestContext{
		Name: "PUT /model request returns 415 (unsupported media type) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "PATCH",
			TargetUrl:   baseUrl + "api/v1/model",
			ContentType: rest.TextMimeType,
			RequestBody: "here is some text that isn't JSON",
		},
		ExpectedResponseStatus: http.StatusUnsupportedMediaType,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)
	muxUnderTest.Shutdown()
}

func TestModePatchRequest_NotJsonBody_BadRequestResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	postContext := TestContext{
		Name: "POST /api/v1/scenario request returns 202 (accepted) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: validScenarioTomlConfig,
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)

	// when
	context := TestContext{
		Name: "PUT /model request returns 400 (bad request) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "PATCH",
			TargetUrl:   baseUrl + "api/v1/model",
			ContentType: rest.JsonMimeType,
			RequestBody: "here is some text that isn't JSON",
		},
		ExpectedResponseStatus: http.StatusBadRequest,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)
	muxUnderTest.Shutdown()
}

func TestModePatchRequest_ValidJsonBody_OkResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()

	// when
	postContext := TestContext{
		Name: "POST /api/v1/scenario request returns 202 (accepted) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "POST",
			TargetUrl:   baseUrl + "api/v1/scenario",
			RequestBody: validScenarioTomlConfig,
			ContentType: rest.TomlMimeType,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, postContext)

	// when

	// when
	requestAttributes := attributes.Attributes{
		attributes.NameValuePair{
			Name:  "Summary",
			Value: "solution summary",
		},
		attributes.NameValuePair{
			Name:  "Encoding",
			Value: "A1",
		},
	}

	attributesAsJson, _ := json.Marshal(requestAttributes)
	attributesAsJsonString := string(attributesAsJson)

	context := TestContext{
		Name: "PUT /model request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "PATCH",
			TargetUrl:   baseUrl + "api/v1/model",
			ContentType: rest.JsonMimeType,
			RequestBody: attributesAsJsonString,
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)

	// when
	getContext := TestContext{
		Name: "GET /api/v1/model request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:      "GET",
			TargetUrl:   baseUrl + "api/v1/model",
			RequestBody: "here is some text",
		},
		ExpectedResponseStatus: http.StatusOK,
	}

	// then
	response := verifyResponseStatusCode(muxUnderTest, getContext)

	modelCopy := muxUnderTest.model.DeepClone()
	referenceModel := toCatchmentModel(modelCopy)
	referenceModel.Initialise(model.AsIs)
	referenceModel.SetManagementAction(0, true)
	referenceModel.SetManagementAction(5, true)
	referenceModel.SetManagementAction(7, true)
	referenceModel.JoiningAttributes(requestAttributes)

	expectedSolution := new(solution.SolutionBuilder).
		WithId(referenceModel.Id()).
		ForModel(referenceModel).
		Build()

	contentAsJsonBytes, _ := json.MarshalIndent(expectedSolution, "", "  ")
	expectedSolutionAsString := string(contentAsJsonBytes)

	g := NewGomegaWithT(t)
	g.Expect(response.RawResponse).To(Equal(expectedSolutionAsString))

	muxUnderTest.Shutdown()
}
