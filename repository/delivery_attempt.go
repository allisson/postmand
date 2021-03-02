package repository

import (
	"github.com/allisson/postmand"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
)

// DeliveryAttempt implements postmand.DeliveryAttemptRepository interface.
type DeliveryAttempt struct {
	db *sqlx.DB
}

// Get returns postmand.DeliveryAttempt by options filter.
func (d DeliveryAttempt) Get(getOptions *postmand.RepositoryGetOptions) (*postmand.DeliveryAttempt, error) {
	deliveryAttempt := postmand.DeliveryAttempt{}
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("*").From("delivery_attempts")
	for key, value := range getOptions.Filters {
		sb.Where(sb.Equal(key, value))
	}
	sql, args := sb.Build()
	err := d.db.Get(&deliveryAttempt, sql, args...)
	return &deliveryAttempt, err
}

// List returns a slice of postmand.DeliveryAttempt by options filter.
func (d DeliveryAttempt) List(listOptions *postmand.RepositoryListOptions) ([]*postmand.DeliveryAttempt, error) {
	deliveryAttempts := []*postmand.DeliveryAttempt{}
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("*").From("delivery_attempts").Limit(listOptions.Limit).Offset(listOptions.Offset)
	for key, value := range listOptions.Filters {
		sb.Where(sb.Equal(key, value))
	}
	if listOptions.OrderBy != "" {
		sb.OrderBy(listOptions.OrderBy)
	}
	sql, args := sb.Build()
	err := d.db.Select(&deliveryAttempts, sql, args...)
	return deliveryAttempts, err
}

// Create postmand.DeliveryAttempt on database.
func (d DeliveryAttempt) Create(deliveryAttempt *postmand.DeliveryAttempt) error {
	sqlStatement := `
		INSERT INTO delivery_attempts (
			"id",
			"webhook_id",
			"delivery_id",
			"response_headers",
			"response_body",
			"response_status_code",
			"execution_duration",
			"success",
			"error",
			"created_at"
		)
		VALUES (
			:id,
			:webhook_id,
			:delivery_id,
			:response_headers,
			:response_body,
			:response_status_code,
			:execution_duration,
			:success,
			:error,
			:created_at
		)
	`
	_, err := d.db.NamedExec(sqlStatement, deliveryAttempt)
	return err
}

// NewDeliveryAttempt returns postmand.DeliveryAttempt with db connection.
func NewDeliveryAttempt(db *sqlx.DB) *DeliveryAttempt {
	return &DeliveryAttempt{db: db}
}
