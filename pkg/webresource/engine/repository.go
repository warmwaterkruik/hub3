package engine

import (
	e "errors"
	"fmt"

	"github.com/delving/webresource/pkg/domain"
	"github.com/pkg/errors"
)

var (
	// ErrNotFound is used when a WebResource could not be found.
	ErrNotFound = e.New("WebResource not found")

	// ErrOperationNotAllowed is used when a Derivative is not allowed to be generated
	ErrOperationNotAllowed = e.New("Operation is not allowed")

	// ErrAccessDenied is used when Access to the source webresource is not allowed
	ErrAccessDenied = e.New("Access to WebResource is not allowed")
)

// Repository provides access to the storage
type Repository interface {
	// Adding
	Add(wr domain.WebResource) error

	// Listing
	Available(wr domain.WebResource) (bool, error)
	RemoteURI(wr domain.WebResource) string

	// Removing
	Remove(wr domain.WebResource) error

	// Getters
	Generator() domain.DerivativeGenerator
	GetWebResource(wr *domain.WebResource) (*domain.WebResource, error)

	GetByPrefix(path string) ([]string, error)
	IsStale(source, derivative domain.WebResource) (bool, error)
}

// SearchEngine is the main contract for interacting with the search engine for
// WebResources
type SearchEngine interface {
	Get(wr domain.WebResource) error
}

// Service provides engine operations
type Service interface {
	// Adding
	Add(wr domain.WebResource) error

	// Listing
	Available(wr domain.WebResource) (bool, error)
	RemoteURI(wr domain.WebResource) string

	// Removing
	Remove(wr domain.WebResource) error

	// Getters
	Generator() domain.DerivativeGenerator
	GetOrCreateRemoteURI(wr domain.WebResource) (string, error)

	//
	GetGraph(wr domain.WebResource) (string, error)
	IsStale(source, derivative domain.WebResource) (bool, error)

	// generating
	GenerateDeepZoom(wr domain.WebResource) error
}

type service struct {
	r Repository
}

// NewService creates engine service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}

// Add stores Reader with key path in the Storage Repository
func (s service) Add(wr domain.WebResource) error {
	return s.r.Add(wr)
}

// Available checks if the WebResource is available in the Storage Repository
func (s service) Available(wr domain.WebResource) (bool, error) {
	return s.r.Available(wr)
}

// RemoteURI returns the public URI for the WebResource
func (s service) RemoteURI(wr domain.WebResource) string {
	return s.r.RemoteURI(wr)
}

// Remove removes the WebResource from the storage Repository
func (s service) Remove(wr domain.WebResource) error {
	return s.r.Remove(wr)
}

// Generator returns a domain.DerivativeGenerator.
// This can be used for generating thumbnails
func (s service) Generator() domain.DerivativeGenerator {
	return s.r.Generator()
}

// GetOrCreateRemoteURI returns the RemoteURI and when the WebResource is not available
// tries to create it from the source WebResource
func (s service) GetOrCreateRemoteURI(wr domain.WebResource) (string, error) {
	ok, err := s.Available(wr)
	if err != nil {
		return "", errors.Wrap(err, "Unable to check if path is available")
	}

	// if available return remote public URI
	if ok {
		return s.RemoteURI(wr), nil
	}

	src := &domain.WebResource{
		OrgID:     wr.OrgID,
		Spec:      wr.Spec,
		SourceKey: wr.SourceKey,
		Kind:      domain.SOURCE,
	}
	src, err = s.r.GetWebResource(src)
	if err != nil {
		return "", errors.Wrapf(err, "Unable to get WebResource %s", src.StoragePath())
	}
	switch wr.Kind {
	case domain.THUMBNAIL:
		wr.Body, err = s.r.Generator().Thumbnail(src.Body, wr.Operation)
		if err != nil {
			return "", errors.Wrapf(err, "Unable to generate thumbnail for: %s", wr.Operation)
		}
		err = s.Add(wr)
		if err != nil {
			return "", errors.Wrapf(err, "Unable to store Thumbnail: %s", wr.StoragePath())
		}
	case domain.DEEPZOOM:
		err := s.GenerateDeepZoom(wr)
		if err != nil {
			return "", errors.Wrapf(err, "Unable to store DeepZoom: %v", wr.StoragePath())
		}

	default:
		return "", fmt.Errorf("type of conversion %s is not supported", wr.Kind)
	}

	return s.RemoteURI(wr), err

}

// GetGraph returns a JSON-LD RDF graph for all the entries for path.
// In case the path ends with '__', then the graph will return all entries that
// are prefixed by this path. It should not contain any duplicate entries when
// multiple versions with different mime-types are stored in the Repository
func (s service) GetGraph(wr domain.WebResource) (string, error) {
	// TODO implement this
	return "", fmt.Errorf("not implemented yet")
}

// IsStale determines if the source is newer than the derivative.
func (s service) IsStale(source, derivative domain.WebResource) (bool, error) {
	return s.r.IsStale(source, derivative)
}

// GenerateDeepZoom generates all DeepZoom tiles and XML descriptor and stores them
// in the storage repository
func (s service) GenerateDeepZoom(wr domain.WebResource) error {
	return fmt.Errorf("not implemented")
}
