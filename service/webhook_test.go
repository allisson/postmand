package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/allisson/postmand"
	"github.com/allisson/postmand/mocks"
)

func TestWebhook(t *testing.T) {
	ctx := context.Background()

	t.Run("Get", func(t *testing.T) {
		webhookRepository := &mocks.WebhookRepository{}
		webhookService := NewWebhook(webhookRepository)
		expectedWebhook := &postmand.Webhook{ID: uuid.New()}
		getOptions := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": expectedWebhook.ID}}

		webhookRepository.On("Get", mock.Anything, getOptions).Return(expectedWebhook, nil)
		webhook, err := webhookService.Get(ctx, getOptions)
		assert.Nil(t, err)
		assert.Equal(t, expectedWebhook, webhook)
		webhookRepository.AssertExpectations(t)
	})

	t.Run("List", func(t *testing.T) {
		webhookRepository := &mocks.WebhookRepository{}
		webhookService := NewWebhook(webhookRepository)
		expectedWebhook := &postmand.Webhook{ID: uuid.New()}
		listOptions := postmand.RepositoryListOptions{Filters: map[string]interface{}{"id": expectedWebhook.ID}, Limit: 1, Offset: 0}

		webhookRepository.On("List", mock.Anything, listOptions).Return([]*postmand.Webhook{expectedWebhook}, nil)
		webhooks, err := webhookService.List(ctx, listOptions)
		assert.Nil(t, err)
		assert.Equal(t, expectedWebhook, webhooks[0])
		webhookRepository.AssertExpectations(t)
	})

	t.Run("Create", func(t *testing.T) {
		webhookRepository := &mocks.WebhookRepository{}
		webhookService := NewWebhook(webhookRepository)
		webhook := &postmand.Webhook{ID: uuid.New()}

		webhookRepository.On("Create", mock.Anything, webhook).Return(nil)
		err := webhookService.Create(ctx, webhook)
		assert.Nil(t, err)
		webhookRepository.AssertExpectations(t)
	})

	t.Run("Update", func(t *testing.T) {
		webhookRepository := &mocks.WebhookRepository{}
		webhookService := NewWebhook(webhookRepository)
		webhook := &postmand.Webhook{ID: uuid.New()}

		getOptions := postmand.RepositoryGetOptions{Filters: map[string]interface{}{"id": webhook.ID}}
		webhookRepository.On("Get", mock.Anything, getOptions).Return(webhook, nil)
		webhookRepository.On("Update", mock.Anything, webhook).Return(nil)
		err := webhookService.Update(ctx, webhook)
		assert.Nil(t, err)
		webhookRepository.AssertExpectations(t)
	})

	t.Run("Delete", func(t *testing.T) {
		webhookRepository := &mocks.WebhookRepository{}
		webhookService := NewWebhook(webhookRepository)
		webhook := &postmand.Webhook{ID: uuid.New()}

		webhookRepository.On("Delete", mock.Anything, webhook.ID).Return(nil)
		err := webhookService.Delete(ctx, webhook.ID)
		assert.Nil(t, err)
		webhookRepository.AssertExpectations(t)
	})
}
