package services

import (
	"encoding/json"
	"fmt"
	"io"
	config "link_shortener/internal/configs"
	"link_shortener/internal/shortenurl"
	"link_shortener/internal/storage"
	"log"
	"net/http"
	"strings"
)

var conf = config.GetConfig()

func HandleGetRequest(w http.ResponseWriter, r *http.Request, storage storage.URLStorage) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	var originalURL string
	var err error

	originalURL, err = storage.GetURL(id)
	log.Println("originalURL is", originalURL)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func HandlePostRequest(w http.ResponseWriter, r *http.Request, storage storage.URLStorage, baseURL string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	originalURL := strings.TrimSpace(string(body))
	shortURL := shortenurl.Shortener(originalURL)

	err = storage.AddURL(shortURL, originalURL)
	if err != nil {
		http.Error(w, "Failed to add URL ghgsdghsghdhsdhshdgh", http.StatusInternalServerError)
		return
	}

	response := fmt.Sprintf("%s/%s", baseURL, shortURL)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

func Ping(storage storage.URLStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := storage.Ping(); err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, "OK")
	}
}

type URL struct {
	URL string `json:"url"`
}

type ShortURL struct {
	Result string `json:"result"`
}

func ShortenHandler(w http.ResponseWriter, r *http.Request, storage storage.URLStorage, baseURL string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var url URL
	err = json.Unmarshal(body, &url)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	shortURL, err := storage.AddURLSH(url.URL)
	if err != nil {
		http.Error(w, "Failed to add URL", http.StatusInternalServerError)
		return
	}

	response := ShortURL{Result: baseURL + "/" + shortURL}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}
