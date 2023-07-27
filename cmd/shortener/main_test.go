package main

import (
	"bytes"
	"fmt"
	config "link_shortener/internal/configs"
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
	URLMap.AddURL("test", "http://example.com")

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		services.HandleGetRequest(w, r, URLMap)
	})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("handleGetRequest returned wrong status code: got %v want %v", rr.Code, http.StatusTemporaryRedirect)
	}

	expectedLocation := "http://example.com"
	location := rr.Header().Get("Location")
	if location != expectedLocation {
		t.Errorf("handleGetRequest returned unexpected location header: got %v want %v", location, expectedLocation)
	}
}

func TestHandleGetRequestInvalidURL(t *testing.T) {
	req, err := http.NewRequest("GET", "/invalid", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		services.HandleGetRequest(w, r, URLMap)
	})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("handleGetRequest returned wrong status code: got %v want %v", rr.Code, http.StatusBadRequest)
	}
}

func TestHandlePostRequest(t *testing.T) {
	body := bytes.NewBufferString("http://example.com")

	req, err := http.NewRequest("POST", "/", body)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		services.HandlePostRequest(w, r, URLMap, config.GetConfig().BaseURL)
	})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("handlePostRequest returned wrong status code: got %v want %v", rr.Code, http.StatusCreated)
	}

	expectedResponse := "http://localhost:8080/a9b9f043"
	response := strings.TrimSuffix(rr.Body.String(), "\n")
	if response != expectedResponse {
		t.Errorf("handlePostRequest returned unexpected response body: got %v want %v", response, expectedResponse)
	}

	shortURL := strings.TrimPrefix(response, "http://localhost:8080/")
	originalURL, er := URLMap.GetURL(shortURL)
	fmt.Println(er)

	if originalURL != "http://example.com" {
		t.Errorf("handlePostRequest added wrong original URL to map: got %v want %v", originalURL, "http://example.com")
	}
}

func TestShortenHandler(t *testing.T) {

	// Create a test request with a JSON body
	requestBody := []byte(`{"url": "http://example.com"}`)
	req, err := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	// Create a test response recorder
	rr := httptest.NewRecorder()

	// Call the ShortenHandler function with the mock storage and test request/response
	services.ShortenHandler(rr, req, URLMap, "http://example.com")

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body
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
