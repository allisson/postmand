package postmand

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

const (
	// DeliveryStatusPending represents the delivery pending status
	DeliveryStatusPending = "pending"
	// DeliveryStatusSucceeded represents the delivery succeeded status
	DeliveryStatusSucceeded = "succeeded"
	// DeliveryStatusFailed represents the delivery failed status
	DeliveryStatusFailed = "failed"
)

// ID represents the primary key for all entities.
type ID = uuid.UUID

// Webhook represents a webhook in the system.
type Webhook struct {
	ID                     ID            `json:"id" db:"id"`
	Name                   string        `json:"name" db:"name"`
	URL                    string        `json:"url" db:"url"`
	ContentType            string        `json:"content_type" db:"content_type"`
	ValidStatusCodes       pq.Int32Array `json:"valid_status_codes" db:"valid_status_codes"`
	SecretToken            string        `json:"secret_token" db:"secret_token"`
	Active                 bool          `json:"active" db:"active"`
	MaxDeliveryAttempts    int           `json:"max_delivery_attempts" db:"max_delivery_attempts"`
	DeliveryAttemptTimeout int           `json:"delivery_attempt_timeout" db:"delivery_attempt_timeout"`
	RetryMinBackoff        int           `json:"retry_min_backoff" db:"retry_min_backoff"`
	RetryMaxBackoff        int           `json:"retry_max_backoff" db:"retry_max_backoff"`
	CreatedAt              time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time     `json:"updated_at" db:"updated_at"`
}

// Validate implements ozzo validation Validatable interface
func (w Webhook) Validate() error {
	return validation.ValidateStruct(&w,
		validation.Field(&w.Name, validation.Required, validation.Length(3, 255)),
		validation.Field(&w.URL, validation.Required, is.URL),
		validation.Field(&w.ContentType, validation.Required, validation.In("application/x-www-form-urlencoded", "application/json")),
		validation.Field(&w.ValidStatusCodes, validation.Required),
		validation.Field(&w.MaxDeliveryAttempts, validation.Required, validation.Min(1)),
		validation.Field(&w.DeliveryAttemptTimeout, validation.Required, validation.Min(1)),
		validation.Field(&w.RetryMinBackoff, validation.Required, validation.Min(1)),
		validation.Field(&w.RetryMaxBackoff, validation.Required, validation.Min(1)),
	)
}

// Delivery represents a payload that must be delivery using webhook context.
type Delivery struct {
	ID               ID        `json:"id" db:"id"`
	WebhookID        ID        `json:"webhook_id" db:"webhook_id"`
	Payload          string    `json:"payload" db:"payload"`
	ScheduledAt      time.Time `json:"scheduled_at" db:"scheduled_at"`
	DeliveryAttempts int       `json:"delivery_attempts" db:"delivery_attempts"`
	Status           string    `json:"status" db:"status"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// Validate implements ozzo validation Validatable interface
func (d Delivery) Validate() error {
	return validation.ValidateStruct(&d,
		validation.Field(&d.WebhookID, validation.Required, is.UUIDv4),
	)
}

// DeliveryAttempt represents a delivery attempt.
type DeliveryAttempt struct {
	ID                 ID        `json:"id" db:"id"`
	WebhookID          ID        `json:"webhook_id" db:"webhook_id"`
	DeliveryID         ID        `json:"delivery_id" db:"delivery_id"`
	RawRequest         string    `json:"raw_request" db:"raw_request"`
	RawResponse        string    `json:"raw_response" db:"raw_response"`
	ResponseStatusCode int       `json:"response_status_code" db:"response_status_code"`
	ExecutionDuration  int       `json:"execution_duration" db:"execution_duration"`
	Success            bool      `json:"success" db:"success"`
	Error              string    `json:"error" db:"error"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}
