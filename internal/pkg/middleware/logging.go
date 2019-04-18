package middleware

import (
	"2019_1_OPG_plus_2/internal/pkg/tsLogger"
	"fmt"
	"net/http"
	"time"
)

func AccessLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		sw := tsLogger.NewStatusWriter(w)
		next.ServeHTTP(sw, r)

		tsLogger.Logger.LogAcc(fmt.Sprintf(
			"%s %q %s %d",
			r.Method,
			r.RequestURI,
			time.Since(start),
			sw.Status,
		))
	})
}
