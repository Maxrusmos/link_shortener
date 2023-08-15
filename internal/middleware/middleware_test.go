package middleware

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
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
	resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, "Hello, World!", string(body))
}

func TestCompressionMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
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
	resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, payload, string(body))
}

func compress(s string) string {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write([]byte(s))
	w.Close()
	return b.String()
}

func readResponseBody(resp *http.Response) string {
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}

type errorReader struct{}

func (er *errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("error reading request body")
}
