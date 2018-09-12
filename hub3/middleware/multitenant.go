package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/delving/rapid-saas/config"
	"github.com/go-chi/chi"
)

type ctxKeyOrgID int
type ctxKeyOrgConfig int

// Keys to used when setting context.Context keys
const (
	OrgIDKey     ctxKeyOrgID     = iota
	OrgConfigKey ctxKeyOrgConfig = iota
)

// MultiTenant sets configuration for the OrgID and the OrgID configuration.
// When the 'OrgID' is missing from the routes the default is inserted from the
// main configuration.
func MultiTenant(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		orgID := chi.URLParam(r, "orgID")
		if orgID == "" {
			orgID = config.Config.OrgID
		}
		fmt.Printf("this is the orgID: %s\n", orgID)
		// TODO set later the orgIDConfig key
		ctx = context.WithValue(ctx, OrgIDKey, orgID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// GetOrgID returns a tentant ID from the given context if one is present.
// Returns the empty string if a tentant ID cannot be found.
func GetOrgID(ctx context.Context) string {
	if ctx == nil {
		fmt.Println("ctx is empty")
		return config.Config.OrgID
	}
	if orgID, ok := ctx.Value(OrgIDKey).(string); ok {
		return orgID
	}
	fmt.Println("can't retrieve key from context")
	return config.Config.OrgID
}
