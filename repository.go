package postmand

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
}

// WebhookRepository is the interface that will be used to iterate with the Webhook data.
type WebhookRepository interface {
	Get(getOptions *RepositoryGetOptions) (*Webhook, error)
	List(listOptions *RepositoryListOptions) ([]*Webhook, error)
	Create(webhook *Webhook) error
	Update(webhook *Webhook) error
	Delete(id ID) error
}

// DeliveryRepository is the interface that will be used to iterate with the Delivery data.
type DeliveryRepository interface {
	Get(getOptions *RepositoryGetOptions) (*Delivery, error)
	List(listOptions *RepositoryListOptions) ([]*Delivery, error)
	Create(delivery *Delivery) error
	Update(delivery *Delivery) error
	Delete(id ID) error
}
