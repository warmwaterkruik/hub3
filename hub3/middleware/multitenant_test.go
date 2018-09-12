package middleware_test

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/delving/rapid-saas/hub3/middleware"
	"github.com/go-chi/chi"
)

func TestMultiTenant(t *testing.T) {
	r := chi.NewRouter()

	// This middleware must be mounted at the top level of the router, not at the end-handler
	// because then it'll be too late and will end up in a 404
	r.Use(middleware.MultiTenant)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("nothing here"))
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("root"))
	})

	//r.Route("/{orgID}/", func(r chi.Router) {
	r.Get("/{orgID}/123", func(w http.ResponseWriter, r *http.Request) {
		orgID := chi.URLParam(r, "orgID")
		w.Write([]byte(orgID))
	})
	//})

	ts := httptest.NewServer(r)
	defer ts.Close()

	//_, orgID, resp := testRequest(t, ts, "GET", "/test-org/123", nil)
	//t.Logf("orgID = %s (%s)", resp)
	//fmt.Printf("orgID = %s (%s)\n", orgID, resp)

	//if _, resp := testRequest(t, ts, "GET", "/", nil); resp != "root" {
	//t.Fatalf(resp)
	//}
	//if _, resp := testRequest(t, ts, "GET", "//", nil); resp != "root" {
	//t.Fatalf(resp)
	//}
	//if _, orgID, resp := testRequest(t, ts, "GET", "/test-org/123", nil); orgID != "test-org" {
	//t.Fatalf("orgID: %s, response: %s\n", orgID, resp)
	//}
	//if _, resp := testRequest(t, ts, "GET", "/accounts/admin/", nil); resp != "admin" {
	//t.Fatalf(resp)
	//}
	//if _, resp := testRequest(t, ts, "GET", "/nothing-here", nil); resp != "nothing here" {
	//t.Fatalf(resp)
	//}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (
	*http.Response, string, string) {

	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return nil, "", ""
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return nil, "", ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return nil, "", ""
	}
	defer resp.Body.Close()

	return resp, middleware.GetOrgID(req.Context()), string(respBody)
}

func TestGetOrgID(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"simple get", args{ctx: context.WithValue(context.TODO(), middleware.OrgIDKey, "test-org")}, "test-org"},
		{"default orgID", args{ctx: context.TODO()}, "rapid"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := middleware.GetOrgID(tt.args.ctx); got != tt.want {
				t.Errorf("GetOrgID() = %v, want %v", got, tt.want)
			}
		})
	}
}
