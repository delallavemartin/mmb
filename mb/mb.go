package main

// This package provides the main entry point for the Message Broker (MB) application.
// It sets up the HTTP server, handles configuration, logging, and graceful shutdown.


import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mllave.com/mllave/mmb/mb/handlers"
	"mllave.com/mllave/mmb/mb/src/model/subscribers"
)

func main() {
	// Get configuration from environment variables
    port := os.Getenv("MB_PORT")
    if port == "" {
        port = "8080"
    }
    logPath := os.Getenv("MB_LOG_PATH")
    if logPath == "" {
        logPath = "log/imb.log"
    }

	// Set up zap logger
    logger, err := zap.NewProduction()
    if err != nil {
        logger.Fatal("can't initialize zap logger: %v", zap.Error(err))
    }
    if err := logger.Sync(); err != nil {
        logger.Error("failed to sync logger", zap.Error(err))
    }

    logger.Info("Server started", zap.String("port", port), zap.String("logPath", logPath))

	
	// SPEC: Only 10 subscribers are supported
	var aListOfSubscribers = subscribers.SubscribersList{Addresses: make([]string, 0, 10)}

    mux := http.NewServeMux()
    mux.HandleFunc("/notify", handlers.PublisherHandler(&aListOfSubscribers, logger))
    mux.HandleFunc("/subscribe", handlers.SubscriberHandler(&aListOfSubscribers, logger))

    server := &http.Server{
        Addr:    ":" + port,
        Handler: mux,
    }

	// Graceful shutdown
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		logger.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			logger.Error("HTTP server Shutdown", zap.Error(err))
		}
		close(idleConnsClosed)
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Error("ListenAndServe", zap.Error(err))
	}

	<-idleConnsClosed
	logger.Info("Server stopped.")
}
