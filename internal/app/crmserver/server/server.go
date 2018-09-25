// Copyright (c) 2018 Australian Rivers Institute.

// See: https://astaxie.gitbooks.io/build-web-application-with-golang/en/03.2.html

package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/LindsayBradford/crm/config"
	"github.com/LindsayBradford/crm/logging/handlers"
	"golang.org/x/net/context"
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

	var adminServer *http.Server
	var apiServer *http.Server

	go func() {

		serverAdminAddress := fmt.Sprintf(":%d", s.configuration.AdminPort)
		s.Logger.Debug("Starting Admin server listening on address [" + serverAdminAddress + "]")

		adminServer = &http.Server{Addr: serverAdminAddress, Handler: s.adminMux}

		if err := adminServer.ListenAndServe(); err != http.ErrServerClosed {
			s.setStatus("DEAD")
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	go func() {
		serverApiAddress := fmt.Sprintf(":%d", s.configuration.ApiPort)
		s.Logger.Debug("Starting API server listening on address [" + serverApiAddress + "]")
		apiServer = &http.Server{Addr: serverApiAddress, Handler: s.apiMux}

		if err := apiServer.ListenAndServe(); err != http.ErrServerClosed {
			s.setStatus("DEAD")
			log.Fatal("ListenAndServe: ", err)
		}

	}()

	s.setStatus("RUNNING")

	s.adminMux.WaitOnShutdown()

	apiServer.Shutdown(context.Background())
	adminServer.Shutdown(context.Background())

	s.Logger.Warn("Shutting down")
}

func (s *RestServer) setStatus(statusMessage string) {
	s.adminMux.setStatus(statusMessage)
}

func (s *RestServer) AddApiMapping(address string, handlerFunction http.HandlerFunc) {
	s.apiMux.handlerMap[address] = handlerFunction
}
