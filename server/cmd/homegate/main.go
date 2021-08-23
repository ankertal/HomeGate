package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gate/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	config, err := server.LoadConfig()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
		os.Exit(1)
	}

	srv := server.NewServer(config)

	errC := make(chan error, 1)

	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		log.Infof("HTTP service is at: http://0.0.0.0:%s/", config.Port)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errC <- err
		}
	}()

	// wait for termination
	<-ctx.Done()

	log.Infof("Shutdown signal received")

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer func() {
		stop()
		cancel()
		close(errC)
	}()

	srv.SetKeepAlivesEnabled(false)

	if err := srv.Shutdown(ctxTimeout); err != nil {
		errC <- err
	}

	log.Infof("Shutdown completed")
}
