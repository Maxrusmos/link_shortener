package main

import (
	"flag"
	"fmt"
	config "link_shortener/internal/configs"
	"link_shortener/internal/services"
	"link_shortener/internal/storage"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	conf := config.GetConfig()
	r := chi.NewRouter()
	flag.StringVar(&conf.Address, "a", "localhost:8080", "HTTP server address")
	flag.StringVar(&conf.BaseURL, "b", "http://localhost:8080", "Base address for shortened URL")
	flag.Parse()

	fmt.Println(conf)

	storage := storage.NewMapURLStorage()

	r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {
		services.HandleGetRequest(w, r, storage)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		services.HandlePostRequest(w, r, storage, conf.BaseURL)
	})

	log.Fatal(http.ListenAndServe(conf.Address, r))
}
