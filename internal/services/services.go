package services

import (
	"fmt"
	"io"
	config "link_shortener/internal/configs"
	"link_shortener/internal/data"
	"link_shortener/internal/shortenurl"
	"net/http"
	"strings"
)

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
