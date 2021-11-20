package service

import (
	"context"

	"go.uber.org/zap"

	"github.com/allisson/postmand"
)

// Migration implements postmand.MigrationService interface.
type Migration struct {
	migrationRepo postmand.MigrationRepository
	logger        *zap.Logger
}

// Run database migrations.
func (m Migration) Run(ctx context.Context) error {
	m.logger.Info("migration-started")
	if err := m.migrationRepo.Run(ctx); err != nil {
		m.logger.Error("migration-error", zap.Error(err))
	}
	m.logger.Info("migration-completed")
	return nil
}

// NewMigration will create an implementation of postmand.MigrationService.
func NewMigration(migrationRepo postmand.MigrationRepository, logger *zap.Logger) *Migration {
	return &Migration{
		migrationRepo: migrationRepo,
		logger:        logger,
	}
}
