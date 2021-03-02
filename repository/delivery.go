package repository

import (
	"github.com/allisson/postmand"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
)

// Delivery implements postmand.DeliveryRepository interface.
type Delivery struct {
	db *sqlx.DB
}

// Get returns postmand.Delivery by options filter.
func (d Delivery) Get(getOptions *postmand.RepositoryGetOptions) (*postmand.Delivery, error) {
	delivery := postmand.Delivery{}
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("*").From("deliveries")
	for key, value := range getOptions.Filters {
		sb.Where(sb.Equal(key, value))
	}
	sql, args := sb.Build()
	err := d.db.Get(&delivery, sql, args...)
	return &delivery, err
}

// List returns a slice of postmand.Delivery by options filter.
func (d Delivery) List(listOptions *postmand.RepositoryListOptions) ([]*postmand.Delivery, error) {
	deliveries := []*postmand.Delivery{}
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("*").From("deliveries").Limit(listOptions.Limit).Offset(listOptions.Offset)
	for key, value := range listOptions.Filters {
		sb.Where(sb.Equal(key, value))
	}
	if listOptions.OrderBy != "" {
		sb.OrderBy(listOptions.OrderBy)
	}
	sql, args := sb.Build()
	err := d.db.Select(&deliveries, sql, args...)
	return deliveries, err
}

// Create postmand.Delivery on database.
func (d Delivery) Create(delivery *postmand.Delivery) error {
	sqlStatement := `
		INSERT INTO deliveries (
			"id",
			"webhook_id",
			"payload",
			"scheduled_at",
			"delivery_attempts",
			"status",
			"created_at",
			"updated_at"
		)
		VALUES (
			:id,
			:webhook_id,
			:payload,
			:scheduled_at,
			:delivery_attempts,
			:status,
			:created_at,
			:updated_at
		)
	`
	_, err := d.db.NamedExec(sqlStatement, delivery)
	return err
}

// Update postmand.Delivery on database.
func (d Delivery) Update(delivery *postmand.Delivery) error {
	sqlStatement := `
		UPDATE deliveries
		SET webhook_id = :webhook_id,
			payload = :payload,
			scheduled_at = :scheduled_at,
			delivery_attempts = :delivery_attempts,
			status = :status,
			created_at = :created_at,
			updated_at = :updated_at
		WHERE id = :id
	`
	_, err := d.db.NamedExec(sqlStatement, delivery)
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
