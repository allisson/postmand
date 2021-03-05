package repository

import (
	"context"
	"database/sql"

	"github.com/allisson/postmand"
	"github.com/jmoiron/sqlx"
)

// Webhook implements postmand.WebhookRepository interface.
type Webhook struct {
	db *sqlx.DB
}

// Get returns postmand.Webhook by options filter.
func (w Webhook) Get(ctx context.Context, getOptions postmand.RepositoryGetOptions) (*postmand.Webhook, error) {
	webhook := postmand.Webhook{}
	query, args := getQuery("webhooks", getOptions)
	err := w.db.GetContext(ctx, &webhook, query, args...)
	if err == sql.ErrNoRows {
		return &webhook, postmand.ErrWebhookNotFound
	}
	return &webhook, err
}

// List returns a slice of postmand.Webhook by options filter.
func (w Webhook) List(ctx context.Context, listOptions postmand.RepositoryListOptions) ([]*postmand.Webhook, error) {
	webhooks := []*postmand.Webhook{}
	query, args := listQuery("webhooks", listOptions)
	err := w.db.SelectContext(ctx, &webhooks, query, args...)
	return webhooks, err
}

// Create postmand.Webhook on database.
func (w Webhook) Create(ctx context.Context, webhook *postmand.Webhook) error {
	query, args := insertQuery("webhooks", webhook)
	_, err := w.db.ExecContext(ctx, query, args...)
	return err
}

// Update postmand.Webhook on database.
func (w Webhook) Update(ctx context.Context, webhook *postmand.Webhook) error {
	query, args := updateQuery("webhooks", webhook.ID, webhook)
	_, err := w.db.ExecContext(ctx, query, args...)
	return err
}

// Delete postmand.Webhook on database.
func (w Webhook) Delete(ctx context.Context, id postmand.ID) error {
	query := `
		DELETE FROM webhooks WHERE id = $1
	`
	_, err := w.db.ExecContext(ctx, query, id)
	return err
}

// NewWebhook will create an implementation of postmand.WebhookRepository.
func NewWebhook(db *sqlx.DB) *Webhook {
	return &Webhook{db: db}
}
