package repository

import (
	"github.com/allisson/postmand"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
)

// Webhook implements postmand.WebhookRepository interface.
type Webhook struct {
	db *sqlx.DB
}

// Get returns postmand.Webhook by options filter.
func (w Webhook) Get(getOptions *postmand.RepositoryGetOptions) (*postmand.Webhook, error) {
	webhook := postmand.Webhook{}
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("*").From("webhooks")
	for key, value := range getOptions.Filters {
		sb.Where(sb.Equal(key, value))
	}
	sql, args := sb.Build()
	err := w.db.Get(&webhook, sql, args...)
	return &webhook, err
}

// List returns a slice of postmand.Webhook by options filter.
func (w Webhook) List(listOptions *postmand.RepositoryListOptions) ([]*postmand.Webhook, error) {
	webhooks := []*postmand.Webhook{}
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("*").From("webhooks").Limit(listOptions.Limit).Offset(listOptions.Offset)
	for key, value := range listOptions.Filters {
		sb.Where(sb.Equal(key, value))
	}
	if listOptions.OrderBy != "" {
		sb.OrderBy(listOptions.OrderBy)
	}
	sql, args := sb.Build()
	err := w.db.Select(&webhooks, sql, args...)
	return webhooks, err
}

// Create webhook on database.
func (w Webhook) Create(webhook *postmand.Webhook) error {
	sqlStatement := `
		INSERT INTO webhooks (
			"id",
			"name",
			"url",
			"content_type",
			"valid_status_codes",
			"secret_token",
			"active",
			"max_delivery_attempts",
			"delivery_attempt_timeout",
			"retry_min_backoff",
			"retry_max_backoff",
			"created_at",
			"updated_at"
		)
		VALUES (
			:id,
			:name,
			:url,
			:content_type,
			:valid_status_codes,
			:secret_token,
			:active,
			:max_delivery_attempts,
			:delivery_attempt_timeout,
			:retry_min_backoff,
			:retry_max_backoff,
			:created_at,
			:updated_at
		)
	`
	_, err := w.db.NamedExec(sqlStatement, webhook)
	return err
}

// Update webhook on database.
func (w Webhook) Update(webhook *postmand.Webhook) error {
	sqlStatement := `
		UPDATE webhooks
		SET name = :name,
			url = :url,
			content_type = :content_type,
			valid_status_codes = :valid_status_codes,
			secret_token = :secret_token,
			active = :active,
			max_delivery_attempts = :max_delivery_attempts,
			delivery_attempt_timeout = :delivery_attempt_timeout,
			retry_min_backoff = :retry_min_backoff,
			retry_max_backoff = :retry_max_backoff,
			created_at = :created_at,
			updated_at = :updated_at
		WHERE id = :id
	`
	_, err := w.db.NamedExec(sqlStatement, webhook)
	return err
}

// Delete webhook on database.
func (w Webhook) Delete(id postmand.ID) error {
	sqlStatement := `
		DELETE FROM webhooks WHERE id = $1
	`
	_, err := w.db.Exec(sqlStatement, id)
	return err
}

// NewWebhook returns Webhook with db connection.
func NewWebhook(db *sqlx.DB) *Webhook {
	return &Webhook{db: db}
}
