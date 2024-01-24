package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hosseintrz/gaterun/config/models"
)

const (
	GateRunUserAgent = "GateRun"
)

var (
	UserAgentHeader = []string{GateRunUserAgent}
)

func RunServer(ctx context.Context, cfg models.ServiceConfig, handler http.Handler) error {
	errChan := make(chan error)
	s := NewServer(cfg, handler)

	go func() {
		errChan <- s.ListenAndServe()
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return s.Shutdown(context.Background())
	}
}

func NewServer(cfg models.ServiceConfig, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:           handler,
		ReadTimeout:       cfg.Timeout,
		WriteTimeout:      cfg.Timeout,
		ReadHeaderTimeout: cfg.Timeout,
		IdleTimeout:       cfg.Timeout,
	}
}
