package metrics

import (
	"context"
	"net/http"
	"time"

	"github.com/pratchaya-maneechot/service-exchange/apps/users/internal/config"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricServer struct {
	server *http.Server
	config config.MetricsConfig
}

func NewServer(cfg *config.Config) *MetricServer {
	mux := http.NewServeMux()
	mux.Handle(cfg.Metrics.Path, promhttp.Handler())
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	return &MetricServer{
		server: &http.Server{
			Addr:         cfg.Metrics.Address,
			Handler:      mux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		config: cfg.Metrics,
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
