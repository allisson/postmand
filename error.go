package postmand

import "errors"

var (
	// ErrWebhookNotFound is returned by any operation that can't load a webhook.
	ErrWebhookNotFound = errors.New("webhook_not_found")
)
