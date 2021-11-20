package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/allisson/postmand"
)

type deliveryAttemptList struct {
	DeliveryAttempts []*postmand.DeliveryAttempt `json:"delivery_attempts"`
	Limit            int                         `json:"limit"`
	Offset           int                         `json:"offset"`
} //@name DeliveryAttemptList

// DeliveryAttempt implements rest interface for delivery attempt.
type DeliveryAttempt struct {
	deliveryAttemptService postmand.DeliveryAttemptService
	logger                 *zap.Logger
}

// List delivery attempts.
// List godoc
// @Summary List delivery attempts
// @Tags delivery-attempts
// @Accept json
// @Produce json
// @Param limit query int false "The limit indicates the maximum number of items to return"
// @Param offset query int false "The offset indicates the starting position of the query in relation to the complete set of unpaginated items"
// @Param webhook_id query string false "Filter by webhook_id"
// @Param delivery_id query string false "Filter by delivery_id"
// @Param success query boolean false "Filter by success"
// @Param created_at.gt query string false "Return results where the created_at field is greater than this value"
// @Param created_at.gte query string false "Return results where the created_at field is greater than or equal to this value"
// @Param created_at.lt query string false "Return results where the created_at field is less than this value"
// @Param created_at.lte query string false "Return results where the created_at field is less than or equal to this value"
// @Success 200 {object} deliveryAttemptList
// @Failure 500 {object} errorResponse
// @Router /delivery-attempts [get]
func (d DeliveryAttempt) List(w http.ResponseWriter, r *http.Request) {
	listOptions := makeListOptions(r, []string{"webhook_id", "delivery_id", "success", "created_at.gt", "created_at.gte", "created_at.lt", "created_at.lte"})
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
// Get godoc
// @Summary Show a delivery attempt
// @Tags delivery-attempts
// @Accept json
// @Produce json
// @Param delivery_attempt_id path string true "Delivery Attempt ID"
// @Success 200 {object} postmand.DeliveryAttempt
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /delivery-attempts/{delivery_attempt_id} [get]
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
