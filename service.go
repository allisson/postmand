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

// DeliveryService is the interface that will be used to perform operations with deliveries.
type DeliveryService interface {
	Get(ctx context.Context, getOptions RepositoryGetOptions) (*Delivery, error)
	List(ctx context.Context, listOptions RepositoryListOptions) ([]*Delivery, error)
	Create(ctx context.Context, delivery *Delivery) error
	Update(ctx context.Context, delivery *Delivery) error
	Delete(ctx context.Context, id ID) error
}

// DeliveryAttemptService is the interface that will be used to perform operations with delivery attempt.
type DeliveryAttemptService interface {
	Get(ctx context.Context, getOptions RepositoryGetOptions) (*DeliveryAttempt, error)
	List(ctx context.Context, listOptions RepositoryListOptions) ([]*DeliveryAttempt, error)
}

// PingService is the interface that will be used to perform ping operation against database.
type PingService interface {
	Run(ctx context.Context) error
}
