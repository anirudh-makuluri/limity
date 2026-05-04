package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	requestsTotal           *prometheus.CounterVec
	requestDuration         *prometheus.HistogramVec
	checkTotal              *prometheus.CounterVec
	analyticsDroppedTotal   prometheus.Counter
	apiKeysTotal            prometheus.Gauge
	redisErrorsTotal        prometheus.Counter
	authFailuresTotal       prometheus.Counter
	ownerLookupSuccessTotal prometheus.Counter
	ownerLookupMissTotal    prometheus.Counter
	ownerLookupErrorTotal   prometheus.Counter
}

func NewMetrics() *Metrics {
	m := &Metrics{
		requestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{Name: "http_requests_total", Help: "Total HTTP requests."},
			[]string{"route", "method", "status"},
		),
		requestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request latency in seconds.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"route", "method"},
		),
		checkTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{Name: "limity_check_total", Help: "Total /check results."},
			[]string{"result"},
		),
		analyticsDroppedTotal: prometheus.NewCounter(
			prometheus.CounterOpts{Name: "analytics_events_dropped_total", Help: "Dropped analytics events due to backpressure/errors."},
		),
		apiKeysTotal: prometheus.NewGauge(
			prometheus.GaugeOpts{Name: "limity_api_keys_total", Help: "Total API keys in users table."},
		),
		redisErrorsTotal: prometheus.NewCounter(
			prometheus.CounterOpts{Name: "limity_redis_errors_total", Help: "Total Redis operation errors."},
		),
		authFailuresTotal: prometheus.NewCounter(
			prometheus.CounterOpts{Name: "limity_auth_failures_total", Help: "Total auth failures."},
		),
		ownerLookupSuccessTotal: prometheus.NewCounter(
			prometheus.CounterOpts{Name: "limity_owner_lookup_success_total", Help: "Total successful owner lookups by API key."},
		),
		ownerLookupMissTotal: prometheus.NewCounter(
			prometheus.CounterOpts{Name: "limity_owner_lookup_miss_total", Help: "Total owner lookup misses by API key."},
		),
		ownerLookupErrorTotal: prometheus.NewCounter(
			prometheus.CounterOpts{Name: "limity_owner_lookup_error_total", Help: "Total owner lookup errors by API key."},
		),
	}

	collectors := []prometheus.Collector{
		m.requestsTotal,
		m.requestDuration,
		m.checkTotal,
		m.analyticsDroppedTotal,
		m.apiKeysTotal,
		m.redisErrorsTotal,
		m.authFailuresTotal,
		m.ownerLookupSuccessTotal,
		m.ownerLookupMissTotal,
		m.ownerLookupErrorTotal,
	}
	for _, c := range collectors {
		_ = prometheus.Register(c)
	}
	return m
}

func (s *Server) RegisterObservabilityRoutes(r chi.Router) {
	r.Handle("/metrics", promhttp.Handler())
}

func (s *Server) ObservabilityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(withRequestMeta(r.Context()))
		start := time.Now()
		rw := &statusWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rw, r)

		if s.metrics == nil {
			return
		}

		route := routePattern(r)
		method := r.Method
		status := strconv.Itoa(rw.statusCode)
		duration := time.Since(start).Seconds()

		s.metrics.requestsTotal.WithLabelValues(route, method, status).Inc()
		s.metrics.requestDuration.WithLabelValues(route, method).Observe(duration)

		if s.analytics != nil && r.URL.Path != "/health" {
			enqueued := s.analytics.Enqueue(RequestEvent{
				Timestamp:      time.Now().UTC(),
				Method:         method,
				Route:          route,
				Path:           r.URL.Path,
				StatusCode:     rw.statusCode,
				DurationMs:     time.Since(start).Milliseconds(),
				ClientIP:       r.RemoteAddr,
				UserAgent:      r.UserAgent(),
				OwnerUserID:    getOwnerUserID(r.Context()),
				APIKey:         getAPIKeyFromCheckRequest(r.Context()),
			})
			if !enqueued {
				s.metrics.analyticsDroppedTotal.Inc()
			}
		}
	})
}

func (s *Server) StartAPIKeyGaugeUpdater(interval time.Duration) {
	if s.metrics == nil || s.pg == nil || interval <= 0 {
		return
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	update := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		count, err := s.pg.CountAPIKeys(ctx)
		cancel()
		if err != nil {
			return
		}
		s.metrics.apiKeysTotal.Set(float64(count))
	}

	update()
	for range ticker.C {
		update()
	}
}

type statusWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *statusWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func routePattern(r *http.Request) string {
	if rc := chi.RouteContext(r.Context()); rc != nil {
		if pattern := rc.RoutePattern(); pattern != "" {
			return pattern
		}
	}
	return r.URL.Path
}
