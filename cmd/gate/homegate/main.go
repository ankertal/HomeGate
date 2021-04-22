package main

import (
	"os"

	"github.com/gate/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	config, err := server.LoadConfig()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
		os.Exit(1)
	}

	// default nil channel for HTTP
	doneHTTP := make(chan error)
	go func() {
		log.Infof("HTTP service is at: http://0.0.0.0:%s/", config.Port)
		srv := server.NewServer(config)
		err := srv.Server.ListenAndServe()
		doneHTTP <- err
	}()

	err = <-doneHTTP
	if err != nil {
		log.Fatal(err)
	}
}
