package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"todolist/backend/internal/app/routes"
	"todolist/backend/internal/infra/config"
	"todolist/backend/internal/infra/db"
	"todolist/backend/internal/infra/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	logg := logger.New(cfg.App.Env)
	defer func() { _ = logg.Sync() }()

	dbConn, err := db.Connect(cfg, logg)
	if err != nil {
		logg.Fatal("failed to connect database", zapError(err))
	}

	if err := db.AutoMigrate(dbConn, logg); err != nil {
		logg.Fatal("failed to run migrations", zapError(err))
	}

	engine := routes.SetupRouter(cfg, logg, dbConn)

	srv := serverConfig(cfg, engine)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logg.Fatal("server error", zapError(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logg.Error("server shutdown failed", zapError(err))
	}
}

const (
	ReadHeaderTimeout = 5 * time.Second
	ReadTimeout       = 10 * time.Second
	WriteTimeout      = 15 * time.Second
	IdleTimeout       = 60 * time.Second
)

func serverConfig(cfg *config.Config, engine http.Handler) *http.Server {
	return &http.Server{
		Addr:              cfg.App.Addr(),
		Handler:           engine,
		ReadHeaderTimeout: ReadHeaderTimeout,
		ReadTimeout:       ReadTimeout,
		WriteTimeout:      WriteTimeout,
		IdleTimeout:       IdleTimeout,
	}
}

func zapError(err error) zap.Field {
	return zap.Error(err)
}
