package http

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gitlab.ubrato.ru/ubrato/cdn/internal/transport/http/middlewares"
)

func registerMiddleware(mux *chi.Mux) {
	mux.Use(newCORShandler())
	mux.Use(middleware.RealIP)
	mux.Use(middleware.RequestID)
	mux.Use(middlewares.RequestIDToResponse)
	mux.Use(middleware.Recoverer)
	mux.Use(middlewares.Logger)
	mux.Use(middleware.Timeout(1 * time.Second))
}
