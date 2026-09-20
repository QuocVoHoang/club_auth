package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/your-org/go-base/internal/config"
	"github.com/your-org/go-base/internal/domain/usecase"
	"github.com/your-org/go-base/internal/framework/route"
	"github.com/your-org/go-base/internal/infrastructure/db/postgres"
	grpcserver "github.com/your-org/go-base/internal/infrastructure/grpc"
	gprcUser "github.com/your-org/go-base/internal/infrastructure/grpc/generated/user"
	"github.com/your-org/go-base/internal/infrastructure/handler"
	"github.com/your-org/go-base/internal/platform/database"
	"github.com/your-org/go-base/pkg/log"
	middlewarepkg "github.com/your-org/go-base/pkg/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 30 * time.Second
	idleTimeout       = 60 * time.Second
	shutdownTimeout   = 30 * time.Second
)

func Run(cfg *config.Config) error {
	if err := log.SetLevel(cfg.LogLevel); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}

	db, err := database.InitDatabase(database.Config{
		User:         cfg.Database.DBUser,
		Password:     cfg.Database.DBPass,
		Host:         cfg.Database.DBHost,
		Port:         cfg.Database.DBPort,
		DatabaseName: cfg.Database.DBName,
		SSLMode:      cfg.Database.SSLMode,
		ENV:          cfg.ENV,
	})
	if err != nil {
		return fmt.Errorf("initialize database: %w", err)
	}
	defer func() {
		sqlDB, err := db.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
	}()

	userRepo := postgres.NewUserRepository(db)
	jwtService := middlewarepkg.NewJWT(cfg.JWTSecret)
	registerUsecase := usecase.NewRegisterUsecase(userRepo, jwtService)
	loginUsecase := usecase.NewLoginUsecase(userRepo, jwtService)
	getCurrentUserUsecase := usecase.NewGetCurrentUserUsecase(userRepo)
	updateCurrentUserUsecase := usecase.NewUpdateCurrentUserUsecase(userRepo)
	authHandler := handler.NewAuthHandler(registerUsecase, loginUsecase)
	userHandler := handler.NewUserHandler(getCurrentUserUsecase, updateCurrentUserUsecase)

	grpcSrv := grpc.NewServer()
	gprcUser.RegisterUserServiceServer(grpcSrv, grpcserver.NewUserServer(userRepo))
	reflection.Register(grpcSrv)
	grpcLis, err := net.Listen("tcp", ":"+cfg.GRPCPort)
	if err != nil {
		return fmt.Errorf("listen gRPC on :%s: %w", cfg.GRPCPort, err)
	}

	router := route.NewRouter(cfg, db, jwtService, authHandler, userHandler)
	server := &http.Server{
		Addr:              ":" + cfg.PORT,
		Handler:           router,
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	serverErr := make(chan error, 1)
	go func() {
		log.Infof("Starting %s HTTP in %s environment on port %s", cfg.AppName, cfg.ENV, cfg.PORT)
		serverErr <- server.ListenAndServe()
	}()

	go func() {
		log.Infof("Starting %s gRPC in %s environment on port %s", cfg.AppName, cfg.ENV, cfg.GRPCPort)
		serverErr <- grpcSrv.Serve(grpcLis)
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("run server: %w", err)
		}
	case <-ctx.Done():
		log.Info("Shutdown signal received")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	grpcSrv.GracefulStop()

	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown HTTP server: %w", err)
	}

	log.Info("Server stopped")
	return nil
}
