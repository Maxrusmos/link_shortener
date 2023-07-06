package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

type Config struct {
 Address     string
 BaseAddress string
}

var urlMap = make(map[string]string)

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

 response := fmt.Sprintf("%s/%s", config.BaseAddress, shortURL)
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

var config Config

func init() {
//  flag.StringVar(&config.Address, "a", "localhost:8080", "HTTP server address")
//  flag.StringVar(&config.BaseAddress, "b", "http://localhost:8080", "Base address for shortened URL")
//  flag.Parse()
}

func main() {
 r := chi.NewRouter()

 r.Get("/{id}", handleGetRequest)
 r.Post("/", handlePostRequest)

 log.Fatal(http.ListenAndServe(config.Address, r))
}