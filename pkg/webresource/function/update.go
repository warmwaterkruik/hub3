package function

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/pkg/errors"
)

var replacer = strings.NewReplacer(" ", "_")

type GCSEvent struct {
	Bucket                  string `json:"bucket"`
	Name                    string `json:"name"`
	ContentType             string `json:"contentType"`
	Crc32c                  string `json:"crc32c"`
	Etag                    string `json:"etag"`
	Generation              string `json:"generation"`
	Id                      string `json:"id"`
	Kind                    string `json:"kind"`
	Md5Hash                 string `json:"md5Hash"`
	MediaLink               string `json:"mediaLink"`
	Metageneration          string `json:"metageneration"`
	SelfLink                string `json:"selfLink"`
	Size                    string `json:"size"`
	StorageClass            string `json:"storageClass"`
	TimeCreated             string `json:"timeCreated"`
	TimeStorageClassUpdated string `json:"timeStorageClassUpdated"`
	Updated                 string `json:"updated"`
}

func CloudUpdate(ctx context.Context, e GCSEvent) error {
	if strings.HasSuffix(e.Name, ".pb") {
		return nil
	}
	sk, err := NewStoreKey(e.Name)
	if err != nil {
		return err
	}
	wr, err := sk.createWebResource(e.Md5Hash)
	if err != nil {
		return err
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	store := NewCloudStore(ctx, e.Bucket, client)

	err = sk.Update(wr, e.Bucket, store)
	if err != nil {
		return errors.Wrapf(err, "cannot save storekey %s to bucket %s", sk.KeyPrefix(), e.Bucket)
	}

	err = sk.DeleteDerivatives(e.Bucket, store)
	if err != nil {
		return errors.Wrapf(err, "cannot delete derivatives for %s", sk.KeyPrefix())
	}

	err = sk.CreateDefaultDerivatives(wr, e.Bucket, store)
	if err != nil {
		return errors.Wrapf(err, "cannot create default derivatives for %s", sk.KeyPrefix())
	}

	fmt.Printf("%s was uploaded to %s: %#v", e.Name, e.Bucket, e)
	return nil
}
