package main

import (
	"bytes"
	"flag"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAFlag(t *testing.T) {
    // Создаем фейковый аргумент командной строки с флагом -a
    flag.Set("a", "localhost:8080")
   
    // Получаем значение флага -a
    a := flag.String("a", "localhost:8080", "Адрес программы")
   
    // Проверяем, что значение флага -a соответствует ожидаемому значению
    if *a != "localhost:8080" {
     t.Errorf("Ожидаемое значение: localhost:8080, получено значение: %s", *a)
    }
   }
   
   func TestBFlag(t *testing.T) {
    // Создаем фейковый аргумент командной строки с флагом -b
    flag.Set("b", "http://localhost:8080")
   
    // Получаем значение флага -b
    b := flag.String("b", "http://localhost:8080/a9b9f043", "Базовый адрес для сокращенных URL")
   
    // Проверяем, что значение флага -b соответствует ожидаемому значению
    if *b != "http://localhost:8080/a9b9f043" {
     t.Errorf("http://localhost:8080/a9b9f043, получено значение: %s", *b)
    }
   }
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
