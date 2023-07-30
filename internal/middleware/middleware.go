package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

func LoggingMiddleware(logger *zap.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		responseWriter := &responseWriter{ResponseWriter: w, status: http.StatusOK, size: 0}
		defer func() {
			elapsed := time.Since(start)
			logger.Info("response",
				zap.String("uri", r.RequestURI),
				zap.String("method", r.Method),
				zap.Int("status", responseWriter.status),
				zap.Int64("size", responseWriter.size),
				zap.Duration("elapsed", elapsed),
			)
		}()

		next.ServeHTTP(responseWriter, r)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int64
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += int64(size)
	return size, err
}

//////////////////////////////////////////////////////////////////////////////////

func CompressionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			reader, err := gzip.NewReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Error reading compressed request body"))
				return
			}
			defer reader.Close()
			body, err := io.ReadAll(reader)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Error decompressing request body"))
				return
			}
			r.Body = io.NopCloser(strings.NewReader(string(body)))
			r.Header.Del("Content-Encoding")
			r.Header.Set("Content-Length", string(rune(len(body))))
		}
		next.ServeHTTP(w, r)
	})
}