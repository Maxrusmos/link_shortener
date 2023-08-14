package main

import (
	"database/sql"
	"flag"
	"fmt"
	config "link_shortener/internal/configs"
	"link_shortener/internal/flagpkg"
	"link_shortener/internal/storage"
	"time"

	"link_shortener/internal/middleware"
	"link_shortener/internal/services"
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
	flag.StringVar(&conf.DBConnect, "d", "user=postgres password=490Sutud dbname=link-shortener sslmode=disable", "db Connection String")
	// user=postgres password=490Sutud dbname=link-shortener sslmode=disable
	flag.Parse()

	storage := storage.NewMapURLStorage()
	flag := flagpkg.GetSharedFlag()
	if conf.DBConnect == "" {
		if conf.FileStore != "" {
			flag.SetValue("f")
		} else {
			flag.SetValue("noF")
		}
	} else {
		flag.SetValue("d")
	}
	fmt.Println(flag.GetValue())
	var db *sql.DB

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
}
