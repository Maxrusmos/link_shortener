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
	flag.StringVar(&conf.Address, "a", "localhost:8080", "адрес")
	flag.StringVar(&conf.BaseURL, "b", "http://localhost:8080", "базовый URL")
	flag.Parse()
	r.Post("/", services.HandlePostRequest)
	r.Get("/{id}", services.HandleGetRequest)
	log.Fatal(http.ListenAndServe(conf.Address, r))
}
