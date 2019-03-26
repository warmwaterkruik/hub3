package function

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
)

type StoreKey struct {
	Path      string
	OrgID     string
	Spec      string
	LocalPath string
	Operation string
}

func NewStoreKey(path string) (*StoreKey, error) {
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return nil, fmt.Errorf("cannot parse storekey from %s", path)
	}
	pathSplitter := 2
	if parts[2] == "source" {
		pathSplitter = 3
	}
	localPath := strings.Join(parts[pathSplitter:], "/")
	return &StoreKey{
		Path:      path,
		OrgID:     parts[0],
		Spec:      parts[1],
		LocalPath: localPath,
		Operation: "",
	}, nil
}

func (sk StoreKey) Key() string {
	path := strings.ToLower(sk.LocalPath)
	ext := filepath.Ext(path)
	path = strings.TrimSuffix(path, ext)
	return replacer.Replace(path)
}

func (sk StoreKey) hasSuffix() bool {
	return filepath.Ext(sk.Path) != ""
}

func (sk StoreKey) Get(bucket string, store *CloudStore) (*WebResource, error) {
	b, err := store.read(bucket, sk.PBPath())
	if err != nil {
		return nil, err
	}

	wr := &WebResource{}
	err = proto.Unmarshal(b, wr)
	if err != nil {
		return nil, err
	}
	return wr, nil
}

func (sk StoreKey) ObjectExists(bucket string, store *CloudStore) (string, bool) {
	objectPath := sk.KeyPrefix()
	ext := filepath.Ext(objectPath)
	switch ext {
	case "":
		objectPath = objectPath + ".jpg"
	default:
		objectPath = objectPath + ext
	}
	_, err := store.attrs(bucket, objectPath)
	if err != nil {
		return "", false
	}

	return objectPath, true
}

func (sk StoreKey) KeyPrefix() string {
	return fmt.Sprintf("%s/%s/%s", sk.OrgID, sk.Spec, sk.Key())
}

func (sk StoreKey) GetPublicURL(bucket, requestDomain, path string) string {
	domain := fmt.Sprintf("storage.googleapis.com/%s", bucket)
	if strings.HasPrefix(requestDomain, "media") {
		domain = strings.Replace(requestDomain, "media.", "static.", 1)
	}

	return fmt.Sprintf("https://%s/%s", domain, path)
}

func (sk StoreKey) PBPath() string {
	return fmt.Sprintf("%s.pb", sk.KeyPrefix())
}

func (sk StoreKey) Update(wr *WebResource, bucket string, store *CloudStore) error {
	out, err := proto.Marshal(wr)
	if err != nil {
		return errors.Wrapf(err, "cannot marshal Webresource: %#v", wr)
	}
	derivativeBucket := store.getDerivativeBucket(bucket)
	return store.Write(derivativeBucket, sk.PBPath(), out, "application/octet-stream")
}

func (sk StoreKey) DeleteDerivatives(bucket string, store *CloudStore) error {
	derivativeBucket := store.getDerivativeBucket(bucket)
	return store.deleteByPrefix(derivativeBucket, sk.KeyPrefix(), "")
}

func (sk StoreKey) CreateDefaultDerivatives(wr *WebResource, bucket string, store *CloudStore) error {
	return nil
}

func (sk StoreKey) createWebResource(hash string) (*WebResource, error) {
	wr := &WebResource{
		OrgID:         sk.OrgID,
		Spec:          sk.Spec,
		RelativePath:  sk.LocalPath,
		NormalisedKey: sk.Key(),
		MD5Hash:       hash,
		SourcePath:    sk.Path,
	}
	return wr, nil
}
