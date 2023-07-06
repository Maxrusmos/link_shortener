package main

import (
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"link_shortener/cmd/shortener/config"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
)

var urlMap = make(map[string]string)
var conf = config.GetConfig()

func handleGetRequest(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	originalURL, found := urlMap[id]
	if found {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Location", originalURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
	}
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	originalURL := string(body)
	shortURL := shortenURL(originalURL)
	urlMap[shortURL] = originalURL

	response := fmt.Sprintf("%s/%s", conf.BaseURL, shortURL)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

func shortenURL(originalURL string) string {
	hasher := md5.New()
	hasher.Write([]byte(originalURL))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash[:8]
}

type Config struct {
	ServerAddr string `env:"SERVER_ADDRESS"`
	BaseURL    string `env:"BASE_URL"`
}

func main() {
	var cfg Config
	if os.Getenv("SERVER_ADDRESS") != "" {
		cfg.ServerAddr = os.Getenv("SERVER_ADDRESS")
	} else {
		flag.StringVar(&conf.Address, "a", "localhost:8080", "HTTP server address")
	}
	if os.Getenv("BASE_URL") != "" {
		cfg.BaseURL = os.Getenv("BASE_URL")
	} else {
		flag.StringVar(&conf.BaseURL, "b", "http://localhost:8080", "Base address for shortened URL")
	}
	r := chi.NewRouter()

	flag.Parse()
	r.Get("/{id}", handleGetRequest)
	r.Post("/", handlePostRequest)

	log.Fatal(http.ListenAndServe(conf.Address, r))
}
