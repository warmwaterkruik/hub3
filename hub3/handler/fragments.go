package handler

import (
	"bytes"
	"fmt"
	log "log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"

	"github.com/delving/rapid-saas/hub3/fragments"
	"github.com/delving/rapid-saas/hub3/index"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

// ListFragments returns a list of matching fragments
// See for more info: http://linkeddatafragments.org/
func ListFragments(w http.ResponseWriter, r *http.Request) {
	fr := fragments.NewFragmentRequest()
	spec := chi.URLParam(r, "spec")
	if spec != "" {
		fr.Spec = spec
	}
	err := fr.ParseQueryString(r.URL.Query())
	if err != nil {
		log.Printf("Unable to list fragments because of: %s", err)
		render.JSON(w, r, APIErrorMessage{
			HTTPStatus: http.StatusBadRequest,
			Message:    fmt.Sprint("Unable to list fragments was not found"),
			Error:      err,
		})
		return
	}
	frags, totalFrags, err := fr.Find(r.Context(), index.ESClient())
	if err != nil || len(frags) == 0 {
		log.Printf("Unable to list fragments because of: %s", err)
		render.JSON(w, r, APIErrorMessage{
			HTTPStatus: http.StatusNotFound,
			Message:    fmt.Sprint("No fragmenst for query were found."),
			Error:      err,
		})
		return
	}
	switch fr.Echo {
	case "raw":
		render.JSON(w, r, frags)
		return
	case "es":
		src, err := fr.BuildQuery().Source()
		if err != nil {
			msg := "Unable to get the query source"
			log.Printf(msg)
			render.JSON(w, r, APIErrorMessage{
				HTTPStatus: http.StatusBadRequest,
				Message:    fmt.Sprint(msg),
				Error:      err,
			})
			return
		}
		render.JSON(w, r, src)
		return
	case "searchResponse":
		res, err := fr.Do(r.Context(), index.ESClient())
		if err != nil {
			msg := fmt.Sprintf("Unable to dump request: %s", err)
			log.Print(msg)
			render.JSON(w, r, APIErrorMessage{
				HTTPStatus: http.StatusBadRequest,
				Message:    fmt.Sprint(msg),
				Error:      err,
			})
			return
		}
		render.JSON(w, r, res)
		return
	case "request":
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			msg := fmt.Sprintf("Unable to dump request: %s", err)
			log.Print(msg)
			render.JSON(w, r, APIErrorMessage{
				HTTPStatus: http.StatusBadRequest,
				Message:    fmt.Sprint(msg),
				Error:      err,
			})
			return
		}

		render.PlainText(w, r, string(dump))
		return
	}

	var buffer bytes.Buffer
	for _, frag := range frags {
		buffer.WriteString(fmt.Sprintln(frag.Triple))
	}
	w.Header().Add("FRAG_COUNT", strconv.Itoa(int(totalFrags)))

	// Add hyperMediaControls
	hmd := fragments.NewHyperMediaDataSet(r, totalFrags, fr)
	controls, err := hmd.CreateControls()
	if err != nil {
		msg := fmt.Sprintf("Unable to create media controls: %s", err)
		log.Print(msg)
		render.JSON(w, r, APIErrorMessage{
			HTTPStatus: http.StatusBadRequest,
			Message:    fmt.Sprint(msg),
			Error:      err,
		})
		return
	}

	if strings.Contains(r.Header.Get("Accept"), "n-triples") {
		w.Header().Add("Content-Type", "application/n-triples")
	} else {
		w.Header().Add("Content-Type", "text/plain")
	}

	w.Write(controls)
	w.Write(buffer.Bytes())

	return
}
