package http

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	api "gitlab.ubrato.ru/ubrato/cdn/api/gen"
	"gitlab.ubrato.ru/ubrato/cdn/internal/config"
)

type Server struct {
	srv    *http.Server
	logger *slog.Logger
}

func NewServer(logger *slog.Logger, cfg config.HTTP, router *Router) (*Server, error) {
	apiServer, err := api.NewServer(
		router,
		router.Auth,
		api.WithErrorHandler(router.HandleError),
	)
	if err != nil {
		return nil, fmt.Errorf("init http server: %w", err)
	}

	mux := chi.NewMux()
	registerMiddleware(mux)

	mux.Mount("/", apiServer)

	registerSwaggerUIHandlers(logger, mux, cfg.SwaggerUIPath)

	server := &Server{
		srv: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Port),
			Handler: mux,
		},
		logger: logger,
	}

	return server, nil
}

func newCORShandler() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://ubrato.ru",
			"https://ubrato.ru",
			"http://dev.ubrato.ru",
			"https://dev.ubrato.ru",
			"http://admin.ubrato.ru",
			"https://admin.ubrato.ru",
			"http://localhost",
			"http://localhost:5174",
			"http://localhost:5173",
			"http://localhost:3000",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"X-File-Type"},
		AllowCredentials: true,
	})
}

func (s *Server) Start() error {
	s.logger.Info("Starting http server", "addr", s.srv.Addr)
	err := s.srv.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	s.logger.Info("Shutting down http server")
	return s.srv.Shutdown(context.Background())
}
