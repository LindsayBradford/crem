// Copyright (c) 2018 Australian Rivers Institute.

// See: https://astaxie.gitbooks.io/build-web-application-with-golang/en/03.2.html

package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/LindsayBradford/crm/config"
	"github.com/LindsayBradford/crm/logging/handlers"
)

type Status struct {
	Name    string
	Version string
	Message string
	Time    string
}

type Response struct {
	ResponseCode int
	Message      string
	Time         string
}

type AddressDelegateMapping struct {
	AddressPattern string
	Delegate       func(w http.ResponseWriter, r *http.Request)
}

type Delegates []AddressDelegateMapping

type RestServer struct {
	adminHandler *HttpHandler
	apiHandler   *HttpHandler

	configuration *config.HttpServerConfig
	Logger        handlers.LogHandler

	Status Status
}

func (s *RestServer) Initialise() *RestServer {
	s.adminHandler = new(HttpHandler)
	s.apiHandler = new(HttpHandler)

	s.AddAdminMapping(AddressDelegateMapping{AddressPattern: "/status", Delegate: s.statusHandler})
	s.AddAdminMapping(AddressDelegateMapping{AddressPattern: "/shutdown", Delegate: s.shutdownHandler})
	return s
}

func (s *RestServer) statusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.ServeMethodNotAllowedError(w, r)
		return
	}

	s.Logger.Debug("Responding with status [" + s.Status.Message + "]")
	s.UpdateStatusTime()

	statusJson, encodeError := json.MarshalIndent(s.Status, "", "  ")
	if encodeError != nil {
		s.Logger.Error(encodeError)
	}

	fmt.Fprintf(w, string(statusJson))
}

func (s *RestServer) shutdownHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.ServeMethodNotAllowedError(w, r)
		return
	}

	s.Status.Message = "SHUTTING_DOWN"
	s.Logger.Debug("Responding with status [" + s.Status.Message + "]")
	s.UpdateStatusTime()

	statusJson, encodeError := json.MarshalIndent(s.Status, "", "  ")
	if encodeError != nil {
		s.Logger.Error(encodeError)
	}

	bufferedWriter := bufio.NewWriter(w)

	fmt.Fprintf(bufferedWriter, string(statusJson))
	bufferedWriter.Flush()

	s.Logger.Warn("Shutting down")
	os.Exit(0)
}

func (s *RestServer) WithLogger(logger handlers.LogHandler) *RestServer {
	s.Logger = logger
	s.adminHandler.WithLogger(logger)
	s.apiHandler.WithLogger(logger)
	return s
}

func (s *RestServer) WithConfig(configuration *config.HttpServerConfig) *RestServer {
	s.configuration = configuration
	return s
}

func (s *RestServer) WithStatus(status Status) *RestServer {
	s.Status = status
	return s
}

func (s *RestServer) Start() {
	s.setStatus("CONFIGURING")

	serverAdminAddress := fmt.Sprintf(":%d", s.configuration.AdminPort)
	s.Logger.Debug("Starting Admin server listening on address [" + serverAdminAddress + "]")

	go func() {
		if err := http.ListenAndServe(serverAdminAddress, s.adminHandler); err != nil {
			s.setStatus("DEAD")
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	serverApiAddress := fmt.Sprintf(":%d", s.configuration.ApiPort)

	s.Logger.Debug("Starting API server listening on address [" + serverApiAddress + "]")

	s.setStatus("RUNNING")

	if err := http.ListenAndServe(serverApiAddress, s.apiHandler); err != nil {
		s.setStatus("DEAD")
		log.Fatal("ListenAndServe: ", err)
	}
}

func (s *RestServer) setStatus(statusMessage string) {
	s.Logger.Debug("Changed server Status to [" + statusMessage + "]")
	s.Status.Message = statusMessage
	s.UpdateStatusTime()
}

func (s *RestServer) UpdateStatusTime() {
	s.Status.Time = FormattedTimestamp()
}

func (s *RestServer) AddAdminMapping(mapping AddressDelegateMapping) {
	s.adminHandler.AddDelegateMapping(mapping)
}

func (s *RestServer) AddApiMapping(mapping AddressDelegateMapping) {
	s.apiHandler.AddDelegateMapping(mapping)
}

func (s *RestServer) ServeNotFoundError(w http.ResponseWriter, r *http.Request) {
	s.apiHandler.ServeError(http.StatusNotFound, "Resource not found", w, r)
}

func (s *RestServer) ServeMethodNotAllowedError(w http.ResponseWriter, r *http.Request) {
	s.apiHandler.ServeError(http.StatusMethodNotAllowed, "Method not allowed", w, r)
}

type HttpHandler struct {
	delegates Delegates
	logger    handlers.LogHandler
}

func (h *HttpHandler) WithLogger(logger handlers.LogHandler) *HttpHandler {
	h.logger = logger
	return h
}

func (h *HttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, delegate := range h.delegates {
		if r.URL.Path == delegate.AddressPattern {
			h.logger.Info("Handling request [" + r.URL.Path + "] from [" + r.RemoteAddr + "].")
			delegate.Delegate(w, r)
			return
		}
	}
	h.ServeNotFoundError(w, r)
}

func (h *HttpHandler) ServeNotFoundError(w http.ResponseWriter, r *http.Request) {
	h.ServeError(http.StatusNotFound, "Resource not found", w, r)
}

func (h *HttpHandler) ServeMethodNotAllowedError(w http.ResponseWriter, r *http.Request) {
	h.ServeError(http.StatusMethodNotAllowed, "Method not allowed", w, r)
}

func (h *HttpHandler) ServeError(responseCode int, responseMsg string, w http.ResponseWriter, r *http.Request) {
	response := Response{ResponseCode: responseCode, Message: responseMsg}
	response.Time = FormattedTimestamp()

	w.WriteHeader(responseCode)

	statusJson, encodeError := json.MarshalIndent(response, "", "  ")
	if encodeError != nil {
		h.logger.Error(encodeError)
	}

	h.logger.Warn("Request method [" + r.Method + "] for request [" + r.URL.Path + "] from [" + r.RemoteAddr + "] Responding with [" + responseMsg + "] error.")
	fmt.Fprintf(w, string(statusJson))
}

func (h *HttpHandler) AddDelegateMapping(mapping AddressDelegateMapping) {
	h.delegates = append(h.delegates, mapping)
}

func FormattedTimestamp() string {
	return fmt.Sprintf("%v", time.Now().Format(time.RFC3339Nano))
}
