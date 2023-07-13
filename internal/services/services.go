package services

import (
	"fmt"
	"io"
	config "link_shortener/internal/configs"
	"link_shortener/internal/storage"
	"net/http"
	"strings"
)

func HandleGetRequest(w http.ResponseWriter, r *http.Request, storage storage.URLStorage) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	fmt.Println(storage)
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
	conf := config.GetConfig()
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

	fmt.Println(conf.BaseURL)
	response := fmt.Sprintf("%s/%s", conf.BaseURL, shortURL)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}
