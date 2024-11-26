package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/snowflake"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gitlab.ubrato.ru/ubrato/cdn/internal/config"
	"gitlab.ubrato.ru/ubrato/cdn/internal/lib/token"
	authService "gitlab.ubrato.ru/ubrato/cdn/internal/service/auth"
	s3Service "gitlab.ubrato.ru/ubrato/cdn/internal/service/s3"
	"gitlab.ubrato.ru/ubrato/cdn/internal/transport/http"
	authHandler "gitlab.ubrato.ru/ubrato/cdn/internal/transport/http/handlers/auth"
	errorHandler "gitlab.ubrato.ru/ubrato/cdn/internal/transport/http/handlers/error"
	fileHandler "gitlab.ubrato.ru/ubrato/cdn/internal/transport/http/handlers/file"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	cfg, err := config.ReadConfig()
	if err != nil {
		logger.Error("Error parsing default config from env", "error", err)
		os.Exit(1)
	}

	if cfg.Debug {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))

		logger.Warn("Debug enabled")
	}

	if err := run(cfg, logger); err != nil {
		logger.Error("Error initializing service", "error", err)
		os.Exit(1)
	}
}

func run(cfg config.Default, logger *slog.Logger) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	tokenAuthorizer, err := token.NewTokenAuthorizer(cfg.Auth.JWT)
	if err != nil {
		return fmt.Errorf("init token authorizer: %w", err)
	}

	node, err := snowflake.NewNode(1)
	if err != nil {
		return fmt.Errorf("init snowflake: %w", err)
	}

	minioClient, err := minio.New(cfg.Gateway.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Gateway.Minio.AccessKeyID, cfg.Gateway.Minio.SecretAccessKey, ""),
		Secure: false,
	})

	authService := authService.New(
		tokenAuthorizer,
	)

	s3Service := s3Service.New(
		minioClient,
		node,
	)

	router := http.NewRouter(http.RouterParams{
		Error: errorHandler.New(logger),
		Auth:  authHandler.New(logger, authService),
		File:  fileHandler.New(logger, s3Service),
	})

	server, err := http.NewServer(logger, cfg.Transport.HTTP, router)
	if err != nil {
		return fmt.Errorf("create http server: %w", err)
	}

	go func() {
		<-sig
		logger.Info("Received termination signal, cleaning up")
		err := server.Stop()
		if err != nil {
			logger.Error("Stop http server", "error", err)
		}
	}()

	err = server.Start()
	if err != nil {
		return fmt.Errorf("serve http: %w", err)
	}

	return nil
}
