package postmand

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestWebhook(t *testing.T) {
	var tests = []struct {
		kind            string
		request         Webhook
		expectedPayload string
	}{
		{
			"required fields",
			Webhook{},
			`{"content_type":"cannot be blank","delivery_attempt_timeout":"cannot be blank","max_delivery_attempts":"cannot be blank","name":"cannot be blank","retry_max_backoff":"cannot be blank","retry_min_backoff":"cannot be blank","url":"cannot be blank","valid_status_codes":"cannot be blank"}`,
		},
		{
			"Short name",
			Webhook{ID: uuid.New(), Name: "A", URL: "https://httpbin.org/post", ContentType: "application/json", ValidStatusCodes: pq.Int32Array{200, 201}, MaxDeliveryAttempts: 1, DeliveryAttemptTimeout: 1, RetryMinBackoff: 1, RetryMaxBackoff: 1},
			`{"name":"the length must be between 3 and 255"}`,
		},
		{
			"Long name",
			Webhook{ID: uuid.New(), Name: strings.Repeat("A", 300), URL: "https://httpbin.org/post", ContentType: "application/json", ValidStatusCodes: pq.Int32Array{200, 201}, MaxDeliveryAttempts: 1, DeliveryAttemptTimeout: 1, RetryMinBackoff: 1, RetryMaxBackoff: 1},
			`{"name":"the length must be between 3 and 255"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			err := tt.request.Validate()
			assert.NotNil(t, err)
			errorPayload, err := json.Marshal(err)
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedPayload, string(errorPayload))
		})
	}

	webhook := Webhook{
		ID:                     uuid.New(),
		Name:                   "AAA",
		URL:                    "https://httpbin.org/post",
		ContentType:            "application/json",
		ValidStatusCodes:       pq.Int32Array{200, 201},
		MaxDeliveryAttempts:    1,
		DeliveryAttemptTimeout: 1,
		RetryMinBackoff:        1,
		RetryMaxBackoff:        1,
	}
	err := webhook.Validate()
	assert.Nil(t, err)
}

func TestDelivery(t *testing.T) {
	var tests = []struct {
		kind            string
		request         Delivery
		expectedPayload string
	}{
		{
			"required fields",
			Delivery{},
			`{"webhook_id":"must be a valid UUID v4"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			err := tt.request.Validate()
			assert.NotNil(t, err)
			errorPayload, err := json.Marshal(err)
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedPayload, string(errorPayload))
		})
	}

	delivery := Delivery{
		ID:          uuid.New(),
		WebhookID:   uuid.New(),
		Payload:     `{"success": true}`,
		ScheduledAt: time.Now().UTC(),
		Status:      DeliveryStatusPending,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
	err := delivery.Validate()
	assert.Nil(t, err)
}
