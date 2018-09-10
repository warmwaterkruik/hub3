package experimental

import (
	"log"
	"net/http"

	"github.com/delving/rapid-saas/hub3/fragments"
	"github.com/delving/rapid-saas/hub3/models"
	"github.com/go-chi/render"
)

func csvDelete(w http.ResponseWriter, r *http.Request) {
	conv := fragments.NewCSVConvertor()
	conv.DefaultSpec = r.FormValue("defaultSpec")

	if conv.DefaultSpec == "" {
		render.Status(r, http.StatusBadRequest)
		render.PlainText(w, r, "defaultSpec is a required field")
		return
	}

	ds, _, err := models.GetOrCreateDataSet(conv.DefaultSpec)
	if err != nil {
		log.Printf("Unable to get DataSet for %s\n", conv.DefaultSpec)
		render.PlainText(w, r, err.Error())
		return
	}
	_, err = ds.DropRecords(r.Context(), nil)
	if err != nil {
		log.Printf("Unable to delete all fragments for %s: %s", conv.DefaultSpec, err.Error())
		render.Status(r, http.StatusBadRequest)
		return
	}

	render.Status(r, http.StatusNoContent)
	return
}

func csvUpload(w http.ResponseWriter, r *http.Request) {
	in, _, err := r.FormFile("csv")
	if err != nil {
		render.PlainText(w, r, err.Error())
		return
	}

	conv := fragments.NewCSVConvertor()
	conv.InputFile = in
	conv.SubjectColumn = r.FormValue("subjectColumn")
	conv.SubjectClass = r.FormValue("subjectClass")
	conv.SubjectURIBase = r.FormValue("subjectURIBase")
	conv.Separator = r.FormValue("separator")
	conv.PredicateURIBase = r.FormValue("predicateURIBase")
	conv.SubjectColumn = r.FormValue("subjectColumn")
	conv.ObjectResourceColumns = []string{r.FormValue("objectResourceColumns")}
	conv.ObjectURIFormat = r.FormValue("objectURIFormat")
	conv.DefaultSpec = r.FormValue("defaultSpec")
	conv.ThumbnailURIBase = r.FormValue("thumbnailURIBase")
	conv.ThumbnailColumn = r.FormValue("thumbnailColumn")
	conv.ManifestColumn = r.FormValue("manifestColumn")
	conv.ManifestURIBase = r.FormValue("manifestURIBase")
	conv.ManifestLocale = r.FormValue("manifestLocale")

	ds, created, err := models.GetOrCreateDataSet(conv.DefaultSpec)
	if err != nil {
		log.Printf("Unable to get DataSet for %s\n", conv.DefaultSpec)
		render.PlainText(w, r, err.Error())
		return
	}

	// todo replace with proper implemtation later
	if created {
		err = fragments.SaveDataSet(conv.DefaultSpec, bp)
		if err != nil {
			log.Printf("Unable to Save DataSet Fragment for %s\n", conv.DefaultSpec)
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

	triplesCreated, rowsSeen, err := conv.IndexFragments(bp, ds.Revision)
	conv.RowsProcessed = rowsSeen
	conv.TriplesCreated = triplesCreated
	log.Printf("Processed %d csv rows\n", rowsSeen)
	if err != nil {
		render.PlainText(w, r, err.Error())
		return
	}

	_, err = ds.DropOrphans(r.Context(), bp, wp)
	if err != nil {
		render.PlainText(w, r, err.Error())
		return
	}

	render.Status(r, http.StatusCreated)
	//render.PlainText(w, r, "ok")
	render.JSON(w, r, conv)
	return
}
