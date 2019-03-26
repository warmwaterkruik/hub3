package function

import (
	fmt "fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"cloud.google.com/go/storage"
)

func CloudResolve(w http.ResponseWriter, r *http.Request) {
	ok, w := handlePreFlight(w, r)
	if ok {
		return
	}

	log.Printf("raw path: %s", r.URL.EscapedPath())
	path, err := cleanPath(r.URL.EscapedPath())
	if err != nil {
		log.Printf("path %s not properly formatter: %#v", path, err)
		http.Error(w, "path not properly formatted", http.StatusBadRequest)
		return
	}

	sk, err := NewStoreKey(path)
	if err != nil {
		log.Printf("Unable to create store key for %s: %#v", path, err)
		http.Error(w, "unable to create store key", http.StatusBadRequest)
		return
	}

	bucket := os.Getenv("DERIVATIVE_BUCKET_NAME")
	if bucket == "" {
		log.Fatal("environment variable DERIVATIVE_BUCKET_NAME must be defined")
	}

	ctx := r.Context()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal("can't create storage client")
		return
	}

	store := NewCloudStore(ctx, bucket, client)

	cleanPath, ok := sk.ObjectExists(bucket, store)
	if !ok {
		// TODO: later add create derivative
		log.Printf("object does not exist: %s", path)
		http.Error(w, fmt.Sprintf("object does not exist: %s", path), http.StatusNotFound)
		return
	}

	url := sk.GetPublicURL(bucket, r.Host, cleanPath)

	http.Redirect(w, r, url, http.StatusFound)
	return
}

func cleanPath(path string) (string, error) {
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		return "", fmt.Errorf("malformed request: %s", path)
	}
	switch parts[1] {
	case "thumbnail", "deepzoom":
		parts = parts[2:]
	default:
		parts = parts[1:]
	}

	return strings.Join(parts, "/"), nil
}
