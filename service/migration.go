package service

import (
	"context"

	"github.com/allisson/postmand"
)

// Migration implements postmand.MigrationService interface.
type Migration struct {
	migrationRepo postmand.MigrationRepository
}

// Run database migrations.
func (m Migration) Run(ctx context.Context) error {
	return m.migrationRepo.Run(ctx)
}

// NewMigration will create an implementation of postmand.MigrationService.
func NewMigration(migrationRepo postmand.MigrationRepository) *Migration {
	return &Migration{migrationRepo: migrationRepo}
}
