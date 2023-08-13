package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	config "link_shortener/internal/configs"
	"link_shortener/internal/dbwork"
	filework "link_shortener/internal/fileWork"
	"link_shortener/internal/shortenurl"
	"link_shortener/internal/storage"
	"net/http"
	"strings"
)

func HandleGetRequest(w http.ResponseWriter, r *http.Request, storage storage.URLStorage, db *sql.DB, flagProvided string) {
	id := strings.TrimPrefix(r.URL.Path, "/")
	var originalURL string
	var err error

	if flagProvided == "d" {
		originalURL, err = dbwork.GetOriginalURL(db, id)
		fmt.Println(originalURL)
		if err != nil {
			fmt.Print("err")
		} else {
			fmt.Println(originalURL)
		}
	} else {
		if flagProvided == "f" {
			conf := config.GetConfig()
			originalURL, err = filework.FindOriginURL(conf.FileStore, id)
			if err != nil {
				fmt.Println("err")
			}
		} else {
			originalURL, err = storage.GetURL(id)
			if err != nil {
				http.Error(w, "Invalid URL", http.StatusBadRequest)
				return
			}
		}

	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func HandlePostRequest(w http.ResponseWriter, r *http.Request, storage storage.URLStorage, baseURL string, db *sql.DB, flagProvided string) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var shortURL string

	originalURL := strings.ReplaceAll(string(body), "\"", "")
	shortURL = shortenurl.Shortener(originalURL)

	if flagProvided == "d" {
		dbwork.CreateTables(db)
		dbwork.AddURL(db, shortURL, originalURL)
	} else {
		if flagProvided == "f" {
			conf := config.GetConfig()
			urlToWrite := filework.JSONURLs{
				ShortURL:  shortURL,
				OriginURL: originalURL,
			}
			filework.WriteURLsToFile(conf.FileStore, urlToWrite)
		} else {
			shortURL, err = storage.AddURLSH(originalURL)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

	}

	response := fmt.Sprintf("%s/%s", baseURL, shortURL)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(response))
}

func Ping(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := db.Ping(); err != nil {
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
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var url URL
	err := json.NewDecoder(r.Body).Decode(&url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	shortURL, err := storage.AddURLSH(url.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := ShortURL{Result: baseURL + "/" + shortURL}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}
