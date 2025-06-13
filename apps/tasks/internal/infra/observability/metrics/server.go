package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/pratchaya-maneechot/service-exchange/apps/tasks/internal/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricServer struct {
	server *http.Server
	config config.MetricsConfig
}

var (
	GrpcRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
		[]string{"method", "status"},
	)

	GrpcRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "Duration of gRPC requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)

	GrpcActiveConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "grpc_active_connections",
			Help: "Number of active gRPC connections",
		},
	)

	AppInfo = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "app_info",
			Help: "Application information",
		},
		[]string{"version", "environment"},
	)
)

func init() {
	prometheus.MustRegister(
		GrpcRequestsTotal,
		GrpcRequestDuration,
		GrpcActiveConnections,
		AppInfo,
	)
}

func NewServer(cfg config.MetricsConfig) *MetricServer {
	mux := http.NewServeMux()
	mux.Handle(cfg.Path, promhttp.Handler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return &MetricServer{
		server: &http.Server{
			Addr:         cfg.Address,
			Handler:      mux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		config: cfg,
	}
}

func (s *MetricServer) Start() error {
	if !s.config.Enabled {
		return nil
	}
	return s.server.ListenAndServe()
}

func (s *MetricServer) Stop(ctx context.Context) error {
	if !s.config.Enabled {
		return nil
	}
	return s.server.Shutdown(ctx)
}
