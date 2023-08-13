package main

import (
	"bytes"
	"database/sql"
	config "link_shortener/internal/configs"
	filework "link_shortener/internal/fileWork"
	"link_shortener/internal/services"
	"link_shortener/internal/shortenurl"
	"link_shortener/internal/storage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var URLMap = storage.NewMapURLStorage()

func TestHandleGetRequest(t *testing.T) {
	testDB, err := sql.Open("postgres", "user=postgres password=490Sutud dbname=link-shortener sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer testDB.Close()
	_, err = testDB.Exec("CREATE TABLE IF NOT EXISTS links (id SERIAL PRIMARY KEY, shortURL TEXT UNIQUE, originalURL TEXT)")
	if err != nil {
		t.Fatal(err)
	}

	URLMap.AddURL("a9b9f043", "http://example.com")

	req, err := http.NewRequest("GET", "/a9b9f043", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	//без флагов
	handlerNof := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		services.HandleGetRequest(w, r, URLMap, testDB, "noF")
	})
	handlerNof.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("handleGetRequest returned wrong status code: got %v want %v", rr.Code, http.StatusTemporaryRedirect)
	}

	// флаг f
	handlerF := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		services.HandleGetRequest(w, r, URLMap, testDB, "f")
	})
	handlerF.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("handleGetRequest returned wrong status code: got %v want %v", rr.Code, http.StatusTemporaryRedirect)
	}

	// флаг d
	handlerD := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		services.HandleGetRequest(w, r, URLMap, testDB, "d")
	})
	handlerD.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("handleGetRequest returned wrong status code: got %v want %v", rr.Code, http.StatusTemporaryRedirect)
	}
	///////////

	expectedLocation := "http://example.com"
	location := rr.Header().Get("Location")
	if location != expectedLocation {
		t.Errorf("handleGetRequest returned unexpected location header: got %v want %v", location, expectedLocation)
	}

	_, err = testDB.Exec("DROP TABLE links")
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandleGetRequestInvalidURL(t *testing.T) {
	testDB, err := sql.Open("postgres", "user=postgres password=490Sutud dbname=link-shortener sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer testDB.Close()
	_, err = testDB.Exec("CREATE TABLE IF NOT EXISTS links (id SERIAL PRIMARY KEY, shortURL TEXT UNIQUE, originalURL TEXT)")
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("GET", "/invalid", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	//без флага
	handlerNof := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		services.HandleGetRequest(w, r, URLMap, testDB, "noF")
	})
	handlerNof.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("handleGetRequest returned wrong status code: got %v want %v", rr.Code, http.StatusBadRequest)
	}

	// флаг f
	handlerF := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		services.HandleGetRequest(w, r, URLMap, testDB, "f")
	})
	handlerF.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("handleGetRequest returned wrong status code: got %v want %v", rr.Code, http.StatusBadRequest)
	}

	// флаг d
	handlerD := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		services.HandleGetRequest(w, r, URLMap, testDB, "d")
	})
	handlerD.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("handleGetRequest returned wrong status code: got %v want %v", rr.Code, http.StatusBadRequest)
	}

	_, err = testDB.Exec("DROP TABLE links")
	if err != nil {
		t.Fatal(err)
	}
}

func TestHandlePostRequest(t *testing.T) {
	testDB, err := sql.Open("postgres", "user=postgres password=490Sutud dbname=link-shortener sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	defer testDB.Close()

	_, err = testDB.Exec("CREATE TABLE IF NOT EXISTS links (id SERIAL PRIMARY KEY, shortURL TEXT UNIQUE, originalURL TEXT)")
	if err != nil {
		t.Fatal(err)
	}
	_, err = testDB.Exec("INSERT INTO links (shortURL, originalURL) VALUES ($1, $2)", "a9b9f043", "http://example.com")
	if err != nil {
		t.Fatal(err)
	}

	body := bytes.NewBufferString("http://example.com")

	req, err := http.NewRequest("POST", "/", body)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	var shortURL, originalURL string

	// нет флага
	handlerNof := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		services.HandlePostRequest(w, r, URLMap, config.GetConfig().BaseURL, testDB, "noF")
	})
	handlerNof.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Errorf("handlePostRequest returned wrong status code: got %v want %v", rr.Code, http.StatusCreated)
	}
	response := strings.TrimSuffix(rr.Body.String(), "\n")
	shortURL = strings.TrimPrefix(response, "http://localhost:8080/")
	originalURL, err = URLMap.GetURL(shortURL)
	if err != nil {
		t.Fatal(err)
	}

	if originalURL != "http://example.com" {
		t.Errorf("handlePostRequest added wrong original URL to map: got %v want %v", originalURL, "http://example.com")
	}

	// флаг f
	handlerF := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		services.HandlePostRequest(w, r, URLMap, config.GetConfig().BaseURL, testDB, "f")
	})
	handlerF.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Errorf("handleGetRequest returned wrong status code: got %v want %v", rr.Code, http.StatusCreated)
	}

	conf := config.GetConfig()
	urlToWrite := filework.JSONURLs{
		ShortURL:  "a9b9f043",
		OriginURL: "http://example.com",
	}
	filework.WriteURLsToFile(conf.FileStore, urlToWrite)
	originalURL, err = filework.FindOriginURL(conf.FileStore, "a9b9f043")
	if err != nil {
		t.Fatal(err)
	}
	if originalURL != "http://example.com" {
		t.Errorf("handlePostRequest added wrong original URL to map: got %v want %v", originalURL, "http://example.com")
	}

	// флаг d
	handlerD := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		services.HandlePostRequest(w, r, URLMap, config.GetConfig().BaseURL, testDB, "d")
	})
	handlerD.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Errorf("handleGetRequest returned wrong status code: got %v want %v", rr.Code, http.StatusCreated)
	}

	err = testDB.QueryRow("SELECT shortURL, originalURL FROM links WHERE shortURL='a9b9f043'").Scan(&shortURL, &originalURL)
	if err != nil {
		t.Fatal(err)
	}
	if originalURL != "http://example.com" {
		t.Errorf("handlePostRequest added wrong original URL to map: got %v want %v", originalURL, "http://example.com")
	}

	expectedResponse := "http://localhost:8080/a9b9f043"
	if response != expectedResponse {
		t.Errorf("handlePostRequest returned unexpected response body: got %v want %v", response, expectedResponse)
	}

	_, err = testDB.Exec("DROP TABLE links")
	if err != nil {
		t.Fatal(err)
	}
}

func TestShortenHandler(t *testing.T) {
	requestBody := []byte(`{"url": "http://example.com"}`)
	req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	services.ShortenHandler(rr, req, URLMap, "http://example.com")

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	expectedResponse := `{"result":"http://example.com/a9b9f043"}`
	if rr.Body.String() != expectedResponse {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expectedResponse)
	}
}

func TestShortener(t *testing.T) {
	shortURL := shortenurl.Shortener("http://example.com")

	if len(shortURL) != 8 {
		t.Errorf("shortenURL returned wrong length: got %v want %v", len(shortURL), 8)
	}
}
