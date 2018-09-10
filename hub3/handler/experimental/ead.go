package experimental

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/delving/rapid-saas/hub3/experimental/ead"
	"github.com/delving/rapid-saas/hub3/fragments"
	"github.com/delving/rapid-saas/hub3/handler"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func eadUpload(w http.ResponseWriter, r *http.Request) {
	spec := r.FormValue("spec")
	if spec == "" {
		render.PlainText(w, r, "spec param is required")
		render.Status(r, http.StatusBadRequest)
		return
	}

	_, err := ead.ProcessUpload(r, w, spec, bp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}

func treeList(w http.ResponseWriter, r *http.Request) {
	spec := chi.URLParam(r, "spec")
	if spec == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, handler.APIErrorMessage{
			HTTPStatus: http.StatusBadRequest,
			Message:    fmt.Sprintln("spec can't be empty."),
			Error:      nil,
		})
		return
	}
	nodeID := chi.URLParam(r, "nodeID")
	if nodeID != "" {
		id, err := url.QueryUnescape(nodeID)
		if err != nil {
			log.Println("Unable to unescape QueryParameters.")
			render.Status(r, http.StatusBadRequest)
			render.PlainText(w, r, err.Error())
			return
		}
		q := r.URL.Query()
		q.Add("byLeaf", id)
		r.URL.RawQuery = q.Encode()
		//r.URL.Query().Add("byLeaf", nodeID)
	}
	searchRequest, err := fragments.NewSearchRequest(r.URL.Query())
	if err != nil {
		log.Println("Unable to create Search request")
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, err.Error())
		return
	}
	searchRequest.ItemFormat = fragments.ItemFormatType_TREE
	searchRequest.AddQueryFilter(fmt.Sprintf("spec:%s", spec))
	switch searchRequest.Tree {
	case nil:
		searchRequest.Tree = &fragments.TreeQuery{
			Depth: []string{"1", "2"},
			Spec:  spec,
		}
	default:
		searchRequest.Tree.Spec = spec
	}
	handler.ProcessSearchRequest(w, r, searchRequest)
	return
}

func treeStats(w http.ResponseWriter, r *http.Request) {
	spec := chi.URLParam(r, "spec")
	if spec == "" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, handler.APIErrorMessage{
			HTTPStatus: http.StatusBadRequest,
			Message:    fmt.Sprintln("spec can't be empty."),
			Error:      nil,
		})
		return
	}
	stats, err := fragments.CreateTreeStats(r.Context(), spec)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// todo return 404 if stats.Leafs == 0
	if stats.Leafs == 0 {
		render.Status(r, http.StatusNotFound)
		return
	}
	render.JSON(w, r, stats)
	return
}
