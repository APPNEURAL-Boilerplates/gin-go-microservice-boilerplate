package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/example/gin-microservice-boilerplate/internal/app"
	"github.com/example/gin-microservice-boilerplate/internal/config"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.SlogLevel(),
	}))

	router := app.NewRouter(cfg, logger)
	server := &http.Server{
		Addr:         cfg.Addr(),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeoutSeconds) * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("http server starting",
			"service", cfg.ServiceName,
			"environment", cfg.Environment,
			"address", cfg.Addr(),
		)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		logger.Error("http server failed", "error", err)
		os.Exit(1)
	case signal := <-shutdown:
		logger.Info("shutdown signal received", "signal", signal.String())
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.ShutdownTimeoutSeconds)*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
		if closeErr := server.Close(); closeErr != nil {
			logger.Error("forced shutdown failed", "error", closeErr)
		}
		os.Exit(1)
	}

	logger.Info("http server stopped")
}
