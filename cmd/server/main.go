package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/juanpblasi/go-template/internal/config"
	"github.com/juanpblasi/go-template/internal/repository"
	"github.com/juanpblasi/go-template/internal/server"
	"github.com/juanpblasi/go-template/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	// 1. Configuración principal
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// 2. Logger setup
	if err := logger.InitLogger(cfg.Logger.Level); err != nil {
		panic("Failed to init logger: " + err.Error())
	}
	defer logger.Log.Sync()

	logger.Info("Starting service", zap.String("app", cfg.App.Name), zap.String("env", cfg.App.Env))

	// 3. Database setup
	db, err := repository.NewDB(cfg.DB)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// 4. Iniciar Servidor (Inyección de dependencias en server.go)
	srv := server.NewServer(cfg, db)

	// Error channel to capture server crashes
	errChan := make(chan error, 2)

	// Rutinas para arrancar servidores
	go func() {
		if err := srv.StartHTTP(context.Background()); err != nil {
			errChan <- err
		}
	}()

	go func() {
		if err := srv.StartGRPC(context.Background()); err != nil {
			errChan <- err
		}
	}()

	// 5. Soporte para Graceful Shutdown usando un canal y context
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		logger.Error("Server error occurred", zap.Error(err))
	case sig := <-quit:
		logger.Info("Shutting down cleanly", zap.String("signal", sig.String()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Graceful shutdown failed", zap.Error(err))
	}

	logger.Info("Server stopped cleanly")
}
