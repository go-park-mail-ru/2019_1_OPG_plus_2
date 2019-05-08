package middleware

import (
	"2019_1_OPG_plus_2/internal/pkg/tsLogger"
	"net/http"
	"time"
)

func AccessLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		sw := NewStatusWriter(w)
		next.ServeHTTP(sw, r)

		tsLogger.LogAcc(
			"%d %q %s %d",
			sw.Status,
			r.Method,
			r.RequestURI,
			time.Since(start),
		)
	})
}
