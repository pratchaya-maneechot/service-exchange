package observability

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricConfig struct {
	Path    string
	Addr    string
	Enabled bool
}
type MetricServer struct {
	server *http.Server
	config MetricConfig
}

func NewMetricServer(cfg MetricConfig) *MetricServer {
	mux := http.NewServeMux()
	mux.Handle(cfg.Path, promhttp.Handler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return &MetricServer{
		server: &http.Server{
			Addr:         cfg.Addr,
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
