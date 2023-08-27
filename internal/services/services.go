package services

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	config "link_shortener/internal/configs"
	"link_shortener/internal/shortenurl"
	"link_shortener/internal/storage"
	"log"
	"net/http"
	"strings"
	"sync"
)

var conf = config.GetConfig()

func HandleGetRequest(w http.ResponseWriter, r *http.Request, storage storage.URLStorage) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	originalURL, err := storage.GetURL(id)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", originalURL)
	log.Println(w.Header().Get("Location"))
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func HandlePostRequest(w http.ResponseWriter, r *http.Request, storage storage.URLStorage, baseURL string) {
	var mutex sync.Mutex
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

	mutex.Lock()
	defer mutex.Unlock()

	_, found := storage.GetOriginalURL(shortURL)
	if found {
		response := fmt.Sprintf("%s/%s", baseURL, shortURL)
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(response))
		return
	}

	storage.AddURL(shortURL, originalURL)

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

	existingURL, found := storage.GetOriginalURL(shortenurl.Shortener(url.URL))
	fmt.Println(existingURL)
	if found {
		response := ShortURL{Result: baseURL + "/" + shortenurl.Shortener(url.URL)}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(jsonResponse))
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

type BatchURLRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type BatchURLResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func HandleBatchShorten(w http.ResponseWriter, r *http.Request, storage storage.URLStorage, baseURL string) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requests []BatchURLRequest
	err := json.NewDecoder(r.Body).Decode(&requests)
	fmt.Println(requests)
	if err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	var responses []BatchURLResponse
	for _, req := range requests {
		shortURL, err := storage.AddURLSH(req.OriginalURL)
		if err != nil {
			http.Error(w, "Failed to add URL", http.StatusInternalServerError)
			return
		}
		response := BatchURLResponse{
			CorrelationID: req.CorrelationID,
			ShortURL:      baseURL + "/" + shortURL,
		}
		responses = append(responses, response)
	}

	jsonResponse, err := json.Marshal(responses)
	if err != nil {
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}

var users = make(map[string][]URL)

func UserUrlsHandler(w http.ResponseWriter, r *http.Request, storage storage.URLStorage) {
	// Scookie, err := cookieswork.GenerateCookie("user")
	// if err != nil {
	// 	return
	// }
	// http.SetCookie(w, Scookie)

	// Получаем куку с уникальным идентификатором пользователя
	cookie, err := r.Cookie("auth_cookie")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Проверяем подлинность куки
	parts := strings.Split(cookie.Value, "|")
	if len(parts) != 2 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	value, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	signature, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	expectedSignature := hmac.New(sha256.New, []byte("secret_key"))
	expectedSignature.Write(value)
	expectedValue := expectedSignature.Sum(nil)
	if !hmac.Equal(signature, expectedValue) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Получаем список сокращенных URL пользователя из базы данных
	urls, err := getUserUrls(cookie.Value, storage)
	if err != nil {
		fmt.Println("ggghgh")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Если список пустой, возвращаем HTTP-статус 204 No Content
	if len(urls) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Отправляем список сокращенных URL в формате JSON
	json.NewEncoder(w).Encode(urls)
	w.Header().Set("Content-Type", "application/json")
}

func getUserUrls(cookieValue string, storage storage.URLStorage) ([]map[string]string, error) {
	parts := strings.Split(cookieValue, "|")
	value, err := base64.StdEncoding.DecodeString(parts[0])
	fmt.Println(value)
	if err != nil {
		return nil, err
	}

	urls, err := storage.GetAllURLs()
	if err != nil {
		return nil, err
	}

	return urls, nil
}
