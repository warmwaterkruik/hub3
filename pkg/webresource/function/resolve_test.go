package function

import (
	fmt "fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestCloudResolve(t *testing.T) {
	bucketName := "webdrv-test-bucket"
	os.Setenv("DERIVATIVE_BUCKET_NAME", bucketName)
	r, err := http.NewRequest("GET", "/thumbnail/brabantcloud/helmond-objecten/14495/220", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(CloudResolve)
	handler.ServeHTTP(w, r)

	resp := w.Result()

	if resp.StatusCode != http.StatusFound {
		t.Errorf("wrong status code: got %v want %v", resp.StatusCode, http.StatusOK)
	}

	loc, _ := resp.Location()
	if loc.Hostname() != "storage.googleapis.com" {
		t.Error("expected default media host")
	}

	if !strings.HasPrefix(loc.Path, fmt.Sprintf("/%s", bucketName)) {
		t.Errorf("expected default path to start with the bucket name")
	}
}
