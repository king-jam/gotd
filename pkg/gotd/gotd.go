package gotd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/king-jam/gotd/dashboard"
	"github.com/king-jam/gotd/gif"
	"github.com/king-jam/gotd/postgres"
	"github.com/king-jam/gotd/slack"
)

type App struct {
	database *gorm.DB
	server   *http.Server
}

// New return a default uninitialized App instance
func New() *App {
	return &App{}
}

// Start creates and starts all necessary services
func (a *App) Start() error {
	db, err := initializeDatabase()
	if err != nil {
		return err
	}

	a.database = db

	gifService, err := initializeGifService(db)
	if err != nil {
		return err
	}

	server, err := initializeHTTPServices(gifService)
	if err != nil {
		return err
	}

	a.server = server

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Shutdown tries to gracefully cleanup and shutdown all services
func (a *App) Shutdown() error {
	// handle the server shutdown process
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	if err := a.database.Close(); err != nil {
		return err
	}
	return nil
}

func initializeDatabase() (*gorm.DB, error) {
	dbString := os.Getenv("DATABASE_URL")
	if dbString == "" {
		return nil, errors.New("$DATABASE_URL must be set")
	}

	// NewClient a connection to the database and configure it
	// this must be
	db, err := postgres.NewClient(dbString)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the database: %s", err)
	}

	return db, nil
}

func initializeGifService(db *gorm.DB) (*gif.Service, error) {
	repo, err := gif.NewGIFRepo(db)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize the repository: %s", err)
	}

	if err = repo.InitDB(); err != nil {
		return nil, fmt.Errorf("unable to initialize the schemas: %s", err)
	}

	gifService := gif.NewGifService(repo)

	return gifService, nil
}

func initializeHTTPServices(gifService *gif.Service) (*http.Server, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return nil, errors.New("$PORT must be set")
	}

	siHandler := slack.New(gifService)
	dashboardHandler := dashboard.New(gifService)

	appMux := http.NewServeMux()
	appMux.Handle("/receive", siHandler)
	appMux.Handle("/gif", dashboardHandler)
	appMux.Handle("/", http.FileServer(http.Dir("./static/dashboard")))

	return &http.Server{
		Addr:    ":" + port,
		Handler: appMux,
	}, nil
}
