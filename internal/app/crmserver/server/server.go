// Copyright (c) 2018 Australian Rivers Institute.

// See: https://astaxie.gitbooks.io/build-web-application-with-golang/en/03.2.html

package server

import (
	"fmt"
	"log"
	"net/http"

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

type RestServer struct {
	adminMux *AdminMux
	apiMux   *ApiMux

	configuration *config.HttpServerConfig
	Logger        handlers.LogHandler
}

func (s *RestServer) Initialise() *RestServer {
	s.adminMux = new(AdminMux).Initialise().WithType("Admin")
	s.apiMux = new(ApiMux).Initialise().WithType("API")
	return s
}

func (s *RestServer) WithLogger(logger handlers.LogHandler) *RestServer {
	s.Logger = logger
	s.adminMux.Logger = logger
	s.apiMux.Logger = logger
	return s
}

func (s *RestServer) WithConfig(configuration *config.HttpServerConfig) *RestServer {
	s.configuration = configuration
	return s
}

func (s *RestServer) WithStatus(status Status) *RestServer {
	s.adminMux.Status = status
	return s
}

func (s *RestServer) Start() {
	s.setStatus("CONFIGURING")

	serverAdminAddress := fmt.Sprintf(":%d", s.configuration.AdminPort)
	s.Logger.Debug("Starting Admin server listening on address [" + serverAdminAddress + "]")

	go func() {
		if err := http.ListenAndServe(serverAdminAddress, s.adminMux); err != nil {
			s.setStatus("DEAD")
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	serverApiAddress := fmt.Sprintf(":%d", s.configuration.ApiPort)

	s.Logger.Debug("Starting API server listening on address [" + serverApiAddress + "]")

	s.setStatus("RUNNING")

	if err := http.ListenAndServe(serverApiAddress, s.apiMux); err != nil {
		s.setStatus("DEAD")
		log.Fatal("ListenAndServe: ", err)
	}
}

func (s *RestServer) setStatus(statusMessage string) {
	s.adminMux.setStatus(statusMessage)
}

func (s *RestServer) AddApiMapping(address string, handlerFunction http.HandlerFunc) {
	s.apiMux.handlerMap[address] = handlerFunction
}
