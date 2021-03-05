package handler

import (
	"net/http"

	"github.com/allisson/postmand"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type deliveryAttemptList struct {
	DeliveryAttempts []*postmand.DeliveryAttempt `json:"delivery_attempts"`
	Limit            int                         `json:"limit"`
	Offset           int                         `json:"offset"`
}

// DeliveryAttempt implements rest interface for delivery attempt.
type DeliveryAttempt struct {
	deliveryAttemptService postmand.DeliveryAttemptService
	logger                 *zap.Logger
}

// List delivery attempts.
func (d DeliveryAttempt) List(w http.ResponseWriter, r *http.Request) {
	listOptions, err := makeListOptions(r, []string{"webhook_id", "delivery_id", "success"})
	if err != nil {
		er := errorResponses["internal_server_error"]
		makeErrorResponse(w, &er, d.logger)
		return
	}
	listOptions.OrderBy = "created_at"
	listOptions.Order = "desc"

	// Call service
	deliveryAttempts, err := d.deliveryAttemptService.List(r.Context(), listOptions)
	if err != nil {
		d.logger.Error(
			"service-error",
			zap.String("name", "DeliveryAttemptService"),
			zap.String("method", "List"),
			zap.Error(err),
		)
		er := errorResponses["internal_server_error"]
		makeErrorResponse(w, &er, d.logger)
		return
	}

	// Return response
	dl := deliveryAttemptList{
		DeliveryAttempts: deliveryAttempts,
		Limit:            listOptions.Limit,
		Offset:           listOptions.Offset,
	}
	makeJSONResponse(w, http.StatusOK, dl, d.logger)
}

// Get delivery attempt.
func (d DeliveryAttempt) Get(w http.ResponseWriter, r *http.Request) {
	deliveryAttemptID, err := uuid.Parse(chi.URLParam(r, "delivery_attempt_id"))
	if err != nil {
		er := errorResponses["invalid_id"]
		makeErrorResponse(w, &er, d.logger)
		return
	}

	// Call service
	getOptions := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": deliveryAttemptID}}
	deliveryAttempt, err := d.deliveryAttemptService.Get(r.Context(), getOptions)
	if err != nil {
		if err == postmand.ErrDeliveryAttemptNotFound {
			er := errorResponses["delivery_attempt_not_found"]
			makeErrorResponse(w, &er, d.logger)
			return
		}
		d.logger.Error(
			"service-error",
			zap.String("name", "DeliveryAttemptService"),
			zap.String("method", "Get"),
			zap.Error(err),
		)
		er := errorResponses["internal_server_error"]
		makeErrorResponse(w, &er, d.logger)
		return
	}

	// Return response
	makeJSONResponse(w, http.StatusOK, deliveryAttempt, d.logger)
}

// NewDeliveryAttempt creates a new DeliveryAttempt.
func NewDeliveryAttempt(deliveryAttemptService postmand.DeliveryAttemptService, logger *zap.Logger) *DeliveryAttempt {
	return &DeliveryAttempt{
		deliveryAttemptService: deliveryAttemptService,
		logger:                 logger,
	}
}
