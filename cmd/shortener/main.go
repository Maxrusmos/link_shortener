package main

import (
	"flag"
	"fmt"
	config "link_shortener/internal/configs"
	"link_shortener/internal/dbwork"
	"link_shortener/internal/middleware"
	"link_shortener/internal/services"
	"link_shortener/internal/storage"
	"net/http"
	"time"

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
	flag.StringVar(&conf.DBConnect, "d", "", "db Connection String")
	// user=postgres password=490Sutud dbname=link-shortener sslmode=disable
	flag.Parse()

	var storageURL storage.URLStorage

	if conf.DBConnect == "" {
		if conf.FileStore != "" {
			storageURL = storage.NewFileURLStorage(conf.FileStore)
		} else {
			storageURL = storage.NewMapURLStorage()
		}
	} else {
		db, err := dbwork.Connect(conf.DBConnect)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer db.Close()
		storageURL = storage.NewDatabaseURLStorage(db)
		err = dbwork.CreateTables(db, `CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			shortURL TEXT UNIQUE,
			originalURL TEXT
		  )`)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		services.HandleGetRequest(w, r, storageURL)
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		services.Ping(storageURL)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		services.HandlePostRequest(w, r, storageURL, conf.BaseURL)
	})
	r.Post("/api/shorten", func(w http.ResponseWriter, r *http.Request) {
		services.ShortenHandler(w, r, storageURL, conf.BaseURL)
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
