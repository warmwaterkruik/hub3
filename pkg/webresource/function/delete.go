package function

import (
	"context"

	"cloud.google.com/go/storage"
)

func CloudDelete(ctx context.Context, e GCSEvent) error {
	sk, err := NewStoreKey(e.Name)
	if err != nil {
		return err
	}

	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	store := NewCloudStore(ctx, e.Bucket, client)

	return sk.DeleteDerivatives(e.Bucket, store)
}
