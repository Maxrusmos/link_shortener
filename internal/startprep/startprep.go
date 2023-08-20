package startprep

import (
	"link_shortener/internal/middleware"
	"link_shortener/internal/services"
	"link_shortener/internal/storage"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func RegisterRoutes(r chi.Router, storageURL storage.URLStorage, baseURL string, logger *zap.Logger) {
	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		services.HandleGetRequest(w, r, storageURL)
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		services.Ping(storageURL)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		services.HandlePostRequest(w, r, storageURL, baseURL)
	})
	r.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		services.ShortenHandler(w, r, storageURL, baseURL)
	})
	r.Post("/api/shorten/batch", func(w http.ResponseWriter, r *http.Request) {
		services.HandleBatchShorten(w, r, storageURL, baseURL)
	})
}

func StartServer(address string, r chi.Router, logger *zap.Logger) {
	server := &http.Server{
		Addr:         address,
		Handler:      middleware.CompressionMiddleware(middleware.LoggingMiddleware(logger, r)),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	logger.Info("server started")
	if err := server.ListenAndServe(); err != nil {
		logger.Error("server stopped", zap.Error(err))
	}
}
