package service

import (
	"context"

	"github.com/crypitor/postmand"
)

// Ping implements postmand.PingService interface.
type Ping struct {
	pingRepository postmand.PingRepository
}

// Run ping operation against the database.
func (p Ping) Run(ctx context.Context) error {
	return p.pingRepository.Run(ctx)
}

// NewPing will create an implementation of postmand.PingService.
func NewPing(pingRepository postmand.PingRepository) *Ping {
	return &Ping{pingRepository: pingRepository}
}
