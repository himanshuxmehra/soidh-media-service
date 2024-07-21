package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"soidh-media-service/internal/api"
	"soidh-media-service/internal/config"
	"soidh-media-service/internal/database"
	"soidh-media-service/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logger := logger.New(cfg.LogLevel)

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
	}
	defer db.Close()

	router := api.SetupRoutes(db, logger)

	server := &http.Server{
		Addr:    cfg.ServerAddr,
		Handler: router,
	}

	go func() {
		logger.Info("Starting server", "addr", cfg.ServerAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exiting")
}