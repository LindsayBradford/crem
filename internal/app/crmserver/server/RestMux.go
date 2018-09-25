// Copyright (c) 2018 Australian Rivers Institute.

package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/LindsayBradford/crm/logging/handlers"
)

type RestMux struct {
	http.ServeMux
	muxType string

	handlerMap HandlerFunctionMap
	Logger     handlers.LogHandler
}

type HandlerFunctionMap map[string]http.HandlerFunc

func (rm *RestMux) Initialise() *RestMux {
	rm.handlerMap = make(HandlerFunctionMap)
	return rm
}

func (rm *RestMux) WithType(muxType string) *RestMux {
	rm.muxType = muxType
	return rm
}

func (rm *RestMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rm.logRequestReceipt(r)
	if handlerFunction, handlerFound := rm.handlerFor(r); handlerFound {
		handlerFunction(w, r)
	} else {
		rm.ServeNotFoundError(w, r)
	}
}

func (rm *RestMux) logRequestReceipt(r *http.Request) {
	rm.Logger.Debug(
		"[" + rm.muxType + "] Received request method [" + r.Method +
			"] for request [" + r.URL.Path + "] from [" + r.RemoteAddr + "].")
}

func (rm *RestMux) handlerFor(r *http.Request) (handlerFunction http.HandlerFunc, found bool) {
	handlerFunction, found = rm.handlerMap[r.URL.String()]
	return
}

func (rm *RestMux) ServeNotFoundError(w http.ResponseWriter, r *http.Request) {
	rm.ServeError(http.StatusNotFound, "Resource not found", w, r)
}

func (rm *RestMux) ServeMethodNotAllowedError(w http.ResponseWriter, r *http.Request) {
	rm.ServeError(http.StatusMethodNotAllowed, "Method not allowed", w, r)
}

func (rm *RestMux) ServeError(responseCode int, responseMsg string, w http.ResponseWriter, r *http.Request) {
	response := Response{ResponseCode: responseCode, Message: responseMsg}
	response.Time = FormattedTimestamp()

	w.WriteHeader(responseCode)

	statusJson, encodeError := json.MarshalIndent(response, "", "  ")
	if encodeError != nil {
		rm.Logger.Error(encodeError)
	}

	rm.logResponseError(r, responseMsg)

	fmt.Fprintf(w, string(statusJson))
}

func (rm *RestMux) logResponseError(r *http.Request, responseMsg string) {
	rm.Logger.Warn(
		"Request method [" + r.Method + "] for request [" + r.URL.Path + "] from [" + r.RemoteAddr +
			"] Responding with [" + responseMsg + "] error.")
}
