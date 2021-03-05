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

// WebhookService is the interface that will be used to perform operations with webhooks.
type WebhookService interface {
	Get(ctx context.Context, getOptions RepositoryGetOptions) (*Webhook, error)
	List(ctx context.Context, listOptions RepositoryListOptions) ([]*Webhook, error)
	Create(ctx context.Context, webhook *Webhook) error
	Update(ctx context.Context, webhook *Webhook) error
	Delete(ctx context.Context, id ID) error
}
