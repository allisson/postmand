package repository

import (
	"github.com/allisson/postmand"
	"github.com/jmoiron/sqlx"
)

// DeliveryAttempt implements postmand.DeliveryAttemptRepository interface.
type DeliveryAttempt struct {
	db *sqlx.DB
}

// Get returns postmand.DeliveryAttempt by options filter.
func (d DeliveryAttempt) Get(getOptions postmand.RepositoryGetOptions) (*postmand.DeliveryAttempt, error) {
	deliveryAttempt := postmand.DeliveryAttempt{}
	sql, args := getQuery("delivery_attempts", getOptions)
	err := d.db.Get(&deliveryAttempt, sql, args...)
	return &deliveryAttempt, err
}

// List returns a slice of postmand.DeliveryAttempt by options filter.
func (d DeliveryAttempt) List(listOptions postmand.RepositoryListOptions) ([]*postmand.DeliveryAttempt, error) {
	deliveryAttempts := []*postmand.DeliveryAttempt{}
	sql, args := listQuery("delivery_attempts", listOptions)
	err := d.db.Select(&deliveryAttempts, sql, args...)
	return deliveryAttempts, err
}

// Create postmand.DeliveryAttempt on database.
func (d DeliveryAttempt) Create(deliveryAttempt *postmand.DeliveryAttempt) error {
	sql, args := insertQuery("delivery_attempts", deliveryAttempt)
	_, err := d.db.Exec(sql, args...)
	return err
}

// NewDeliveryAttempt returns postmand.DeliveryAttempt with db connection.
func NewDeliveryAttempt(db *sqlx.DB) *DeliveryAttempt {
	return &DeliveryAttempt{db: db}
}
