package middleware

import (
	"2019_1_OPG_plus_2/internal/pkg/monitoring"
	"net/http"
	"strconv"
)

func AccessMonitoringMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := NewStatusWriter(w)
		next.ServeHTTP(sw, r)
		monitoring.AccessCounter.WithLabelValues(
			strconv.Itoa(sw.Status),
			r.URL.String(),
			r.Method).Inc()
	})
}
