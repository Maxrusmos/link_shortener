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
	"net/url"
	"regexp"
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

func isValidURL(u string) bool {
	parsedURL, err := url.Parse(u)
	fmt.Println("parse", parsedURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return false
	}
	return true
}

func RemoveControlCharacters(input string) string {
	// Создаем регулярное выражение для поиска управляющих символов
	controlCharRegex := regexp.MustCompile(`[[:cntrl:]]`)

	// Заменяем управляющие символы на пустую строку
	cleanedString := controlCharRegex.ReplaceAllString(input, "")

	return cleanedString
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

	originalURL := RemoveControlCharacters(strings.TrimSpace(string(body)))
	if !isValidURL(originalURL) {
		http.Error(w, "Invalid URL hhhhhhhhhhhhhhhhhhhhhhhhhhhhh", http.StatusBadRequest)
		return
	}
	// shortURL := shortenurl.Shortener(originalURL)
	var shortURL string

	shortURL, err = storage.AddURLSH(originalURL)
	fmt.Println("shortURL after ADD:::", shortURL)
	if err != nil {
		http.Error(w, "Failed to add URL", http.StatusInternalServerError) // Обработка ошибки добавления URL
		return
	}

	response := fmt.Sprintf("%s/%s", baseURL, shortURL)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated) // Возвращаем корректный статус для успешного добавления
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

	var url URL
	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortURL := shortenurl.Shortener(url.URL)
	err = storage.AddURL(shortURL, url.URL)
	if err != nil {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	response := ShortURL{Result: baseURL + "/" + shortURL}
	log.Println("ЬГВлитыповмолымв")
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}
