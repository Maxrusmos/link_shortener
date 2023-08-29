package startprep

import (
	config "link_shortener/internal/configs"
	"link_shortener/internal/dbwork"
	"link_shortener/internal/middleware"
	"link_shortener/internal/services"
	"link_shortener/internal/storage"
	"net/http"
	"os"
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
	r.Get("/api/user/urls", func(w http.ResponseWriter, r *http.Request) {
		services.UserUrlsHandler(w, r, storageURL)
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
	r.Delete("/api/user/urls", func(w http.ResponseWriter, r *http.Request) {
		services.DeleteURLsHandler(w, r, storageURL)
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

func GetStorageURL(conf config.Config) storage.URLStorage {
	if conf.DBConnect == "" && conf.FileStore == "" && os.Getenv("FILE_STORAGE_PATH") == "" {
		return storage.NewMapURLStorage()
	}

	if conf.FileStore != "" || os.Getenv("FILE_STORAGE_PATH") != "" {
		return storage.NewFileURLStorage(conf.FileStore)
	}

	db, err := dbwork.Connect(conf.DBConnect)
	if err != nil {
		panic(err)
	}
	err = dbwork.CreateTables(db, `CREATE TABLE IF NOT EXISTS shortened_urls  (
        id SERIAL PRIMARY KEY,
    	short_url TEXT NOT NULL,
    	original_url TEXT NOT NULL,
   	 	user_id TEXT,
		deleted_flag BOOLEAN DEFAULT false,
    	UNIQUE (original_url)
      )`)
	if err != nil {
		panic(err)
	}
	return storage.NewDatabaseURLStorage(db)
}
