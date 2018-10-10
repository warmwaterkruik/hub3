package rest

import (
	"net/http"

	"github.com/delving/rapid-saas/pkg/engine"
	"github.com/go-chi/chi"
)

// Handler returns a REST router
func Handler(service engine.Service) http.Handler {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("you are rocking rapid"))
	})

	// Setup the static directory
	//FileServer(r, "/static", getAbsolutePathToFileDir(s.))
	return r
}
