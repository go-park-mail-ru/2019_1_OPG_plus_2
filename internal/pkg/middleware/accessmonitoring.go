package middleware

import (
	"2019_1_OPG_plus_2/internal/pkg/monitoring"
	"2019_1_OPG_plus_2/internal/pkg/tsLogger"
	"net/http"
	"strconv"
	"time"
)

func AccessMonitoringMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := NewStatusWriter(w)
		start := time.Now()
		next.ServeHTTP(sw, r)
		delta := time.Since(start)

		monitoring.AccessSummary.WithLabelValues(
			strconv.Itoa(sw.Status),
			r.URL.String(),
			r.Method,
		).Observe(float64(delta))

		tsLogger.LogAcc(
			"%d %q %s %d",
			sw.Status,
			r.Method,
			r.RequestURI,
			delta,
		)
	})
}
