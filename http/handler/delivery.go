package handler

import (
	"net/http"

	"github.com/allisson/postmand"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type deliveryList struct {
	Deliveries []*postmand.Delivery `json:"deliveries"`
	Limit      int                  `json:"limit"`
	Offset     int                  `json:"offset"`
} //@name DeliveryList

// Delivery implements rest interface for delivery.
type Delivery struct {
	deliveryService postmand.DeliveryService
	logger          *zap.Logger
}

// List deliveries.
// List godoc
// @Summary List deliveries
// @Tags deliveries
// @Accept json
// @Produce json
// @Param limit query int false "The limit indicates the maximum number of items to return"
// @Param offset query int false "The offset indicates the starting position of the query in relation to the complete set of unpaginated items"
// @Param webhook_id query string false "Filter by webhook_id field"
// @Param status query string false "Filter by status field"
// @Param created_at.gt query string false "Return results where the created_at field is greater than this value"
// @Param created_at.gte query string false "Return results where the created_at field is greater than or equal to this value"
// @Param created_at.lt query string false "Return results where the created_at field is less than this value"
// @Param created_at.lte query string false "Return results where the created_at field is less than or equal to this value"
// @Success 200 {object} deliveryList
// @Failure 500 {object} errorResponse
// @Router /deliveries [get]
func (d Delivery) List(w http.ResponseWriter, r *http.Request) {
	listOptions := makeListOptions(r, []string{"webhook_id", "status", "created_at.gt", "created_at.gte", "created_at.lt", "created_at.lte"})
	listOptions.OrderBy = "created_at"
	listOptions.Order = "desc"

	// Call service
	deliveries, err := d.deliveryService.List(r.Context(), listOptions)
	if err != nil {
		d.logger.Error(
			"service-error",
			zap.String("name", "DeliveryService"),
			zap.String("method", "List"),
			zap.Error(err),
		)
		er := errorResponses["internal_server_error"]
		makeErrorResponse(w, &er, d.logger)
		return
	}

	// Return response
	dl := deliveryList{
		Deliveries: deliveries,
		Limit:      listOptions.Limit,
		Offset:     listOptions.Offset,
	}
	makeJSONResponse(w, http.StatusOK, dl, d.logger)
}

// Get delivery.
// Get godoc
// @Summary Show a delivery
// @Tags deliveries
// @Accept json
// @Produce json
// @Param delivery_id path string true "Delivery ID"
// @Success 200 {object} postmand.Delivery
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /deliveries/{delivery_id} [get]
func (d Delivery) Get(w http.ResponseWriter, r *http.Request) {
	deliveryID, err := uuid.Parse(chi.URLParam(r, "delivery_id"))
	if err != nil {
		er := errorResponses["invalid_id"]
		makeErrorResponse(w, &er, d.logger)
		return
	}

	// Call service
	getOptions := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": deliveryID}}
	delivery, err := d.deliveryService.Get(r.Context(), getOptions)
	if err != nil {
		if err == postmand.ErrDeliveryNotFound {
			er := errorResponses["delivery_not_found"]
			makeErrorResponse(w, &er, d.logger)
			return
		}
		d.logger.Error(
			"service-error",
			zap.String("name", "DeliveryService"),
			zap.String("method", "Get"),
			zap.Error(err),
		)
		er := errorResponses["internal_server_error"]
		makeErrorResponse(w, &er, d.logger)
		return
	}

	// Return response
	makeJSONResponse(w, http.StatusOK, delivery, d.logger)
}

// Create delivery.
// Create godoc
// @Summary Add an delivery
// @Tags deliveries
// @Accept json
// @Produce json
// @Param delivery body postmand.Delivery true "Add delivery"
// @Success 201 {object} postmand.Delivery
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /deliveries [post]
func (d Delivery) Create(w http.ResponseWriter, r *http.Request) {
	// Parse request
	delivery := postmand.Delivery{}
	if er := readBodyJSON(r, &delivery, d.logger); er != nil {
		makeErrorResponse(w, er, d.logger)
		return
	}

	// Call service
	if err := d.deliveryService.Create(r.Context(), &delivery); err != nil {
		d.logger.Error(
			"service-error",
			zap.String("name", "DeliveryService"),
			zap.String("method", "Create"),
			zap.Error(err),
		)
		er := errorResponses["internal_server_error"]
		makeErrorResponse(w, &er, d.logger)
	}

	// Return response
	makeJSONResponse(w, http.StatusCreated, delivery, d.logger)
}

// Delete delivery.
// Delete godoc
// @Summary Delete an delivery
// @Tags deliveries
// @Accept json
// @Produce json
// @Param delivery_id path string true "Delivery ID"
// @Success 204 "No Content"
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /deliveries/{delivery_id} [delete]
func (d Delivery) Delete(w http.ResponseWriter, r *http.Request) {
	deliveryID, err := uuid.Parse(chi.URLParam(r, "delivery_id"))
	if err != nil {
		er := errorResponses["invalid_id"]
		makeErrorResponse(w, &er, d.logger)
		return
	}

	// Call service
	if err := d.deliveryService.Delete(r.Context(), deliveryID); err != nil {
		d.logger.Error(
			"service-error",
			zap.String("name", "DeliveryService"),
			zap.String("method", "Delete"),
			zap.Error(err),
		)
		er := errorResponses["internal_server_error"]
		makeErrorResponse(w, &er, d.logger)
	}

	// Return response
	makeResponse(w, []byte(""), http.StatusNoContent, "application/json", d.logger)
}

// NewDelivery creates a new Delivery.
func NewDelivery(deliveryService postmand.DeliveryService, logger *zap.Logger) *Delivery {
	return &Delivery{
		deliveryService: deliveryService,
		logger:          logger,
	}
}
