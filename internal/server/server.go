package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	pb "github.com/juanpblasi/go-template/api/proto/v1"
	"github.com/juanpblasi/go-template/internal/config"
	grpchandler "github.com/juanpblasi/go-template/internal/handler/grpc"
	httphandler "github.com/juanpblasi/go-template/internal/handler/http"
	"github.com/juanpblasi/go-template/internal/repository"
	"github.com/juanpblasi/go-template/internal/service"
	"github.com/juanpblasi/go-template/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type Server struct {
	httpServer *http.Server
	grpcServer *grpc.Server
	cfg        *config.Config
}

func NewServer(cfg *config.Config, db *gorm.DB) *Server {
	// Initialize Dependencies
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	// Setup HTTP Router
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Register HTTP Handlers
	httphandler.RegisterHealthRoutes(r)
	httphandler.RegisterUserRoutes(r, userService)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler: r,
	}

	// Setup gRPC Server
	grpcServer := grpc.NewServer()
	userGrpcHandler := grpchandler.NewUserGrpcHandler(userService)
	pb.RegisterUserServiceServer(grpcServer, userGrpcHandler)

	return &Server{
		httpServer: httpServer,
		grpcServer: grpcServer,
		cfg:        cfg,
	}
}

func (s *Server) StartHTTP(ctx context.Context) error {
	logger.Info("Starting HTTP server", zap.Int("port", s.cfg.HTTP.Port))
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Server) StartGRPC(ctx context.Context) error {
	logger.Info("Starting gRPC server", zap.Int("port", s.cfg.GRPC.Port))
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.GRPC.Port))
	if err != nil {
		return err
	}
	if err := s.grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down gRPC server")
	s.grpcServer.GracefulStop()

	logger.Info("Shutting down HTTP server")
	return s.httpServer.Shutdown(ctx)
}
