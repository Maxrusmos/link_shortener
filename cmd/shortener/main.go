package main

import (
	"flag"
	"fmt"
	"io"
	config "link_shortener/internal/configs"
	"link_shortener/internal/data"
	"link_shortener/internal/shortenurl"
	"strings"

	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

type Config struct {
	ServerAddr string `env:"SERVER_ADDRESS"`
	BaseURL    string `env:"BASE_URL"`
}

func HandleGetRequest(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	originalURL, found := data.URLMap[id]
	if found {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
	}
}

func HandlePostRequest(w http.ResponseWriter, r *http.Request) {
	conf := config.GetConfig()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	originalURL := string(body)
	shortURL := shortenurl.Shortener(originalURL)
	data.URLMap[shortURL] = originalURL

	response := fmt.Sprintf("%s/%s", conf.BaseURL, shortURL)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
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
	r.Post("/", HandlePostRequest)
	r.Get("/{id}", HandleGetRequest)
	log.Fatal(http.ListenAndServe(conf.Address, r))
}
