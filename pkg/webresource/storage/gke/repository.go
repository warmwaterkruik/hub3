package gke

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/storage"
	"github.com/delving/webresource/pkg/domain"
	"github.com/pkg/errors"
)

// Storage  interacts with Google Cloud Storage
type Storage struct {
	client     *storage.Client
	bucketName string
	bucket     *storage.BucketHandle
	generator  domain.DerivativeGenerator

	ctx context.Context
}

// NewStorage  returns a new Google Cloud Storage
func NewStorage(bucketName, projectID string, generator domain.DerivativeGenerator) (*Storage, error) {
	var err error

	s := &Storage{
		ctx:        context.Background(),
		bucketName: bucketName,
		generator:  generator,
	}

	// Creates a client.
	s.client, err = storage.NewClient(s.ctx)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return nil, err
	}

	// Creates a Bucket instance.
	s.bucket = s.client.Bucket(bucketName)

	if _, err = s.bucket.Attrs(s.ctx); err != nil {
		// Creates the new bucket.
		attrs := &storage.BucketAttrs{
			//StorageClass: "MULTI-REGIONAL",
			Location: "europe-west1",
		}

		if err := s.bucket.Create(s.ctx, projectID, attrs); err != nil {
			log.Printf("Failed to create bucket: %v", err)
			return nil, errors.Wrap(err, "Failed to create bucket")
		}

		log.Printf("Bucket %v created.\n", bucketName)
	}

	return s, nil
}

// Add adds a WebResource to the  Storage
func (s Storage) Add(wr domain.WebResource) error {
	path := wr.StoragePath()
	wc := s.bucket.Object(path).NewWriter(s.ctx)
	wc.ContentType = wr.MimeType
	wc.Metadata = wr.GetMetadata("x-goog-meta-")

	_, err := wc.Write(wr.Body)
	if err != nil {
		return errors.Wrapf(err, "Unable to write data to bucket %s, file %s", s.bucketName, path)
	}
	if err := wc.Close(); err != nil {
		return errors.Wrapf(err, "unable to close bucket %q, file %q", s.bucketName, path)
	}

	switch wr.Kind {
	case domain.DEEPZOOM, domain.THUMBNAIL:
		err := s.bucket.Object(path).ACL().Set(s.ctx, storage.AllUsers, storage.RoleReader)
		if err != nil {
			return errors.Wrapf(err, "putACLRule: unable to save ACL rule for bucket %q, file %q", s.bucketName, path)
		}
	}

	return nil
}

// Remove a WebResource from storage
func (s Storage) Remove(wr domain.WebResource) error {
	path := wr.StoragePath()
	if err := s.bucket.Object(path).Delete(s.ctx); err != nil {
		return errors.Wrapf(err, "deleteFiles: unable to delete bucket %q, file %q", s.bucketName, path)
	}
	return nil
}

// RemoteURI returns the public path for the WebResource
func (s Storage) RemoteURI(wr domain.WebResource) string {
	return fmt.Sprintf("https://%s/storage.googleapis.com/%s", s.bucketName, wr.StoragePath())
}

// Available checks if the WebResource exists in the Storage Repository
func (s Storage) Available(wr domain.WebResource) (bool, error) {
	_, err := s.Metadata(wr)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Metadata returns the custom metadata for a Webresource
func (s Storage) Metadata(wr domain.WebResource) (map[string]string, error) {
	path := wr.StoragePath()
	obj, err := s.bucket.Object(path).Attrs(s.ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "statFile: unable to stat file from bucket %q, file %q", s.bucketName, path)
	}
	return obj.Metadata, nil
}

// Generator returns a domain.DerivativeGenerator
func (s Storage) Generator() domain.DerivativeGenerator {
	return s.generator
}

// IsStale determines if the source is newer than the derivative. In that case,
// the derivatives should be discarded
func (s Storage) IsStale(src, drv domain.WebResource) (bool, error) {
	return true, nil
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
