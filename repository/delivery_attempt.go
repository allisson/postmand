package repository

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/crypitor/postmand"
)

// DeliveryAttempt implements postmand.DeliveryAttemptRepository interface.
type DeliveryAttempt struct {
	db *sqlx.DB
}

// Get returns postmand.DeliveryAttempt by options filter.
func (d DeliveryAttempt) Get(ctx context.Context, getOptions postmand.RepositoryGetOptions) (*postmand.DeliveryAttempt, error) {
	deliveryAttempt := postmand.DeliveryAttempt{}
	query, args := getQuery("delivery_attempts", getOptions)
	err := d.db.GetContext(ctx, &deliveryAttempt, query, args...)
	if err == sql.ErrNoRows {
		return &deliveryAttempt, postmand.ErrDeliveryAttemptNotFound
	}
	return &deliveryAttempt, err
}

// List returns a slice of postmand.DeliveryAttempt by options filter.
func (d DeliveryAttempt) List(ctx context.Context, listOptions postmand.RepositoryListOptions) ([]*postmand.DeliveryAttempt, error) {
	deliveryAttempts := []*postmand.DeliveryAttempt{}
	query, args := listQuery("delivery_attempts", listOptions)
	err := d.db.SelectContext(ctx, &deliveryAttempts, query, args...)
	return deliveryAttempts, err
}

// Create postmand.DeliveryAttempt on database.
func (d DeliveryAttempt) Create(ctx context.Context, deliveryAttempt *postmand.DeliveryAttempt) error {
	query, args := insertQuery("delivery_attempts", deliveryAttempt)
	_, err := d.db.ExecContext(ctx, query, args...)
	return err
}

// NewDeliveryAttempt returns postmand.DeliveryAttempt with db connection.
func NewDeliveryAttempt(db *sqlx.DB) *DeliveryAttempt {
	return &DeliveryAttempt{db: db}
}
