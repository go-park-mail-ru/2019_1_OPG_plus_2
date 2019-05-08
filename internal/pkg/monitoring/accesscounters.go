package monitoring

import "github.com/prometheus/client_golang/prometheus"

var AccessSummary = prometheus.NewSummaryVec(prometheus.SummaryOpts{
	Name:      "access_time_summaries",
	Namespace: "colors_core",
}, []string{"status", "path", "method"})

var ActiveRooms = prometheus.NewGauge(prometheus.GaugeOpts{
	Name:      "active_rooms",
	Namespace: "colors_game",
})

var ActiveConns = prometheus.NewGauge(prometheus.GaugeOpts{
	Name:      "active_cons",
	Namespace: "colors_game",
})
