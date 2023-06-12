package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"link_shortener/cmd/shortener/config"

	"github.com/go-chi/chi"
)

var urlMap = make(map[string]string)
var cfg = config.GetConfig()

func handleGetRequest(w http.ResponseWriter, r *http.Request) {
 id := strings.TrimPrefix(r.URL.Path, "/")
 originalURL, found := urlMap[id]
 if found {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", originalURL)
	w.WriteHeader(307)
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

 response := fmt.Sprintf(cfg.BaseURL, "/%s", shortURL)
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

func main() {
 r := chi.NewRouter()

 r.Get("/{id}", handleGetRequest)
 r.Post("/", handlePostRequest)
 log.Fatal(http.ListenAndServe(cfg.Address, r))
}