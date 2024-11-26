package http

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gitlab.ubrato.ru/ubrato/cdn/api"
)

func serveBytes(b []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(b)
	}
}

func registerSwaggerUIHandlers(logger *slog.Logger, mux *chi.Mux, path string) {
	if path == "" {
		logger.Warn("mux: serve swagger ui: swagger_ui_path is empty; skip")

		return
	}

	if len(api.OpenapiSpec) == 0 {
		logger.Error("swagger spec is not set")

		return
	}

	logger.Info("mux: serve swagger ui files", slog.String("path", path))

	mux.Route("/swagger", func(r chi.Router) {
		r.Mount("/", http.StripPrefix("/swagger", http.FileServer(http.Dir(path))))
		r.Mount("/spec", serveBytes(api.OpenapiSpec))
	})
}
