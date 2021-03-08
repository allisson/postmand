package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// Ping implements postmand.PingRepository interface.
type Ping struct {
	db *sqlx.DB
}

// Run ping operation against the database.
func (p Ping) Run(ctx context.Context) error {
	return p.db.PingContext(ctx)
}

// NewPing will create an implementation of postmand.PingRepository.
func NewPing(db *sqlx.DB) *Ping {
	return &Ping{db: db}
}
