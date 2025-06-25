// infra/observability/metrics_recorder.go
package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricsRecorder interface {
	RecordRoleCacheInitSuccess()
	RecordRoleCacheRefreshSuccess()
	RecordRoleCacheRefreshFailed()
	RecordRoleCacheLoadCount(count float64)
	RecordRoleCacheHit()
	RecordRoleCacheMiss(reason string)
	RecordGrpcRequestTotal(fullMethod string, statusCode string)
	RecordGrpcRequestDuration(fullMethod string, durationSeconds float64)
}

type prometheusMetricsRecorder struct {
	roleCacheInitSuccessCounter    prometheus.Counter
	roleCacheRefreshSuccessCounter prometheus.Counter
	roleCacheRefreshFailedCounter  prometheus.Counter
	roleCacheLoadCountGauge        prometheus.Gauge
	roleCacheHitCounter            prometheus.Counter
	roleCacheMissCounter           *prometheus.CounterVec
	grpcRequestsTotal              *prometheus.CounterVec
	grpcRequestDuration            *prometheus.HistogramVec
}

func NewPrometheusMetricsRecorder() MetricsRecorder {
	return &prometheusMetricsRecorder{
		roleCacheInitSuccessCounter: promauto.NewCounter(prometheus.CounterOpts{
			Name: "app_role_cache_init_success_total",
			Help: "Total number of successful role cache initializations.",
		}),
		roleCacheRefreshSuccessCounter: promauto.NewCounter(prometheus.CounterOpts{
			Name: "app_role_cache_refresh_success_total",
			Help: "Total number of successful role cache refreshes.",
		}),
		roleCacheRefreshFailedCounter: promauto.NewCounter(prometheus.CounterOpts{
			Name: "app_role_cache_refresh_failed_total",
			Help: "Total number of failed role cache refreshes.",
		}),
		roleCacheLoadCountGauge: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "app_role_cache_load_total_roles",
			Help: "Total number of roles loaded into the cache.",
		}),
		roleCacheHitCounter: promauto.NewCounter(prometheus.CounterOpts{
			Name: "app_role_cache_hit_total",
			Help: "Total number of role cache hits.",
		}),
		roleCacheMissCounter: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "app_role_cache_miss_total",
			Help: "Total number of role cache misses by reason.",
		}, []string{"reason"}),
		grpcRequestsTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "grpc_server_requests_total",
			Help: "Total number of gRPC requests by method and status code.",
		}, []string{"method", "code"}),
		grpcRequestDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "grpc_server_request_duration_seconds",
			Help:    "Histogram of gRPC server request latencies.",
			Buckets: prometheus.DefBuckets,
		}, []string{"method"}),
	}
}

func (r *prometheusMetricsRecorder) RecordRoleCacheInitSuccess() {
	r.roleCacheInitSuccessCounter.Inc()
}
func (r *prometheusMetricsRecorder) RecordRoleCacheRefreshSuccess() {
	r.roleCacheRefreshSuccessCounter.Inc()
}
func (r *prometheusMetricsRecorder) RecordRoleCacheRefreshFailed() {
	r.roleCacheRefreshFailedCounter.Inc()
}
func (r *prometheusMetricsRecorder) RecordRoleCacheLoadCount(count float64) {
	r.roleCacheLoadCountGauge.Set(count)
}
func (r *prometheusMetricsRecorder) RecordRoleCacheHit() {
	r.roleCacheHitCounter.Inc()
}
func (r *prometheusMetricsRecorder) RecordRoleCacheMiss(reason string) {
	r.roleCacheMissCounter.WithLabelValues(reason).Inc()
}
func (r *prometheusMetricsRecorder) RecordGrpcRequestTotal(fullMethod string, statusCode string) {
	r.grpcRequestsTotal.WithLabelValues(fullMethod, statusCode).Inc()
}
func (r *prometheusMetricsRecorder) RecordGrpcRequestDuration(fullMethod string, durationSeconds float64) {
	r.grpcRequestDuration.WithLabelValues(fullMethod).Observe(durationSeconds)
}
