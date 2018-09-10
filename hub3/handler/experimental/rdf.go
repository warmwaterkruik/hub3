package experimental

import (
	"fmt"
	"log"
	"net/http"

	"github.com/delving/rapid-saas/config"
	"github.com/delving/rapid-saas/hub3/fragments"
	"github.com/delving/rapid-saas/hub3/models"
	"github.com/go-chi/render"
	"github.com/gorilla/schema"
)

type rdfUploadForm struct {
	Spec          string `json:"spec"`
	RecordType    string `json:"recordType"`
	TypePredicate string `json:"typePredicate"`
	IDSplitter    string `json:"idSplitter"`
}

func (ruf *rdfUploadForm) isValid() error {
	if ruf.Spec == "" {
		return fmt.Errorf("spec param is required")
	}
	if ruf.RecordType == "" {
		return fmt.Errorf("recordType param is required")
	}
	if ruf.TypePredicate == "" {
		return fmt.Errorf("typePredicate param is required")
	}
	if ruf.IDSplitter == "" {
		return fmt.Errorf("idSplitter param is required")
	}
	return nil
}

var decoder = schema.NewDecoder()

func rdfUpload(w http.ResponseWriter, r *http.Request) {
	in, _, err := r.FormFile("turtle")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var form rdfUploadForm
	err = decoder.Decode(&form, r.PostForm)
	if err != nil {
		log.Printf("Unable to decode form %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = form.isValid()
	if err != nil {
		log.Printf("form is not valid; %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// todo handle when no form.Spec is given
	ds, created, err := models.GetOrCreateDataSet(form.Spec)
	if err != nil {
		log.Printf("Unable to get DataSet for %s\n", form.Spec)
		render.PlainText(w, r, err.Error())
		return
	}
	if created {
		err = fragments.SaveDataSet(form.Spec, bp)
		if err != nil {
			log.Printf("Unable to Save DataSet Fragment for %s\n", form.Spec)
			if err != nil {
				render.PlainText(w, r, err.Error())
				return
			}
		}
	}

	ds, err = ds.IncrementRevision()
	if err != nil {
		render.PlainText(w, r, err.Error())
		return
	}

	upl := fragments.NewRDFUploader(
		config.Config.OrgID,
		form.Spec,
		form.RecordType,
		form.TypePredicate,
		form.IDSplitter,
		ds.Revision,
	)

	go func() {
		defer in.Close()
		log.Print("Start creating resource map")
		_, err := upl.Parse(in)
		if err != nil {
			log.Printf("Can't read turtle file: %v", err)
			return
		}
		log.Printf("Start saving fragments.")
		//processed, err := upl.IndexFragments(bp)
		//if err != nil {
		//log.Printf("Can't save fragments: %v", err)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		//return
		//}
		//log.Printf("Saved %d fragments for %s", processed, upl.Spec)
		processed, err := upl.SaveFragmentGraphs(bp)
		if err != nil {
			log.Printf("Can't save records: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Saved %d records for %s", processed, upl.Spec)
		ds.DropOrphans(r.Context(), bp, nil)
	}()

	render.Status(r, http.StatusCreated)
	render.PlainText(w, r, "ok")
	return
}
