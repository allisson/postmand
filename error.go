package postmand

import "errors"

var (
	// ErrWebhookNotFound is returned by any operation that can't load a webhook.
	ErrWebhookNotFound = errors.New("webhook_not_found")
	// ErrDeliveryNotFound is returned by any operation that can't load a delivery.
	ErrDeliveryNotFound = errors.New("delivery_not_found")
	// ErrDeliveryAttemptNotFound is returned by any operation that can't load a delivery attempt.
	ErrDeliveryAttemptNotFound = errors.New("delivery_attempt_not_found")
)
