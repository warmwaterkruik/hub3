package function

import (
	"context"
	"io/ioutil"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// CloudStore is a struct that provides functions to manipulate Google Cloud Storage objects
type CloudStore struct {
	client *storage.Client
	ctx    context.Context
	bucket string
}

// NewCloudStore creates a new CloudStore
func NewCloudStore(ctx context.Context, bucket string, client *storage.Client) *CloudStore {
	return &CloudStore{
		client: client,
		ctx:    ctx,
		bucket: bucket,
	}
}

func (cs CloudStore) Write(object string, b []byte, mimeType string) error {

	wc := cs.client.Bucket(cs.bucket).Object(object).NewWriter(cs.ctx)
	if mimeType != "" {
		wc.ContentType = mimeType
	}

	if _, err := wc.Write(b); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}

func (cs CloudStore) list(bucket string) *storage.ObjectIterator {
	return cs.client.Bucket(bucket).Objects(cs.ctx, nil)
}

func (cs CloudStore) listByPrefix(bucket, prefix, delim string) *storage.ObjectIterator {
	// Prefixes and delimiters can be used to emulate directory listings.
	// Prefixes can be used filter objects starting with prefix.
	// The delimiter argument can be used to restrict the results to only the
	// objects in the given "directory". Without the delimeter, the entire  tree
	// under the prefix is returned.
	//
	// For example, given these blobs:
	//   /a/1.txt
	//   /a/b/2.txt
	//
	// If you just specify prefix="a/", you'll get back:
	//   /a/1.txt
	//   /a/b/2.txt
	//
	// However, if you specify prefix="a/" and delim="/", you'll get back:
	//   /a/1.txt
	return cs.client.Bucket(bucket).Objects(cs.ctx, &storage.Query{
		Prefix:    prefix,
		Delimiter: delim,
	})
}

func (cs CloudStore) deleteByPrefix(bucket, prefix, delim string) error {
	it := cs.listByPrefix(bucket, prefix, delim)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		err = cs.delete(bucket, attrs.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cs CloudStore) read(bucket, object string) ([]byte, error) {
	ctx := context.Background()
	// [START download_file]
	rc, err := cs.client.Bucket(bucket).Object(object).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (cs CloudStore) attrs(bucket, object string) (*storage.ObjectAttrs, error) {
	ctx := context.Background()
	// [START get_metadata]
	o := cs.client.Bucket(bucket).Object(object)
	attrs, err := o.Attrs(ctx)
	if err != nil {
		return nil, err
	}

	return attrs, nil
}

func (cs CloudStore) makePublic(bucket, object string) error {
	ctx := context.Background()

	acl := cs.client.Bucket(bucket).Object(object).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return err
	}

	return nil
}

func (cs CloudStore) makePrivate(bucket, object string) error {
	ctx := context.Background()

	acl := cs.client.Bucket(bucket).Object(object).ACL()
	if err := acl.Delete(ctx, storage.AllUsers); err != nil {
		return err
	}

	return nil
}

func (cs CloudStore) move(bucket, object string) error {
	ctx := context.Background()

	dstName := object + "-rename"

	src := cs.client.Bucket(bucket).Object(object)
	dst := cs.client.Bucket(bucket).Object(dstName)

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return err
	}
	if err := src.Delete(ctx); err != nil {
		return err
	}

	return nil
}

func (cs *CloudStore) copyToBucket(dstBucket, srcBucket, srcObject string) error {
	ctx := context.Background()

	dstObject := srcObject + "-copy"
	src := cs.client.Bucket(srcBucket).Object(srcObject)
	dst := cs.client.Bucket(dstBucket).Object(dstObject)

	if _, err := dst.CopierFrom(src).Run(ctx); err != nil {
		return err
	}

	return nil
}

func (cs CloudStore) delete(bucket, object string) error {
	ctx := context.Background()
	// [START delete_file]
	o := cs.client.Bucket(bucket).Object(object)
	if err := o.Delete(ctx); err != nil {
		return err
	}
	// [END delete_file]
	return nil
}

func (cs CloudStore) getDerivativeBucket(sourceBucket string) string {
	return strings.Replace(sourceBucket, "websrc-", "webdrv-", 1)
}

// writeEncryptedObject writes an object encrypted with user-provided AES key to a bucket.
func (cs CloudStore) writeEncryptedObject(bucket, object string, secretKey []byte) error {
	ctx := context.Background()

	// [START storage_upload_encrypted_file]
	obj := cs.client.Bucket(bucket).Object(object)
	// Encrypt the object's contents.
	wc := obj.Key(secretKey).NewWriter(ctx)
	if _, err := wc.Write([]byte("top secret")); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	// [END storage_upload_encrypted_file]
	return nil
}

// writeWithKMSKey writes an object encrypted with KMS-provided key to a bucket.
func (cs CloudStore) writeWithKMSKey(bucket, object string, keyName string) error {
	ctx := context.Background()

	obj := cs.client.Bucket(bucket).Object(object)
	// Encrypt the object's contents
	wc := obj.NewWriter(ctx)
	wc.KMSKeyName = keyName
	if _, err := wc.Write([]byte("top secret")); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}

func (cs CloudStore) readEncryptedObject(bucket, object string, secretKey []byte) ([]byte, error) {
	ctx := context.Background()

	// [START storage_download_encrypted_file]
	obj := cs.client.Bucket(bucket).Object(object)
	rc, err := obj.Key(secretKey).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	// [END storage_download_encrypted_file]
	return data, nil
}

func (cs CloudStore) rotateEncryptionKey(bucket, object string, key, newKey []byte) error {
	ctx := context.Background()
	// [START storage_rotate_encryption_key]
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	obj := client.Bucket(bucket).Object(object)
	// obj is encrypted with key, we are encrypting it with the newKey.
	_, err = obj.Key(newKey).CopierFrom(obj.Key(key)).Run(ctx)
	if err != nil {
		return err
	}
	// [END storage_rotate_encryption_key]
	return nil
}
