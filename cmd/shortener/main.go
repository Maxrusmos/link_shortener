package main

import (
	"flag"
	config "link_shortener/internal/configs"
	"link_shortener/internal/services"

	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

type Config struct {
	ServerAddr string `env:"SERVER_ADDRESS"`
	BaseURL    string `env:"BASE_URL"`
}

func main() {
	var cfg Config
	conf := config.GetConfig()
	if os.Getenv("SERVER_ADDRESS") != "" {
		cfg.ServerAddr = os.Getenv("SERVER_ADDRESS")
		conf.Address = cfg.ServerAddr
	}
	if os.Getenv("BASE_URL") != "" {
		cfg.BaseURL = os.Getenv("BASE_URL")
		conf.BaseURL = cfg.BaseURL
	}
	r := chi.NewRouter()
	flag.StringVar(&conf.Address, "a", "localhost:8080", "адрес")
	flag.StringVar(&conf.BaseURL, "b", "http://localhost:8080", "базовый URL")
	flag.Parse()
	r.Post("/", services.HandlePostRequest)
	r.Get("/{id}", services.HandleGetRequest)
	log.Fatal(http.ListenAndServe(conf.Address, r))
}
