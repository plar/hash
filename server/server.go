package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/plar/hash/server/config"
	"github.com/plar/hash/service/hasher"
	"github.com/plar/hash/service/health"
	"github.com/plar/hash/service/stats"
)

// Server implements HTTP server
type Server struct {
	cfg  config.Config
	done chan error

	hasherSvc hasher.Service
	healthSvc health.Service

	router     *router
	httpServer *http.Server
}

// New creates a new HTTP server with config
func New(cfg config.Config, hasherSvc hasher.Service, statsSvc stats.Service, healthSvc health.Service) (*Server, chan error) {

	// initialize specific handlers for /hash and /stats endpoints
	hasherHandler := &hasherHandler{svc: hasherSvc}
	statsHandler := &statsHandler{svc: statsSvc}

	// initialize server
	s := &Server{
		cfg:       cfg,
		done:      make(chan error),
		hasherSvc: hasherSvc,
		healthSvc: healthSvc,
	}

	// initialize routes
	s.router = newRouter([]route{
		newRoute(http.MethodPost, "/hash", hasherHandler.createHash),
		newRoute(http.MethodGet, "/hash/([0-9]+)", hasherHandler.getHash),

		newRoute(http.MethodGet, "/stats", statsHandler.stats),

		newRoute(http.MethodGet, "/shutdown", s.shutdownHandler),
	})

	// and finally initialize HTTP server
	s.httpServer = &http.Server{
		Handler:      s,
		Addr:         s.cfg.ListenAddr(),
		ReadTimeout:  s.cfg.ReadTimeout(),
		WriteTimeout: s.cfg.WriteTimeout(),
		IdleTimeout:  s.cfg.IdleTimeout(),
	}

	return s, s.done
}

// Config returns the server configuration
func (s *Server) Config() config.Config {
	return s.cfg
}

// Run executes the server main loop
func (s *Server) Run() error {
	s.healthSvc.Healthy()
	log.Printf("The server is ready to handle requests at %v\n", s.cfg.ListenAddr())
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() {
	s.healthSvc.Unhealthy() // the server is sick, we won't get here again
	log.Println("the server is shutting down")

	s.hasherSvc.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ShutDownTimeout()*time.Second)
	defer cancel()

	s.httpServer.SetKeepAlivesEnabled(false) // release idle connections
	s.done <- s.httpServer.Shutdown(ctx)     // give 15 seconds for normal shuts down
}

// ServeHTTP handles all requests
// if the server is healthy then all requests will be sent to the router handler
// otherwise an error (503) will be returned
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.healthSvc.IsHealthy() {
		s.router.handler(w, r)
	} else {
		http.Error(w, "the server is shutting down", http.StatusServiceUnavailable)
	}
}

func (s *Server) shutdownHandler(w http.ResponseWriter, r *http.Request) {
	if s.healthSvc.IsHealthy() {
		// use goroutine to release http listeners in s.Shutdown()
		// and immediately quit from the shutdown handler
		go s.Shutdown()
	}
}
