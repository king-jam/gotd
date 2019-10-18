package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/king-jam/gotd/app"
	log "github.com/sirupsen/logrus"
)

func main() {
	// cause the new instance to be created
	gotd := app.New()

	// Catch signal so we can shutdown gracefully
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		// service connections
		log.Infof("GOTD Starting")
		if err := gotd.Start(); err != nil {
			log.Fatalf("GOTD Run Error: %s\n", err)
		}
	}()
	// defer will handle all the cleanup
	defer func() {
		err := gotd.Shutdown()
		if err != nil {
			log.Fatalf("GOTD Shutdown Error: %s\n", err)
		}
	}()

	// Wait for a signal before shutting down
	sig := <-sigCh
	log.Infof("%s Signal received. Shutting down GOTD\n", sig.String())
}
