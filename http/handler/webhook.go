package handler

import (
	"net/http"

	"github.com/allisson/postmand"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type webhookList struct {
	Webhooks []*postmand.Webhook `json:"webhooks"`
	Limit    int                 `json:"limit"`
	Offset   int                 `json:"offset"`
}

// Webhook implements rest interface for webhook.
type Webhook struct {
	webhookService postmand.WebhookService
	logger         *zap.Logger
}

// List webhooks.
func (wh Webhook) List(w http.ResponseWriter, r *http.Request) {
	listOptions, err := makeListOptions(r, []string{})
	if err != nil {
		er := errorResponses["internal_server_error"]
		makeErrorResponse(w, &er, wh.logger)
		return
	}
	listOptions.OrderBy = "name"
	listOptions.Order = "asc"

	// Call service
	webhooks, err := wh.webhookService.List(r.Context(), listOptions)
	if err != nil {
		wh.logger.Error(
			"service-error",
			zap.String("name", "WebhookService"),
			zap.String("method", "List"),
			zap.Error(err),
		)
		er := errorResponses["internal_server_error"]
		makeErrorResponse(w, &er, wh.logger)
		return
	}

	// Return response
	wl := webhookList{
		Webhooks: webhooks,
		Limit:    listOptions.Limit,
		Offset:   listOptions.Offset,
	}
	makeJSONResponse(w, http.StatusOK, wl, wh.logger)
}

// Get webhook.
func (wh Webhook) Get(w http.ResponseWriter, r *http.Request) {
	webhookID, err := uuid.Parse(chi.URLParam(r, "webhook_id"))
	if err != nil {
		er := errorResponses["invalid_id"]
		makeErrorResponse(w, &er, wh.logger)
		return
	}

	// Call service
	getOptions := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": webhookID}}
	webhook, err := wh.webhookService.Get(r.Context(), getOptions)
	if err != nil {
		if err == postmand.ErrWebhookNotFound {
			er := errorResponses["webhook_not_found"]
			makeErrorResponse(w, &er, wh.logger)
			return
		}
		wh.logger.Error(
			"service-error",
			zap.String("name", "WebhookService"),
			zap.String("method", "Get"),
			zap.Error(err),
		)
		er := errorResponses["internal_server_error"]
		makeErrorResponse(w, &er, wh.logger)
		return
	}

	// Return response
	makeJSONResponse(w, http.StatusOK, webhook, wh.logger)
}

// Create webhook.
func (wh Webhook) Create(w http.ResponseWriter, r *http.Request) {
	// Parse request
	webhook := postmand.Webhook{}
	if er := readBodyJSON(r, &webhook, wh.logger); er != nil {
		makeErrorResponse(w, er, wh.logger)
		return
	}

	// Call service
	if err := wh.webhookService.Create(r.Context(), &webhook); err != nil {
		wh.logger.Error(
			"service-error",
			zap.String("name", "WebhookService"),
			zap.String("method", "Create"),
			zap.Error(err),
		)
		er := errorResponses["internal_server_error"]
		makeErrorResponse(w, &er, wh.logger)
	}

	// Return response
	makeJSONResponse(w, http.StatusCreated, webhook, wh.logger)
}

// Update webhook.
func (wh Webhook) Update(w http.ResponseWriter, r *http.Request) {
	webhookID, err := uuid.Parse(chi.URLParam(r, "webhook_id"))
	if err != nil {
		er := errorResponses["invalid_id"]
		makeErrorResponse(w, &er, wh.logger)
		return
	}

	// Parse request
	webhook := postmand.Webhook{}
	if er := readBodyJSON(r, &webhook, wh.logger); er != nil {
		makeErrorResponse(w, er, wh.logger)
		return
	}
	webhook.ID = webhookID

	// Call service
	if err := wh.webhookService.Update(r.Context(), &webhook); err != nil {
		wh.logger.Error(
			"service-error",
			zap.String("name", "WebhookService"),
			zap.String("method", "Update"),
			zap.Error(err),
		)
		er := errorResponses["internal_server_error"]
		makeErrorResponse(w, &er, wh.logger)
	}

	// Return response
	makeJSONResponse(w, http.StatusOK, webhook, wh.logger)
}

// Delete webhook.
func (wh Webhook) Delete(w http.ResponseWriter, r *http.Request) {
	webhookID, err := uuid.Parse(chi.URLParam(r, "webhook_id"))
	if err != nil {
		er := errorResponses["invalid_id"]
		makeErrorResponse(w, &er, wh.logger)
		return
	}

	// Call service
	if err := wh.webhookService.Delete(r.Context(), webhookID); err != nil {
		wh.logger.Error(
			"service-error",
			zap.String("name", "WebhookService"),
			zap.String("method", "Delete"),
			zap.Error(err),
		)
		er := errorResponses["internal_server_error"]
		makeErrorResponse(w, &er, wh.logger)
	}

	// Return response
	makeResponse(w, []byte(""), http.StatusNoContent, "application/json", wh.logger)
}

// NewWebhook creates a new Webhook.
func NewWebhook(webhookService postmand.WebhookService, logger *zap.Logger) *Webhook {
	return &Webhook{
		webhookService: webhookService,
		logger:         logger,
	}
}
