package api

import (
	_ "embed"
	"net/http"
	"testing"

	"github.com/LindsayBradford/crem/internal/pkg/server/rest"
	httptest "github.com/LindsayBradford/crem/internal/pkg/server/test"
	. "github.com/onsi/gomega"
)

const (
	baseActionsUrl        = baseUrl + "api/v1/model/actions"
	applicableActionsPath = "applicable"
)

func TestFirstApplicableActionsGetRequest_NotFoundResponse(t *testing.T) {
	// given
	muxUnderTest := buildMuxUnderTest()
	actionsUrlUnderTest := baseActionsUrl + rest.UrlPathSeparator + applicableActionsPath

	// when
	context := TestContext{
		Name: http.MethodGet + " " + actionsUrlUnderTest + " request returns 404 (not found) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:    http.MethodGet,
			TargetUrl: actionsUrlUnderTest,
		},
		ExpectedResponseStatus: http.StatusNotFound,
	}

	// then
	verifyResponseStatusCode(muxUnderTest, context)
	muxUnderTest.Shutdown()
}

func TestGetApplicableActionsResource_OkResponse(t *testing.T) {
	g := NewGomegaWithT(t)

	// given
	muxUnderTest := buildMuxUnderTest()
	buildValidScenario(t, muxUnderTest)

	actionsUrlUnderTest := baseActionsUrl + rest.UrlPathSeparator + applicableActionsPath

	// when
	getContext := TestContext{
		Name: http.MethodGet + " " + actionsUrlUnderTest + " request returns 200 (ok) response",
		T:    t,
		Request: httptest.HttpTestRequestContext{
			Method:    http.MethodGet,
			TargetUrl: actionsUrlUnderTest,
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
	rawActions := jsonResponse["ApplicableActions"]
	if actionsMap, isMap := rawActions.(map[string]interface{}); isMap {
		g.Expect(len(actionsMap)).To(BeNumerically("==", 7))

		pu20actions := actionsMap["20"]
		if actionsArray, isArray := pu20actions.([]string); isArray {
			g.Expect(len(actionsArray)).To(BeNumerically("==", 1))
			g.Expect(pu20actions).To(ConsistOf("RiverBankRestoration"))
		}

		pu22actions := actionsMap["22"]
		if actionsArray, isArray := pu22actions.([]string); isArray {
			g.Expect(len(actionsArray)).To(BeNumerically("==", 2))
			g.Expect(pu20actions).To(ConsistOf("RiverBankRestoration", "WetlandsEstablishment"))
		}
	}

	muxUnderTest.Shutdown()
}
