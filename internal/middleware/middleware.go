package middleware

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

func JSONContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			next.ServeHTTP(w, r)

			logger.Infof(
				"Method: %s, Path: %s, RemoteAddr: %s, UserAgent: %s Duration: %s ",
				r.Method,
				r.URL.Path,
				r.UserAgent(),
				r.RemoteAddr,
				time.Since(start),
			)
		})
	}
}
