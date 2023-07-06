package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleGetRequest(t *testing.T) {
	urlMap["test"] = "http://example.com"

	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handleGetRequest)

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
	handler := http.HandlerFunc(handleGetRequest)

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
	handler := http.HandlerFunc(handlePostRequest)

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
	originalURL, found := urlMap[shortURL]
	if !found {
		t.Errorf("handlePostRequest did not add short URL to map")
	}

	if originalURL != "http://example.com" {
		t.Errorf("handlePostRequest added wrong original URL to map: got %v want %v", originalURL, "http://example.com")
	}
}

func TestShortenURL(t *testing.T) {
	shortURL := shortenURL("http://example.com")

	if len(shortURL) != 8 {
		t.Errorf("shortenURL returned wrong length: got %v want %v", len(shortURL), 8)
	}
}
