package main

import (
	"flag"
	config "link_shortener/internal/configs"
	"link_shortener/internal/dbwork"
	filework "link_shortener/internal/fileWork"
	"time"

	"link_shortener/internal/middleware"
	"link_shortener/internal/services"
	"link_shortener/internal/storage"
	"net/http"

	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	conf := config.GetConfig()
	r := chi.NewRouter()
	flag.StringVar(&conf.Address, "a", "localhost:8080", "HTTP server address")
	flag.StringVar(&conf.BaseURL, "b", "http://localhost:8080", "Base address for shortened URL")
	flag.StringVar(&conf.FileStore, "f", "short-url-db.json", "File storage")
	flag.StringVar(&conf.FileStore, "d", "user=postgres password=490Sutud dbname=link-shortener sslmode=disable", "db Connection String")
	flag.Parse()

	storage := storage.NewMapURLStorage()
	db, err := dbwork.Connect(conf.DBConnect)
	if err != nil {
		logger.Error("failed DB connection", zap.Error(err))
	}

	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		services.HandleGetRequest(w, r, storage)
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		services.Ping(db)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		services.HandlePostRequest(w, r, storage, conf.BaseURL)
	})
	r.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		services.ShortenHandler(w, r, storage, conf.BaseURL)
	})
	server := &http.Server{
		Addr:         conf.Address,
		Handler:      middleware.CompressionMiddleware(middleware.LoggingMiddleware(logger, r)),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	logger.Info("server started")
	if err := server.ListenAndServe(); err != nil {
		logger.Error("server stopped", zap.Error(err))
	}
	filework.ReadDataFromFile(conf.FileStore, storage)
}
