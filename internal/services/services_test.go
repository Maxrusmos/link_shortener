package services

import (
	config "link_shortener/internal/configs"
	"link_shortener/internal/storage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert" // Используем библиотеку testify для утверждений
)

var URLMap = storage.NewMapURLStorage()

func TestHandleGetRequest(t *testing.T) {
	URLMap.AddURL("a9b9f043", "http://example.com")
	req, err := http.NewRequest("GET", "/a9b9f043", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	HandleGetRequest(rr, req, URLMap)

	assert.Equal(t, http.StatusTemporaryRedirect, rr.Code)
	assert.Equal(t, "http://example.com", rr.Header().Get("Location"))
}

func TestHandlePostRequest(t *testing.T) {
	body := strings.NewReader(`"https://www.example.com"`)
	req, err := http.NewRequest("POST", "/shorten", body)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	HandlePostRequest(rr, req, URLMap, config.GetConfig().BaseURL)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Contains(t, rr.Body.String(), "http://localhost/abc123")
}
