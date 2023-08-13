package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestLoggingMiddleware(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	})

	req := httptest.NewRequest("GET", "/test", nil)
	recorder := httptest.NewRecorder()

	middleware := LoggingMiddleware(logger, handler)
	middleware.ServeHTTP(recorder, req)

	resp := recorder.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Hello, World!", string(body))
}

func TestCompressionMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		w.Write(body)
	})

	payload := "Hello, World!"
	gzipBuffer := new(bytes.Buffer)
	gzipWriter := gzip.NewWriter(gzipBuffer)
	gzipWriter.Write([]byte(payload))
	gzipWriter.Close()

	req := httptest.NewRequest("POST", "/test", gzipBuffer)
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Content-Length", string(rune(len(gzipBuffer.Bytes()))))

	recorder := httptest.NewRecorder()

	middleware := CompressionMiddleware(handler)
	middleware.ServeHTTP(recorder, req)

	resp := recorder.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, payload, string(body))
}
