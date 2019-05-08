package middleware

import (
	"2019_1_OPG_plus_2/internal/pkg/monitoring"
	"net/http"
	"strconv"
)

func AccessMonitoringMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := NewStatusWriter(w)
		next.ServeHTTP(w, r)

		monitoring.AccessCounter.WithLabelValues(
			strconv.FormatInt(int64(sw.Status), 10),
			r.URL.String(),
			r.Method).Inc()
	})
}
