package internal

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ShutdownHandler struct {
	logger          *slog.Logger
	ShutdownTimeout time.Duration
}

func NewShutdownHandler(logger *slog.Logger, timeout time.Duration) *ShutdownHandler {
	return &ShutdownHandler{
		logger:          logger,
		ShutdownTimeout: timeout,
	}
}

func (h *ShutdownHandler) HandleSignals(cancel context.CancelFunc) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	sig := <-sigCh
	h.logger.Info("received shutdown signal", "signal", sig.String())

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), h.ShutdownTimeout)
	defer shutdownCancel()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		remaining := h.ShutdownTimeout
		for {
			select {
			case <-shutdownCtx.Done():
				return
			case <-ticker.C:
				remaining -= 5 * time.Second
				if remaining > 0 {
					h.logger.Info("shutdown in progress", "remaining", remaining.String())
				}
			}
		}
	}()

	cancel()

	<-shutdownCtx.Done()
	if shutdownCtx.Err() == context.DeadlineExceeded {
		h.logger.Warn("shutdown timeout exceeded, forcing exit")
		os.Exit(1)
	}
}
