package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/plar/hash/domain/repository"
	"github.com/plar/hash/infra/persistence/memory"
	"github.com/plar/hash/server"
	"github.com/plar/hash/server/config"
	"github.com/plar/hash/service/hasher"
	"github.com/plar/hash/service/health"
	"github.com/plar/hash/service/stats"
)

var (
	// Better to use Google Wire DI framework
	healthSvc health.Service
	statsSvc  stats.Service
	hashRepo  repository.HashRepository
	hasherSvc hasher.Service
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Configuration error: %v\n", err)
	}

	// set the number of spare OS threads
	runtime.GOMAXPROCS(int(cfg.TotalWorkers()))

	// create services, repos
	statsSvc = stats.New()
	statsSvc = stats.NewLoggingService(statsSvc)

	hashRepo = memory.NewHashRepository()

	hasherSvc = hasher.New(hashRepo, cfg)
	hasherSvc = hasher.NewInstrumentingService(hasherSvc, statsSvc)
	hasherSvc = hasher.NewLoggingService(hasherSvc)

	healthSvc = health.NewService()

	// create server
	server, done := server.New(cfg, hasherSvc, statsSvc, healthSvc)

	// handle OS signals...
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	go func() {
		log.Printf("Signal received: %v\n", <-quit)
		server.Shutdown()
	}()

	// start HTTP server
	err = server.Run()
	if err != nil && err != http.ErrServerClosed /* happens after Shutdown, ignore it */ {
		log.Fatalf("Could not listen on %s: %v\n", server.Config().ListenAddr(), err)
	}

	// quit, maybe with error
	err = <-done
	if err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v", err)
	}

	log.Println("The server has been shutdown")
}
