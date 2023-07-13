package services

import (
	"fmt"
	"io"
	config "link_shortener/internal/configs"
	"link_shortener/internal/storage"
	"net/http"
	"strings"
)

var conf = config.GetConfig()

func HandleGetRequest(w http.ResponseWriter, r *http.Request, storage storage.URLStorage) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	originalURL, err := storage.GetURL(id)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func HandlePostRequest(w http.ResponseWriter, r *http.Request, storage storage.URLStorage) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	originalURL := string(body)
	shortURL, err := storage.AddURLSH(originalURL)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("%s/%s", conf.BaseURL, shortURL)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}
