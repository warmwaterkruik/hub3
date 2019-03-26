package fs

import (
	"fmt"

	"github.com/delving/webresource/pkg/domain"
)

// Storage  interacts with Google Cloud Storage
type Storage struct {
	generator domain.DerivativeGenerator
	basePath  string
	baseURL   string
}

// NewStorage  returns a new Google Cloud Storage
func NewStorage(basePath, baseURL string, generator domain.DerivativeGenerator) (*Storage, error) {
	var err error

	s := &Storage{
		baseURL:   baseURL,
		basePath:  basePath,
		generator: generator,
	}
	return s, err
}

// Add adds a WebResource to the  Storage
func (s Storage) Add(wr domain.WebResource) error {
	//path := wr.StoragePath()
	return fmt.Errorf("not implemented yet")
}

// Remove a WebResource from storage
func (s Storage) Remove(wr domain.WebResource) error {
	//path := wr.StoragePath()
	return fmt.Errorf("not implemented yet")
}

// RemoteURI returns the public path for the WebResource
func (s Storage) RemoteURI(wr domain.WebResource) string {
	return fmt.Sprintf("https://%s/static/%s", s.baseURL, wr.StoragePath())
}

// Available checks if the WebResource exists in the Storage Repository
func (s Storage) Available(wr domain.WebResource) (bool, error) {
	return true, fmt.Errorf("not implemented yet")
}

// Metadata returns the custom metadata for a Webresource
func (s Storage) Metadata(wr domain.WebResource) (map[string]string, error) {
	//path := wr.StoragePath()
	return nil, fmt.Errorf("not implemented yet")
}

// Generator returns a domain.DerivativeGenerator
func (s Storage) Generator() domain.DerivativeGenerator {
	return s.generator
}

// IsStale determines if the source is newer than the derivative. In that case,
// the derivatives should be discarded
func (s Storage) IsStale(src, drv domain.WebResource) (bool, error) {
	return true, fmt.Errorf("not implemented yet")
}

// GetByPrefix returns an array of all the Repository paths that match the path
func (s Storage) GetByPrefix(path string) ([]string, error) {
	// TODO implement this
	return nil, fmt.Errorf("not implemented yet")
}

// GetWebResource returns a WebResource with its []bytes Body from the storage.
// This body can be used for rendering or transformations
func (s Storage) GetWebResource(wr *domain.WebResource) (*domain.WebResource, error) {
	return nil, fmt.Errorf("not implemented yet")
}
