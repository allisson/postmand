package http

import (
	"context"
	"fmt"
	nethttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	mw "github.com/allisson/postmand/http/middleware"
)

// NewRouter returns *chi.Mux with base middlewares.
func NewRouter(logger *zap.Logger) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(mw.NewZapMiddleware("router", logger))
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(mw.NewSecureHeadersMiddleware())
	r.Use(middleware.Recoverer)
	return r
}

// Server implements a http server.
type Server struct {
	mux               *chi.Mux
	httpPort          int
	readHeaderTimeout time.Duration
	logger            *zap.Logger
}

// Run starts a http server.
func (s Server) Run() {
	httpServer := &nethttp.Server{
		Addr:              fmt.Sprintf(":%d", s.httpPort),
		Handler:           s.mux,
		ReadHeaderTimeout: s.readHeaderTimeout,
	}
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)

		// interrupt signal sent from terminal
		signal.Notify(sigint, os.Interrupt)
		// sigterm signal sent from kubernetes
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint

		// We received an interrupt signal, shut down.
		s.logger.Info("http-server-shutdown-started")
		if err := httpServer.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			s.logger.Error("http-server-shutdown", zap.Error(err))
		}
		close(idleConnsClosed)
		s.logger.Info("http-server-shutdown-finished")
	}()

	s.logger.Info("http-server-listen-and-server")
	if err := httpServer.ListenAndServe(); err != nil {
		if err != nethttp.ErrServerClosed {
			s.logger.Error("http-server-listen-and-server-error", zap.Error(err))
			return
		}
	}

	<-idleConnsClosed
}

// NewServer creates a new Server.
func NewServer(mux *chi.Mux, httpPort int, logger *zap.Logger) *Server {
	return &Server{
		mux:               mux,
		httpPort:          httpPort,
		readHeaderTimeout: time.Second * 60,
		logger:            logger,
	}
}
