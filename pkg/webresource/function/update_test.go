package function

import (
	"context"
	"testing"
)

func TestCloudUpdate(t *testing.T) {
	t.SkipNow()
	ctx := context.Background()

	e := GCSEvent{
		Bucket: "websrc-hub3-saas-nl-prod",
		Name:   "test-org/test-spec/123.jpg",
	}

	err := CloudUpdate(ctx, e)
	if err != nil {
		t.Fatal(err)
	}
}
