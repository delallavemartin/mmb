package main

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"mllave.com/mllave/mmb/mb/src/messengerservice/mail"
	"mllave.com/mllave/mmb/mb/src/messengerservice/postoffice"
	"mllave.com/mllave/mmb/mb/src/messengerservice/delivery"
	"mllave.com/mllave/mmb/mb/src/reader"
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
        log.Fatalf("can't initialize zap logger: %v", err)
    }
    defer logger.Sync()

    logger.Info("Server started", zap.String("port", port), zap.String("logPath", logPath))

	
	// SPEC: Only 10 subscribers are supported
	var aListOfSubscribers = subscribers.SubscribersList{Addresses: make([]string, 0, 10)}

    mux := http.NewServeMux()
    mux.HandleFunc("/notify", PublisherHandler(&aListOfSubscribers, logger))
    mux.HandleFunc("/subscribe", SubscriberHandler(&aListOfSubscribers, logger))

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
