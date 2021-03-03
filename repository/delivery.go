package repository

import (
	"github.com/allisson/postmand"
	"github.com/jmoiron/sqlx"
)

// Delivery implements postmand.DeliveryRepository interface.
type Delivery struct {
	db *sqlx.DB
}

// Get returns postmand.Delivery by options filter.
func (d Delivery) Get(getOptions postmand.RepositoryGetOptions) (*postmand.Delivery, error) {
	delivery := postmand.Delivery{}
	sql, args := getQuery("deliveries", getOptions)
	err := d.db.Get(&delivery, sql, args...)
	return &delivery, err
}

// List returns a slice of postmand.Delivery by options filter.
func (d Delivery) List(listOptions postmand.RepositoryListOptions) ([]*postmand.Delivery, error) {
	deliveries := []*postmand.Delivery{}
	sql, args := listQuery("deliveries", listOptions)
	err := d.db.Select(&deliveries, sql, args...)
	return deliveries, err
}

// Create postmand.Delivery on database.
func (d Delivery) Create(delivery *postmand.Delivery) error {
	sql, args := insertQuery("deliveries", delivery)
	_, err := d.db.Exec(sql, args...)
	return err
}

// Update postmand.Delivery on database.
func (d Delivery) Update(delivery *postmand.Delivery) error {
	sql, args := updateQuery("deliveries", delivery)
	_, err := d.db.Exec(sql, args...)
	return err
}

// Delete postmand.Delivery on database.
func (d Delivery) Delete(id postmand.ID) error {
	sqlStatement := `
		DELETE FROM deliveries WHERE id = $1
	`
	_, err := d.db.Exec(sqlStatement, id)
	return err
}

// NewDelivery returns postmand.Delivery with db connection.
func NewDelivery(db *sqlx.DB) *Delivery {
	return &Delivery{db: db}
}
