package experimental

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/deiu/rdf2go"
	"github.com/delving/rapid-saas/hub3/fragments"
	"github.com/delving/rapid-saas/hub3/models"
	"github.com/go-chi/render"
)

func skosUpload(w http.ResponseWriter, r *http.Request) {

	// create byte.buffer from the input file
	// get dataset param
	// get dataset object
	// get subjectClass
	// create map subject map[string]bool
	// matcher string for rdf:Type
	// create resourceMap
	// use n-triple / turtle parse to build line by line, see rdf2go libraries for Graph
	// addTriple per line
	// gather subject per type
	// check subject map
	// next

	// alternative approach
	// store all fragments
	// get scanner for spec and rdfType to get subjects back
	// make nested call for elasticsearch: get all objects, do mget on fragments,
	// parse into resource map
	// do next level mget on resource objects
	// parse into resource map
	// add to elastic bulk processor
	in, _, err := r.FormFile("skos")
	if err != nil {
		render.PlainText(w, r, err.Error())
		return
	}
	//io.Copy(w, in)
	var buff bytes.Buffer
	fileSize, err := buff.ReadFrom(in)
	//fmt.Println(fileSize) // this will return you a file size.
	//if err != nil {
	//render.PlainText(w, r, err.Error())
	//return
	//}
	render.PlainText(w, r, fmt.Sprintf("The file is %d bytes long", fileSize))

	jsonld := []map[string]interface{}{}
	err = json.Unmarshal(buff.Bytes(), &jsonld)
	if err != nil {
		render.PlainText(w, r, err.Error())
		return
	}

	log.Printf("found %#v resources", jsonld[0])
	log.Printf("found %d resources", len(jsonld))

	defer in.Close()
	//g := rdf2go.NewGraph("")
	//err = g.Parse(in, "application/ld+json")
	//if err != nil {
	//render.PlainText(w, r, err.Error())
	//return
	//}

	//render.PlainText(w, r, fmt.Sprintf("processed triples: %d", g.Len()))
	return
}

func skosSync(w http.ResponseWriter, r *http.Request) {
	targetURL := r.URL.Query().Get("uri")
	spec := r.URL.Query().Get("spec")

	ds, created, err := models.GetOrCreateDataSet(spec)
	if err != nil {
		log.Printf("Unable to get DataSet for %s\n", spec)
		render.PlainText(w, r, err.Error())
		return
	}
	if created {
		err = fragments.SaveDataSet(spec, bp)
		if err != nil {
			log.Printf("Unable to Save DataSet Fragment for %s\n", spec)
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

	g := rdf2go.NewGraph("")
	err = g.LoadURI(targetURL)
	if err != nil {
		log.Printf("Unable to get skos for %s\n", targetURL)
		render.PlainText(w, r, err.Error())
		return
	}

	render.PlainText(w, r, fmt.Sprintf("processed triples: %d", g.Len()))
	return
}
