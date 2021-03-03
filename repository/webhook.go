package repository

import (
	"github.com/allisson/postmand"
	"github.com/jmoiron/sqlx"
)

// Webhook implements postmand.WebhookRepository interface.
type Webhook struct {
	db *sqlx.DB
}

// Get returns postmand.Webhook by options filter.
func (w Webhook) Get(getOptions postmand.RepositoryGetOptions) (*postmand.Webhook, error) {
	webhook := postmand.Webhook{}
	sql, args := getQuery("webhooks", getOptions)
	err := w.db.Get(&webhook, sql, args...)
	return &webhook, err
}

// List returns a slice of postmand.Webhook by options filter.
func (w Webhook) List(listOptions postmand.RepositoryListOptions) ([]*postmand.Webhook, error) {
	webhooks := []*postmand.Webhook{}
	sql, args := listQuery("webhooks", listOptions)
	err := w.db.Select(&webhooks, sql, args...)
	return webhooks, err
}

// Create postmand.Webhook on database.
func (w Webhook) Create(webhook *postmand.Webhook) error {
	sql, args := insertQuery("webhooks", webhook)
	_, err := w.db.Exec(sql, args...)
	return err
}

// Update postmand.Webhook on database.
func (w Webhook) Update(webhook *postmand.Webhook) error {
	sql, args := updateQuery("webhooks", webhook)
	_, err := w.db.Exec(sql, args...)
	return err
}

// Delete postmand.Webhook on database.
func (w Webhook) Delete(id postmand.ID) error {
	sqlStatement := `
		DELETE FROM webhooks WHERE id = $1
	`
	_, err := w.db.Exec(sqlStatement, id)
	return err
}

// NewWebhook returns postmand.Webhook with db connection.
func NewWebhook(db *sqlx.DB) *Webhook {
	return &Webhook{db: db}
}
