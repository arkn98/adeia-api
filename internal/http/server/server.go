/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/. */

package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"adeia/internal/config"
	"adeia/pkg/log"
	"adeia/pkg/util/constants"

	"github.com/go-chi/chi"
)

// Controller is the interface for all methods of the controllers.
type Controller interface {
	Pattern() string
	Handler() http.Handler
}

// Server represents the server.
type Server struct {
	config      *config.ServerConfig
	controllers []Controller
	log         log.Logger
	srv         chi.Router
}

// New creates a new *Server.
func New(conf *config.ServerConfig, log log.Logger, controllers ...Controller) *Server {
	log.Debug("initializing new API server...")
	return &Server{
		config:      conf,
		controllers: controllers,
		log:         log,
		srv:         chi.NewRouter(),
	}
}

// BindControllers binds all the controllers to the Server.
func (s *Server) BindControllers() {
	s.log.Debug("binding handles to router...")

	s.srv.Route("/"+constants.APIVersion, func(r chi.Router) {
		for _, controller := range s.controllers {
			r.Mount(controller.Pattern(), controller.Handler())
		}
	})
}

// Serve starts serving the API. All server-related errors are handled here.
func (s *Server) Serve() {
	addr := s.config.Host + ":" + strconv.Itoa(s.config.Port)
	// TODO: add timeouts
	// TODO: add TLS support
	// TODO: add rate-limiter
	srv := &http.Server{
		Addr:    addr,
		Handler: s.srv,
	}

	// catch server errors in a channel
	serverErrs := make(chan error, 1)

	go func() {
		s.log.Infof("starting server on %q", addr)
		serverErrs <- srv.ListenAndServe()
	}()

	// chan to listen for ctrl+c, SIGTERM
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	// gracefully shutdown server
	select {
	case err := <-serverErrs:
		if err != http.ErrServerClosed {
			s.log.Errorf("error while serving: %v", err)
		}

	case sig := <-interruptChan:
		s.log.Infof("received: %v; starting shutdown...", sig)

		// wait for 5 seconds for pending requests to complete
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// cancel requests that exceed the context deadline (5 seconds)
		if err := srv.Shutdown(ctx); err != nil {
			s.log.Errorf("failed to gracefully shutdown server: %v", err)
		} else {
			s.log.Info("server gracefully stopped")
		}
	}
}
