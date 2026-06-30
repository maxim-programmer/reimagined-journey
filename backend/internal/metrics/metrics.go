package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	SearchRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "search_requests_total",
		Help: "Total number of requests to /api/v1/search",
	})

	SearchRequestDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "search_request_duration_seconds",
		Help:    "Duration of /api/v1/search requests in seconds",
		Buckets: prometheus.DefBuckets,
	})

	SearchRequestErrors = promauto.NewCounter(prometheus.CounterOpts{
		Name: "search_request_errors_total",
		Help: "Total number of failed /api/v1/search requests",
	})
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func InstrumentSearch(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		SearchRequestsTotal.Inc()

		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)

		SearchRequestDuration.Observe(time.Since(start).Seconds())
		if rec.status >= 400 {
			SearchRequestErrors.Inc()
		}
	}
}

func Handler() http.Handler {
	return promhttp.Handler()
}