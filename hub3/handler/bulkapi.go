// Copyright © 2017 Delving B.V. <info@delving.eu>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	c "github.com/delving/rapid-saas/config"
	"github.com/delving/rapid-saas/hub3/fragments"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// LODResource is the router struct for LOD
type BulkAPIResource struct{}

// Routes returns the chi.Router
func (rs BulkAPIResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/rdf", bulkAPI)
	r.Get("/sync", bulkSyncList)
	r.Post("/sync", bulkSyncStart)
	r.Get("/sync/{id}", bulkSyncProgress)
	r.Delete("/sync/{id}", bulkSyncCancel)

	// backwards compatibility
	r.Post("/bulk", bulkAPI)
	r.Post("/fuzzed", generateFuzzed)
	return r
}

func bulkSyncStart(w http.ResponseWriter, r *http.Request) {
	//host := r.URL.Query().Get("host")
	//index := r.URL.Query().Get("index")
	http.Error(w, "not implemented", http.StatusBadRequest)
	return

}

func bulkSyncList(w http.ResponseWriter, r *http.Request) {
	//host := r.URL.Query().Get("host")
	//index := r.URL.Query().Get("index")
	http.Error(w, "not implemented", http.StatusBadRequest)
	return

}

func bulkSyncProgress(w http.ResponseWriter, r *http.Request) {

	http.Error(w, "not implemented", http.StatusBadRequest)
	return
}

func bulkSyncCancel(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusBadRequest)
	return

}

// bulkApi receives bulkActions in JSON form (1 per line) and processes them in
// ingestion pipeline.
func bulkAPI(w http.ResponseWriter, r *http.Request) {
	// todo insert bp and wp in function below
	//response, err := ReadActions(r.Contetx(), r.Body, bp, wp)
	response, err := ReadActions(r.Context(), r.Body, nil, nil)
	if err != nil {
		log.Println("Unable to read actions")
		errR := ErrRender(err)
		// todo fix errr renderer for better narthex consumption.
		_ = errR.Render(w, r)
		render.Render(w, r, errR)
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, response)
	return
}

func generateFuzzed(w http.ResponseWriter, r *http.Request) {
	in, _, err := r.FormFile("file")
	if err != nil {
		render.PlainText(w, r, err.Error())
		return
	}
	spec := r.FormValue("spec")
	number := r.FormValue("number")
	baseURI := r.FormValue("baseURI")
	subjectType := r.FormValue("rootType")
	n, err := strconv.Atoi(number)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	recDef, err := fragments.NewRecDef(in)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fz, err := fragments.NewFuzzer(recDef)
	fz.BaseURL = baseURI
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	records, err := fz.CreateRecords(n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	typeLabel, err := c.Config.NameSpaceMap.GetSearchLabel(subjectType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	actions := []string{}
	for idx, rec := range records {
		hubID := fmt.Sprintf("%s_%s_%d", c.Config.OrgID, spec, idx)
		action := &BulkAction{
			HubID:         hubID,
			OrgID:         c.Config.OrgID,
			LocalID:       fmt.Sprintf("%d", idx),
			Spec:          spec,
			NamedGraphURI: fmt.Sprintf("%s/graph", fz.NewURI(typeLabel, idx)),
			Action:        "index",
			GraphMimeType: "application/ld+json",
			SubjectType:   subjectType,
			RecordType:    "mdr",
			Graph:         rec,
		}
		bytes, err := json.Marshal(action)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			log.Printf("Unable to create Bulkactions: %s\n", err.Error())
			render.PlainText(w, r, err.Error())
			return
		}
		actions = append(actions, string(bytes))
	}
	render.PlainText(w, r, strings.Join(actions, "\n"))
	//w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	return
}
