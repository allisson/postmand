package handler

import (
	"net/http"

	"github.com/allisson/postmand"
	"go.uber.org/zap"
)

type pingResponse struct {
	Success bool `json:"success"`
}

// Ping implements interface for health check.
type Ping struct {
	pingService postmand.PingService
	logger      *zap.Logger
}

// Healthz returns health check response.
func (p Ping) Healthz(w http.ResponseWriter, r *http.Request) {
	pr := pingResponse{}

	if err := p.pingService.Run(r.Context()); err != nil {
		p.logger.Error(
			"service-error",
			zap.String("name", "PingService"),
			zap.String("method", "Run"),
			zap.Error(err),
		)
		makeJSONResponse(w, http.StatusInternalServerError, &pr, p.logger)
		return
	}

	pr.Success = true
	makeJSONResponse(w, http.StatusOK, &pr, p.logger)
}

// NewPing creates a new Ping.
func NewPing(pingService postmand.PingService, logger *zap.Logger) *Ping {
	return &Ping{
		pingService: pingService,
		logger:      logger,
	}
}
