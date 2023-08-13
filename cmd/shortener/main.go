package main

import (
	"database/sql"
	"flag"
	"fmt"
	config "link_shortener/internal/configs"
	"link_shortener/internal/dbwork"
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
	flag.StringVar(&conf.FileStore, "f", "", "File storage")
	// short-url-db.json
	flag.StringVar(&conf.DBConnect, "d", "", "db Connection String")
	// user=postgres password=490Sutud dbname=link-shortener sslmode=disable
	flag.Parse()

	var flagProvided string
	storage := storage.NewMapURLStorage()

	if conf.DBConnect == "" {
		if conf.FileStore == "" {
			flagProvided = "noF"
		} else {
			flagProvided = "f"
		}
	} else {
		flagProvided = "d"
	}
	var db *sql.DB

	switch flagProvided {
	case "f":
		fmt.Println("f")
	case "d":
		db, err = dbwork.Connect(conf.DBConnect)
		if err != nil {
			logger.Error("failed DB connection", zap.Error(err))
		}
		fmt.Println("d")
	case "noF":
		fmt.Println("noF")
	}

	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		services.HandleGetRequest(w, r, storage, db, flagProvided)
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		services.Ping(db)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		services.HandlePostRequest(w, r, storage, conf.BaseURL, db, flagProvided)
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
