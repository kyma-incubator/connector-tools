package main

import (
	"context"
	"github.com/kyma-incubator/connector-tools/litmos-event-gw/internal/config"
	"github.com/kyma-incubator/connector-tools/litmos-event-gw/internal/logger"
	"github.com/kyma-incubator/connector-tools/litmos-event-gw/internal/router"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger.Initialize()
	defer logger.Logger.Sync()

	logger.Logger.Info("starting service...")

	config.ParseFlags()

	rtr := router.New()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	srv := http.Server{
		Addr:    ":" + "8080",
		Handler: rtr,
	}

	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	killSignal := <-interrupt

	switch killSignal {
	case os.Interrupt:
		log.Println("got os interrupt")
	case syscall.SIGTERM:
		log.Println("got SIGTERM")
	}
	log.Println("system is shutting down")

	srv.Shutdown(context.Background())

	logger.Logger.Info("done...")
}
