package monitoring

import "github.com/prometheus/client_golang/prometheus"

var AccessCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "colors_core_access_counts",
}, []string{"status", "path", "method"})
