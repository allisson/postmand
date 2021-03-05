package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Adapted from https://github.com/766b/chi-logger/blob/master/middleware.go
type chilogger struct {
	logger *zap.Logger
	name   string
}

func (c chilogger) middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		var requestID string
		if reqID := r.Context().Value(middleware.RequestIDKey); reqID != nil {
			requestID = reqID.(string)
		}
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(ww, r)
		latency := time.Since(start)
		fields := []zapcore.Field{
			zap.Int("status", ww.Status()),
			zap.Duration("took", latency),
			zap.Int64(fmt.Sprintf("measure#%s.latency", c.name), latency.Nanoseconds()),
			zap.String("remote", r.RemoteAddr),
			zap.String("request", r.RequestURI),
			zap.String("method", r.Method),
		}
		if requestID != "" {
			fields = append(fields, zap.String("request-id", requestID))
		}
		c.logger.Info("request-completed", fields...)
	}
	return http.HandlerFunc(fn)
}

// NewZapMiddleware returns a new Zap Middleware handler.
func NewZapMiddleware(name string, logger *zap.Logger) func(next http.Handler) http.Handler {
	return chilogger{
		logger: logger,
		name:   name,
	}.middleware
}
