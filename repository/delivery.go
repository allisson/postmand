package repository

import (
	"bytes"
	"database/sql"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/allisson/postmand"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/jpillora/backoff"
)

type dispatchResponse struct {
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
func (d Delivery) Get(getOptions postmand.RepositoryGetOptions) (*postmand.Delivery, error) {
	delivery := postmand.Delivery{}
	sql, args := getQuery("deliveries", getOptions)
	err := d.db.Get(&delivery, sql, args...)
	return &delivery, err
}

// List returns a slice of postmand.Delivery by options filter.
func (d Delivery) List(listOptions postmand.RepositoryListOptions) ([]*postmand.Delivery, error) {
	deliveries := []*postmand.Delivery{}
	sql, args := listQuery("deliveries", listOptions)
	err := d.db.Select(&deliveries, sql, args...)
	return deliveries, err
}

// Create postmand.Delivery on database.
func (d Delivery) Create(delivery *postmand.Delivery) error {
	sql, args := insertQuery("deliveries", delivery)
	_, err := d.db.Exec(sql, args...)
	return err
}

// Update postmand.Delivery on database.
func (d Delivery) Update(delivery *postmand.Delivery) error {
	sql, args := updateQuery("deliveries", delivery.ID, delivery)
	_, err := d.db.Exec(sql, args...)
	return err
}

// Delete postmand.Delivery on database.
func (d Delivery) Delete(id postmand.ID) error {
	sqlStatement := `
		DELETE FROM deliveries WHERE id = $1
	`
	_, err := d.db.Exec(sqlStatement, id)
	return err
}

// Dispatch fetchs a delivery and send to url destination.
func (d Delivery) Dispatch() error {
	sqlStatement := `
		SELECT
			deliveries.*
		FROM
			deliveries
		INNER JOIN webhooks
			ON deliveries.webhook_id = webhooks.id
		WHERE
			webhooks.active = true AND deliveries.status = 'pending' AND deliveries.scheduled_at <= now()
		ORDER BY
			deliveries.created_at ASC
		FOR UPDATE SKIP LOCKED
		LIMIT
			1
	`

	// Starts a new transaction
	tx, err := d.db.Beginx()
	if err != nil {
		return err
	}

	// Get delivery
	delivery := postmand.Delivery{}
	err = tx.Get(&delivery, sqlStatement)
	if err != nil {
		// Skip if no result
		if err == sql.ErrNoRows {
			rollback("delivery not found", tx)
			return nil
		}
		rollback("get delivery", tx)
		return err
	}

	// Get webhook
	webhook := postmand.Webhook{}
	getOptions := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": delivery.WebhookID}}
	sql, args := getQuery("webhooks", getOptions)
	err = tx.Get(&webhook, sql, args...)
	if err != nil {
		rollback("get webhook", tx)
		return err
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
	sql, args = updateQuery("deliveries", delivery.ID, delivery)
	_, err = tx.Exec(sql, args...)
	if err != nil {
		rollback("update delivery", tx)
		return err
	}

	// Create delivery attempt
	deliveryAttempt := postmand.DeliveryAttempt{
		ID:                 uuid.New(),
		WebhookID:          webhook.ID,
		DeliveryID:         delivery.ID,
		RawResponse:        dr.RawResponse,
		ResponseStatusCode: dr.ResponseStatusCode,
		ExecutionDuration:  dr.ExecutionDuration,
		Success:            dr.Success,
		Error:              dr.Error,
		CreatedAt:          time.Now().UTC(),
	}
	sql, args = insertQuery("delivery_attempts", deliveryAttempt)
	_, err = tx.Exec(sql, args...)
	if err != nil {
		rollback("create delivery attempt", tx)
		return err
	}

	commit("dispatch", tx)
	return nil
}

// NewDelivery returns postmand.Delivery with db connection.
func NewDelivery(db *sqlx.DB) *Delivery {
	return &Delivery{db: db}
}
