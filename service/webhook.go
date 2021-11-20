package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/allisson/postmand"
)

// Webhook implements postmand.WebhookService interface.
type Webhook struct {
	webhookRepository postmand.WebhookRepository
}

// Get returns postmand.Webhook by options filter.
func (w Webhook) Get(ctx context.Context, getOptions postmand.RepositoryGetOptions) (*postmand.Webhook, error) {
	return w.webhookRepository.Get(ctx, getOptions)
}

// List returns a slice of postmand.Webhook by options filter.
func (w Webhook) List(ctx context.Context, listOptions postmand.RepositoryListOptions) ([]*postmand.Webhook, error) {
	return w.webhookRepository.List(ctx, listOptions)
}

// Create postmand.Webhook on database.
func (w Webhook) Create(ctx context.Context, webhook *postmand.Webhook) error {
	now := time.Now().UTC()
	webhook.ID = uuid.New()
	webhook.CreatedAt = now
	webhook.UpdatedAt = now
	return w.webhookRepository.Create(ctx, webhook)
}

// Update postmand.Webhook on database.
func (w Webhook) Update(ctx context.Context, webhook *postmand.Webhook) error {
	getOptions := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": webhook.ID}}
	storedWebhook, err := w.webhookRepository.Get(ctx, getOptions)
	if err != nil {
		return err
	}
	webhook.CreatedAt = storedWebhook.CreatedAt
	webhook.UpdatedAt = time.Now().UTC()
	return w.webhookRepository.Update(ctx, webhook)
}

// Delete postmand.Webhook on database.
func (w Webhook) Delete(ctx context.Context, id postmand.ID) error {
	return w.webhookRepository.Delete(ctx, id)
}

// NewWebhook will create an implementation of postmand.WebhookService.
func NewWebhook(webhookRepository postmand.WebhookRepository) *Webhook {
	return &Webhook{webhookRepository: webhookRepository}
}
