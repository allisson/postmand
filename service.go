package postmand

import "context"

// WorkerService is the interface that will be used on workers to dispatch webhooks.
type WorkerService interface {
	Run(ctx context.Context)
	Shutdown(ctx context.Context)
}

// MigrationService is the interface that will be used to execute database migrations.
type MigrationService interface {
	Run(ctx context.Context) error
}
