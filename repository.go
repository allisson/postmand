package postmand

import "context"

// RepositoryGetOptions contains options used in the Get methods.
type RepositoryGetOptions struct {
	Filters map[string]interface{}
}

// RepositoryListOptions contains options used in the List methods.
type RepositoryListOptions struct {
	Filters map[string]interface{}
	Limit   int
	Offset  int
	OrderBy string
	Order   string
}

// WebhookRepository is the interface that will be used to iterate with the Webhook data.
type WebhookRepository interface {
	Get(ctx context.Context, getOptions RepositoryGetOptions) (*Webhook, error)
	List(ctx context.Context, listOptions RepositoryListOptions) ([]*Webhook, error)
	Create(ctx context.Context, webhook *Webhook) error
	Update(ctx context.Context, webhook *Webhook) error
	Delete(ctx context.Context, id ID) error
}

// DeliveryRepository is the interface that will be used to iterate with the Delivery data.
type DeliveryRepository interface {
	Get(ctx context.Context, getOptions RepositoryGetOptions) (*Delivery, error)
	List(ctx context.Context, listOptions RepositoryListOptions) ([]*Delivery, error)
	Create(ctx context.Context, delivery *Delivery) error
	Update(ctx context.Context, delivery *Delivery) error
	Delete(ctx context.Context, id ID) error
	Dispatch(ctx context.Context) (*DeliveryAttempt, error)
}

// DeliveryAttemptRepository is the interface that will be used to iterate with the DeliveryAttempt data.
type DeliveryAttemptRepository interface {
	Get(ctx context.Context, getOptions RepositoryGetOptions) (*DeliveryAttempt, error)
	List(ctx context.Context, listOptions RepositoryListOptions) ([]*DeliveryAttempt, error)
	Create(ctx context.Context, deliveryAttempt *DeliveryAttempt) error
}

// MigrationRepository is the interface that will be used to run database migrations.
type MigrationRepository interface {
	Run(ctx context.Context) error
}
