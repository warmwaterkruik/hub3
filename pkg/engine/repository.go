package engine

import (
	"context"

	"github.com/delving/rapid-saas/pkg/domain"
)

// Repository provides access to the storage
type Repository interface {
	Add(sr *domain.StoreRequest) error
	QueuePostHook(sr *domain.StoreRequest) error
	Flush() error
	Reset(ctx context.Context, index string) error
}

// Service is the storage abstraction layer
type Service interface {
	Add(sr *domain.StoreRequest) error
	QueuePostHook(sr *domain.StoreRequest) error
	Flush() error
	Reset(ctx context.Context, index string) error
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

// Reset drops all storage and re-initializes it. Usefull for automated regression testing
func (s service) Reset(ctx context.Context, index string) error {
	return s.r.Reset(ctx, index)
}
