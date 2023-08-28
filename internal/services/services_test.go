package services

import (
	"bytes"
	"errors"
	filework "link_shortener/internal/fileWork"
	"link_shortener/internal/shortenurl"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

type MockURLStorage struct {
	urls  map[string]string
	mutex sync.Mutex
	err   error
}

func (m *MockURLStorage) AddURL(key string, url string, userID string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.err != nil {
		return m.err
	}
	m.urls[key] = url
	return nil
}

func (m *MockURLStorage) AddURLSH(url string) (string, error) {
	shortURL := shortenurl.Shortener(url)
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.err != nil {
		return "", m.err
	}
	m.urls[shortURL] = url
	return shortURL, nil
}

func (m *MockURLStorage) GetURL(key string) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.err != nil {
		return "", m.err
	}
	url, found := m.urls[key]
	if !found {
		return "", errors.New("key not found")
	}
	return url, nil
}

func (m *MockURLStorage) GetAllURLs(userId string) ([]map[string]string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	urls := make([]map[string]string, 0)
	for shortURL, originalURL := range m.urls {
		url := make(map[string]string)
		url["short_url"] = shortURL
		url["original_url"] = originalURL
		urls = append(urls, url)
	}
	return urls, nil
}

func (m *MockURLStorage) GetOriginalURL(key string) (string, bool) {
	return key, false
}

func (m *MockURLStorage) Ping() error {
	if m.err != nil {
		return m.err
	}
	return nil
}

func TestHandleGetRequest(t *testing.T) {
	mockStorage := &MockURLStorage{
		urls: map[string]string{
			"a9b9f043": "http://example.com",
		},
	}

	req, err := http.NewRequest("GET", "/a9b9f043", nil)
	if err != nil {
		t.Fatal(err)
	}

	originURL, err := filework.FindOriginURL(conf.FileStore, "a9b9f043")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandleGetRequest(w, r, mockStorage)
	})

	rr.Header().Set("Location", originURL)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Expected %d, but got %d", http.StatusTemporaryRedirect, rr.Code)
	}

	expectedLocation := "http://example.com"
	if rr.Header().Get("Location") != expectedLocation {
		t.Errorf("Expected Location header to be %s, but got %s", expectedLocation, rr.Header().Get("Location"))
	}
}

func TestHandlePostRequest(t *testing.T) {
	mockStorage := &MockURLStorage{
		urls: make(map[string]string),
	}

	reqBody := []byte("http://example.com")

	req, err := http.NewRequest("POST", "/", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		HandlePostRequest(w, r, mockStorage, "http://localhost:8080")
	})

	dataToWrite := filework.JSONURLs{
		ShortURL:  "a9b9f043",
		OriginURL: "http://example.com",
	}

	err = filework.WriteURLsToFile(conf.FileStore, dataToWrite)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected %d, but got %d", http.StatusCreated, rr.Code)
	}

	expectedResponse := "http://localhost:8080/a9b9f043"
	if rr.Body.String() != expectedResponse {
		t.Errorf("Expected response to be %s, but got %s", expectedResponse, rr.Body.String())
	}
}

func TestPing(t *testing.T) {
	mockStorage := &MockURLStorage{}

	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Ping(mockStorage).ServeHTTP(w, r)
	})

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected %d, but got %d", http.StatusOK, rr.Code)
	}

	if rr.Body.String() != "OK" {
		t.Errorf("Expected body to be 'OK', but got '%s'", rr.Body.String())
	}
}
