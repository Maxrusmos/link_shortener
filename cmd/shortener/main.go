package main

import (
	"flag"
	config "link_shortener/internal/configs"
	"link_shortener/internal/services"
	"link_shortener/internal/storage"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

func main() {
	conf := config.GetConfig()

	if os.Getenv(conf.ServerAddrENV) != "" {
		conf.ServerAddrENV = os.Getenv(conf.ServerAddrENV)
		conf.Address = conf.ServerAddrENV
	}
	if os.Getenv(conf.BaseURLENV) != "" {
		conf.BaseURLENV = os.Getenv(conf.BaseURLENV)
		conf.BaseURL = conf.BaseURLENV
	}

	r := chi.NewRouter()
	flag.StringVar(&conf.Address, "a", "localhost:8080", "HTTP server address")
	flag.StringVar(&conf.BaseURL, "b", "http://localhost:8080", "Base address for shortened URL")
	flag.Parse()

	storage := storage.NewMapURLStorage()

	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		services.HandleGetRequest(w, r, storage)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		services.HandlePostRequest(w, r, storage)
	})

	log.Fatal(http.ListenAndServe(conf.Address, r))
}
