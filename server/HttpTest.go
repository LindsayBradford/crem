// Copyright (c) 2018 Australian Rivers Institute.

package server

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

type HttpTestRequestContext struct {
	method    string
	targetUrl string
	handler   http.HandlerFunc
}

type JsonResponseContainer struct {
	statusCode int
	jsonBody   map[string]interface{}
}

func (context *HttpTestRequestContext) BuildJsonResponse() JsonResponseContainer {
	response := context.buildResponse()

	return JsonResponseContainer{
		statusCode: response.StatusCode,
		jsonBody:   context.buildJsonMapFrom(response.Body),
	}
}

func (context *HttpTestRequestContext) buildResponse() *http.Response {
	w := httptest.NewRecorder()
	context.handler(w, context.newRequest())
	return w.Result()
}

func (context *HttpTestRequestContext) newRequest() *http.Request {
	return httptest.NewRequest(context.method, context.targetUrl, nil)
}

func (context *HttpTestRequestContext) buildJsonMapFrom(responseBody io.ReadCloser) map[string]interface{} {
	jsonMap := make(map[string]interface{})

	responseBodyBytes, _ := ioutil.ReadAll(responseBody)
	json.Unmarshal(responseBodyBytes, &jsonMap)

	return jsonMap
}
