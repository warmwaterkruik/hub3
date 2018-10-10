package engine

import "github.com/delving/rapid-saas/pkg/domain"

// Repository provides access to the storage
type Repository interface {
	Add(sr *domain.StoreRequest) error
	QueuePostHook(sr *domain.StoreRequest) error
	Flush() error
}

// Service is the storage abstraction layer
type Service interface {
	Add(sr *domain.StoreRequest) error
	QueuePostHook(sr *domain.StoreRequest) error
	Flush() error
}

type service struct {
	r Repository
}

// NewService creates engine service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}

// Add adds a StoreRquest to the Repository Storage
func (s service) Add(sr *domain.StoreRequest) error {
	return s.r.Add(sr)
}

// QueuePostHook adds the StoreRequest to a Queue for processing after the StoreRequest
// is stored in the Repository Storage
func (s service) QueuePostHook(sr *domain.StoreRequest) error {
	return s.r.QueuePostHook(sr)
}

// Flush flushes all from the queue to the Repository Storage
func (s service) Flush() error {
	return s.r.Flush()
}
