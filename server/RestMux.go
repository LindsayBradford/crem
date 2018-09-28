// Copyright (c) 2018 Australian Rivers Institute.

package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/LindsayBradford/crm/logging/handlers"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

type RestMux interface {
	Start(portAddress string)
	Shutdown()

	SetLogger(handler handlers.LogHandler)
	AddHandler(address string, handler http.HandlerFunc)
}

type BaseMux struct {
	http.ServeMux
	muxType string
	server  http.Server

	handlerMap HandlerFunctionMap
	logger     handlers.LogHandler
}

type HandlerFunctionMap map[string]http.HandlerFunc

func (rm *BaseMux) Initialise() *BaseMux {
	rm.handlerMap = make(HandlerFunctionMap)
	return rm
}

func (rm *BaseMux) WithType(muxType string) *BaseMux {
	rm.muxType = muxType
	return rm
}

func (rm *BaseMux) SetLogger(logger handlers.LogHandler) {
	rm.logger = logger
}

func (rm *BaseMux) Logger() handlers.LogHandler {
	return rm.logger
}

func (rm *BaseMux) AddHandler(address string, handler http.HandlerFunc) {
	rm.handlerMap[address] = handler
}

func (rm *BaseMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rm.logRequestReceipt(r)
	if handlerFunction, handlerFound := rm.handlerFor(r); handlerFound {
		handlerFunction(w, r)
	} else {
		rm.ServeNotFoundError(w, r)
	}
}

func (rm *BaseMux) logRequestReceipt(r *http.Request) {
	rm.logger.Info(
		"[" + rm.muxType + "] Received request method [" + r.Method +
			"] for request [" + r.URL.Path + "] from [" + r.RemoteAddr + "].")
}

func (rm *BaseMux) handlerFor(r *http.Request) (handlerFunction http.HandlerFunc, found bool) {
	handlerFunction, found = rm.handlerMap[r.URL.String()]
	return
}

func (rm *BaseMux) ServeNotFoundError(w http.ResponseWriter, r *http.Request) {
	rm.ServeError(http.StatusNotFound, "Resource not found", w, r)
}

func (rm *BaseMux) ServeMethodNotAllowedError(w http.ResponseWriter, r *http.Request) {
	rm.ServeError(http.StatusMethodNotAllowed, "Method not allowed", w, r)
}

func (rm *BaseMux) ServeError(responseCode int, responseMsg string, w http.ResponseWriter, r *http.Request) {
	response := Response{ResponseCode: responseCode, Message: responseMsg}
	response.Time = FormattedTimestamp()

	w.WriteHeader(responseCode)

	statusJson, encodeError := json.MarshalIndent(response, "", "  ")
	if encodeError != nil {
		rm.logger.Error(encodeError)
	}

	rm.logResponseError(r, responseMsg)

	fmt.Fprintf(w, string(statusJson))
}

func (rm *BaseMux) logResponseError(r *http.Request, responseMsg string) {
	rm.logger.Warn(
		"Request method [" + r.Method + "] for request [" + r.URL.Path + "] from [" + r.RemoteAddr +
			"] Responding with [" + responseMsg + "] error.")
}

func (rm *BaseMux) Start(address string) {
	rm.logger.Debug("Starting [" + rm.muxType + "] server on address [" + address + "]")

	rm.server = http.Server{Addr: address, Handler: rm}

	if err := rm.server.ListenAndServe(); err != http.ErrServerClosed {
		wrappedErr := errors.Wrap(err, "ListenAndServe")
		rm.logger.Error(wrappedErr)
	}
}

func (rm *BaseMux) Shutdown() {
	rm.server.Shutdown(context.Background())
}
