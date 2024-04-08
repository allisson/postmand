package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/crypitor/postmand"
)

type webhookList struct {
	Webhooks []*postmand.Webhook `json:"webhooks"`
	Limit    int                 `json:"limit"`
	Offset   int                 `json:"offset"`
} //@name WebhookList

// Webhook implements rest interface for webhook.
type Webhook struct {
	webhookService postmand.WebhookService
	logger         *zap.Logger
}

// List webhooks.
// List godoc
// @Summary List webhooks
// @Tags webhooks
// @Accept json
// @Produce json
// @Param limit query int false "The limit indicates the maximum number of items to return"
// @Param offset query int false "The offset indicates the starting position of the query in relation to the complete set of unpaginated items"
// @Param active query boolean false "Filter by active field"
// @Param created_at.gt query string false "Return results where the created_at field is greater than this value"
// @Param created_at.gte query string false "Return results where the created_at field is greater than or equal to this value"
// @Param created_at.lt query string false "Return results where the created_at field is less than this value"
// @Param created_at.lte query string false "Return results where the created_at field is less than or equal to this value"
// @Success 200 {object} webhookList
// @Failure 500 {object} errorResponse
// @Router /webhooks [get]
func (wh Webhook) List(w http.ResponseWriter, r *http.Request) {
	listOptions := makeListOptions(r, []string{"active", "created_at.gt", "created_at.gte", "created_at.lt", "created_at.lte"})
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
// Get godoc
// @Summary Show a webhook
// @Tags webhooks
// @Accept json
// @Produce json
// @Param webhook_id path string true "Webhook ID"
// @Success 200 {object} postmand.Webhook
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /webhooks/{webhook_id} [get]
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
// Create godoc
// @Summary Add an webhook
// @Tags webhooks
// @Accept json
// @Produce json
// @Param webhook body postmand.Webhook true "Add webhook"
// @Success 201 {object} postmand.Webhook
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /webhooks [post]
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
// Update godoc
// @Summary Update an webhook
// @Tags webhooks
// @Accept json
// @Produce json
// @Param webhook_id path string true "Webhook ID"
// @Param webhook body postmand.Webhook true "Update webhook"
// @Success 200 {object} postmand.Webhook
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /webhooks/{webhook_id} [put]
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
// Delete godoc
// @Summary Delete an webhook
// @Tags webhooks
// @Accept json
// @Produce json
// @Param webhook_id path string true "Webhook ID"
// @Success 204 "No Content"
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /webhooks/{webhook_id} [delete]
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
