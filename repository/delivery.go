package repository

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jpillora/backoff"

	"github.com/allisson/postmand"
)

type dispatchResponse struct {
	RawRequest         string
	RawResponse        string
	ResponseStatusCode int
	ExecutionDuration  int
	Success            bool
	Error              string
}

func dispatchToURL(webhook *postmand.Webhook, delivery *postmand.Delivery) dispatchResponse {
	dr := dispatchResponse{}

	// Prepare request
	httpClient := &http.Client{Timeout: time.Duration(webhook.DeliveryAttemptTimeout) * time.Second}
	request, err := http.NewRequest("POST", webhook.URL, bytes.NewBufferString(delivery.Payload))
	if err != nil {
		dr.Success = false
		dr.Error = err.Error()
		return dr
	}
	request.Header.Set("Content-Type", webhook.ContentType)
	if webhook.SecretToken != "" {
		hash := hmac.New(sha256.New, []byte(webhook.SecretToken))
		_, err := hash.Write([]byte(delivery.Payload))
		if err != nil {
			dr.Success = false
			dr.Error = err.Error()
			return dr
		}
		request.Header.Set("X-Hub-Signature", hex.EncodeToString(hash.Sum(nil)))
	}

	// Create request dump
	requestDump, err := httputil.DumpRequest(request, true)
	if err != nil {
		dr.Success = false
		dr.Error = err.Error()
		return dr
	}

	// Make request
	start := time.Now()
	response, err := httpClient.Do(request)
	if err != nil {
		dr.Success = false
		dr.Error = err.Error()
		return dr
	}
	latency := time.Since(start)

	// Create response dump
	responseDump, err := httputil.DumpResponse(response, true)
	if err != nil {
		dr.Success = false
		dr.Error = err.Error()
		return dr
	}

	// Verify response status code
	success := false
	for _, statusCode := range webhook.ValidStatusCodes {
		if response.StatusCode == int(statusCode) {
			success = true
		}
	}

	// Update dispatch response
	dr.RawRequest = string(requestDump)
	dr.RawResponse = string(responseDump)
	dr.ResponseStatusCode = response.StatusCode
	dr.ExecutionDuration = int(latency.Milliseconds())
	dr.Success = success

	return dr
}

// Delivery implements postmand.DeliveryRepository interface.
type Delivery struct {
	db *sqlx.DB
}

// Get returns postmand.Delivery by options filter.
func (d Delivery) Get(ctx context.Context, getOptions postmand.RepositoryGetOptions) (*postmand.Delivery, error) {
	delivery := postmand.Delivery{}
	query, args := getQuery("deliveries", getOptions)
	err := d.db.GetContext(ctx, &delivery, query, args...)
	if err == sql.ErrNoRows {
		return &delivery, postmand.ErrDeliveryNotFound
	}
	return &delivery, err
}

// List returns a slice of postmand.Delivery by options filter.
func (d Delivery) List(ctx context.Context, listOptions postmand.RepositoryListOptions) ([]*postmand.Delivery, error) {
	deliveries := []*postmand.Delivery{}
	query, args := listQuery("deliveries", listOptions)
	err := d.db.SelectContext(ctx, &deliveries, query, args...)
	return deliveries, err
}

// Create postmand.Delivery on database.
func (d Delivery) Create(ctx context.Context, delivery *postmand.Delivery) error {
	query, args := insertQuery("deliveries", delivery)
	_, err := d.db.ExecContext(ctx, query, args...)
	return err
}

// Update postmand.Delivery on database.
func (d Delivery) Update(ctx context.Context, delivery *postmand.Delivery) error {
	query, args := updateQuery("deliveries", delivery.ID, delivery)
	_, err := d.db.ExecContext(ctx, query, args...)
	return err
}

// Delete postmand.Delivery on database.
func (d Delivery) Delete(ctx context.Context, id postmand.ID) error {
	query := `
		DELETE FROM deliveries WHERE id = $1
	`
	_, err := d.db.ExecContext(ctx, query, id)
	return err
}

// Dispatch fetchs a delivery and send to url destination.
func (d Delivery) Dispatch(ctx context.Context) (*postmand.DeliveryAttempt, error) {
	query := `
		SELECT
			deliveries.*
		FROM
			deliveries
		INNER JOIN webhooks
			ON deliveries.webhook_id = webhooks.id
		WHERE
			webhooks.active = true AND deliveries.status = $1 AND deliveries.scheduled_at <= $2
		ORDER BY
			deliveries.created_at ASC
		FOR UPDATE SKIP LOCKED
		LIMIT
			1
	`

	// Starts a new transaction
	tx, err := d.db.Beginx()
	if err != nil {
		return nil, err
	}

	// Get delivery
	delivery := postmand.Delivery{}
	err = tx.GetContext(ctx, &delivery, query, postmand.DeliveryStatusPending, time.Now().UTC())
	if err != nil {
		// Skip if no result
		if err == sql.ErrNoRows {
			rollback("delivery not found", tx)
			return nil, nil
		}
		rollback("get delivery", tx)
		return nil, err
	}

	// Get webhook
	webhook := postmand.Webhook{}
	getOptions := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": delivery.WebhookID}}
	sql, args := getQuery("webhooks", getOptions)
	err = tx.GetContext(ctx, &webhook, sql, args...)
	if err != nil {
		rollback("get webhook", tx)
		return nil, err
	}

	// Dispatch webhook
	dr := dispatchToURL(&webhook, &delivery)

	// Update delivery
	newDeliveryAttempts := delivery.DeliveryAttempts + 1
	newStatus := postmand.DeliveryStatusPending
	newScheduledAt := delivery.ScheduledAt
	if dr.Success {
		newStatus = postmand.DeliveryStatusSucceeded
	} else {
		if newDeliveryAttempts >= webhook.MaxDeliveryAttempts {
			newStatus = postmand.DeliveryStatusFailed
		} else {
			b := &backoff.Backoff{
				Min:    time.Duration(webhook.RetryMinBackoff) * time.Second,
				Max:    time.Duration(webhook.RetryMaxBackoff) * time.Second,
				Factor: 2,
				Jitter: false,
			}
			newScheduledAt = time.Now().UTC().Add(b.ForAttempt(float64(delivery.DeliveryAttempts)))
		}
	}
	delivery.DeliveryAttempts = newDeliveryAttempts
	delivery.Status = newStatus
	delivery.ScheduledAt = newScheduledAt
	delivery.UpdatedAt = time.Now().UTC()
	query, args = updateQuery("deliveries", delivery.ID, delivery)
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		rollback("update delivery", tx)
		return nil, err
	}

	// Create delivery attempt
	deliveryAttempt := postmand.DeliveryAttempt{
		ID:                 uuid.New(),
		WebhookID:          webhook.ID,
		DeliveryID:         delivery.ID,
		RawRequest:         dr.RawRequest,
		RawResponse:        dr.RawResponse,
		ResponseStatusCode: dr.ResponseStatusCode,
		ExecutionDuration:  dr.ExecutionDuration,
		Success:            dr.Success,
		Error:              dr.Error,
		CreatedAt:          time.Now().UTC(),
	}
	query, args = insertQuery("delivery_attempts", deliveryAttempt)
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		rollback("create delivery attempt", tx)
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		rollback("unable to commit", tx)
		return nil, err
	}

	return &deliveryAttempt, nil
}

// NewDelivery will create an implementation of postmand.DeliveryRepository.
func NewDelivery(db *sqlx.DB) *Delivery {
	return &Delivery{db: db}
}
