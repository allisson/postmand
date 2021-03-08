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
}

// Delivery implements rest interface for delivery.
type Delivery struct {
	deliveryService postmand.DeliveryService
	logger          *zap.Logger
}

// List deliveries.
func (d Delivery) List(w http.ResponseWriter, r *http.Request) {
	listOptions, err := makeListOptions(r, []string{"webhook_id", "status"})
	if err != nil {
		er := errorResponses["internal_server_error"]
		makeErrorResponse(w, &er, d.logger)
		return
	}
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
