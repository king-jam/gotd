package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/king-jam/gotd/dashboard"
	"github.com/king-jam/gotd/postgres"
	"github.com/king-jam/gotd/slack_integration"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	dbString := os.Getenv("DATABASE_URL")
	if dbString == "" {
		log.Fatal("$DATABASE_URL must be set")
	}

	// Catch signal so we can shutdown gracefully
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	dbURL, err := url.Parse(dbString)
	if err != nil {
		log.Fatal("Invalid Database URL format")
	}

	db, err := postgres.InitDatabase(dbURL)
	if err != nil {
		log.Fatalf("Unable to initialize the Database: %s", err)
	}
	defer db.Close()

	siHandler := slack_integration.New(db)

	dashboardHandler := dashboard.New(db)

	appMux := http.NewServeMux()
	appMux.Handle("/receive", siHandler)
	appMux.Handle("/gif", dashboardHandler)
	appMux.Handle("/", http.FileServer(http.Dir("./static/dashboard")))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: appMux,
	}

	go func() {
		// service connections
		fmt.Println("[INFO] Server listening")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen Error: %s\n", err)
		}
	}()
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
	}()

	// Wait for a signal
	sig := <-sigCh
	log.Printf("%s Signal received. Shutting down Application.", sig.String())
}
