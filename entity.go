package postmand

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/google/uuid"
)

const (
	// DeliveryStatusPending represents the delivery pending status
	DeliveryStatusPending = "pending"
	// DeliveryStatusFailed represents the delivery failed status
	DeliveryStatusFailed = "failed"
	// DeliveryStatusCompleted represents the delivery completed status
	DeliveryStatusCompleted = "completed"
)

// Webhook represents a webhook in the system.
type Webhook struct {
	ID                     uuid.UUID `json:"id" db:"id"`
	Name                   string    `json:"name" db:"name"`
	URL                    string    `json:"url" db:"url"`
	ContentType            string    `json:"content_type" db:"content_type"`
	SecretToken            string    `json:"secret_token" db:"secret_token"`
	MaxDeliveryAttempts    int       `json:"max_delivery_attempts" db:"max_delivery_attempts"`
	DeliveryAttemptTimeout int       `json:"delivery_attempt_timeout" db:"delivery_attempt_timeout"`
	RetryMinBackoff        int       `json:"retry_min_backoff" db:"retry_min_backoff"`
	RetryMaxBackoff        int       `json:"retry_max_backoff" db:"retry_max_backoff"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`
}

// Validate implements ozzo validation Validatable interface
func (w Webhook) Validate() error {
	return validation.ValidateStruct(&w,
		validation.Field(&w.ID, validation.Required, is.UUIDv4),
		validation.Field(&w.Name, validation.Required, validation.Length(3, 255)),
		validation.Field(&w.URL, validation.Required, is.URL),
		validation.Field(&w.ContentType, validation.Required, validation.In("application/x-www-form-urlencoded", "application/json")),
		validation.Field(&w.MaxDeliveryAttempts, validation.Required, validation.Min(1)),
		validation.Field(&w.DeliveryAttemptTimeout, validation.Required, validation.Min(1)),
		validation.Field(&w.RetryMinBackoff, validation.Required, validation.Min(1)),
		validation.Field(&w.RetryMaxBackoff, validation.Required, validation.Min(1)),
	)
}

// Delivery represents a payload that must be delivery using webhook context.
type Delivery struct {
	ID               uuid.UUID `json:"id" db:"id"`
	WebhookID        uuid.UUID `json:"webhook_id" db:"webhook_id"`
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
		validation.Field(&d.ID, validation.Required, is.UUIDv4),
		validation.Field(&d.WebhookID, validation.Required, is.UUIDv4),
		validation.Field(&d.ScheduledAt, validation.Required),
		validation.Field(&d.Status, validation.Required, validation.In("pending", "completed", "failed")),
	)
}

// DeliveryAttempt represents a delivery attempt.
type DeliveryAttempt struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	WebhookID          uuid.UUID `json:"destination_id" db:"destination_id"`
	DeliveryID         uuid.UUID `json:"delivery_id" db:"delivery_id"`
	ResponseHeaders    string    `json:"response_headers" db:"response_headers"`
	ResponseBody       string    `json:"response_body" db:"response_body"`
	ResponseStatusCode int       `json:"response_status_code" db:"response_status_code"`
	ExecutionDuration  int       `json:"execution_duration" db:"execution_duration"`
	Success            bool      `json:"success" db:"success"`
	Error              string    `json:"error" db:"error"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}
