package posthook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	c "github.com/delving/rapid-saas/config"
	"github.com/delving/rapid-saas/hub3/fragments"
	"github.com/gammazero/workerpool"
	r "github.com/kiivihal/rdf2go"
	ld "github.com/linkeddata/gojsonld"
	"github.com/parnurzeal/gorequest"
)

// PostHookJob  holds the info for building a crea
type PostHookJob struct {
	Graph   string
	Spec    string
	Deleted bool
	Subject string
	jsonld  []map[string]interface{}
}

type PostHookCounter struct {
	ToIndex           int  `json:"toIndex"`
	ToDelete          int  `json:"toDelete"`
	InError           int  `json:"inError"`
	LifeTimeQueued    int  `json:"lifeTimeQueued"`
	LifeTimeProcessed int  `json:"lifeTimeProcessed"`
	IsActive          bool `json:"isActive"`
}

type PostHookGauge struct {
	Created        time.Time                   `json:"created"`
	QueueSize      int                         `json:"queueSize"`
	ActiveDatasets int                         `json:"activeDatasets"`
	Counters       map[string]*PostHookCounter `json:"counters"`
	sync.Mutex
}

func (phg *PostHookGauge) SetActive(counter *PostHookCounter) {
	if counter.LifeTimeProcessed != counter.LifeTimeQueued {
		if counter.IsActive {
			return
		}
		counter.IsActive = true
		return
	}
	if counter.IsActive {
		counter.IsActive = false
		return
	}
	return
}

func (phg *PostHookGauge) Done(ph *PostHookJob) error {
	counter, ok := phg.Counters[ph.Spec]
	if !ok {
		counter = &PostHookCounter{}
		phg.Counters[ph.Spec] = counter
	}
	phg.Lock()
	defer phg.Unlock()
	phg.QueueSize--
	switch ph.Deleted {
	case true:
		counter.ToDelete--
	default:
		counter.ToIndex--
	}
	counter.LifeTimeProcessed++
	phg.SetActive(counter)
	return nil
}

func (phg *PostHookGauge) Error(ph *PostHookJob) error {
	counter, ok := phg.Counters[ph.Spec]
	if !ok {
		counter = &PostHookCounter{}
		phg.Counters[ph.Spec] = counter
	}
	phg.Lock()
	defer phg.Unlock()
	switch ph.Deleted {
	case true:
		counter.ToDelete--
	default:
		counter.ToIndex--
	}
	counter.InError++
	counter.LifeTimeProcessed++
	phg.QueueSize--
	phg.SetActive(counter)
	return nil
}

func (phg *PostHookGauge) Queue(ph *PostHookJob) error {
	//log.Println("queing posthook")
	counter, ok := phg.Counters[ph.Spec]
	if !ok {
		counter = &PostHookCounter{}
		phg.Counters[ph.Spec] = counter
	}
	phg.Lock()
	defer phg.Unlock()
	counter.LifeTimeQueued++
	phg.QueueSize++
	phg.SetActive(counter)

	if ph.Deleted {
		counter.ToDelete++
		return nil
	}

	counter.ToIndex++
	return nil
}

var gauge PostHookGauge
var wp *workerpool.WorkerPool

func init() {
	wp = workerpool.New(100)
	gauge = PostHookGauge{
		Created:  time.Now(),
		Counters: make(map[string]*PostHookCounter),
	}
}

// NewPostHookJob creates a new PostHookJob and populates the rdf2go Graph
func NewPostHookJob(g, spec string, delete bool, subject, hubID string) (*PostHookJob, error) {
	//add foaf about

	ph := &PostHookJob{
		Graph:   g,
		Spec:    spec,
		Deleted: delete,
		Subject: subject,
	}
	if !delete {
		// setup the cleanup
		err := ph.parseJsonLD()
		if err != nil {
			return nil, err
		}

		ph.addNarthexDefaults(hubID)
		ph.cleanPostHookGraph()
		//log.Printf("%#v", ph.jsonld)

		err = ph.updateJsonLD()
		if err != nil {
			return nil, err
		}
	}
	return ph, nil
}

func (ph *PostHookJob) parseJsonLD() error {
	var jsonld []map[string]interface{}
	err := json.Unmarshal([]byte(ph.Graph), &jsonld)
	if err != nil {
		return err
	}
	ph.jsonld = jsonld
	return nil
}

func (ph *PostHookJob) updateJsonLD() error {
	b, err := json.Marshal(ph.jsonld)
	if err != nil {
		return err
	}
	ph.Graph = string(b)
	return nil
}

func (ph *PostHookJob) addNarthexDefaults(hubID string) {
	//log.Printf("adding defaults for %s", ph.Subject)
	parts := strings.Split(hubID, "_")
	localID := parts[2]
	subject := ph.Subject + "/about"
	defaults := make(map[string]interface{})
	defaults["@id"] = subject
	defaults["@type"] = []string{"http://xmlns.com/foaf/0.1/Document"}
	defaults["http://schemas.delving.eu/narthex/terms/localId"] = []string{localID}
	defaults["http://schemas.delving.eu/narthex/terms/hubID"] = []string{hubID}
	defaults["http://schemas.delving.eu/narthex/terms/spec"] = []string{ph.Spec}

	ph.jsonld = append(ph.jsonld, defaults)
}

// Valid determines if the posthok is valid to apply.
func (ph PostHookJob) Valid() bool {
	return ProcessSpec(ph.Spec)
}

// ProcessSpec determines if a PostHookJob should be applied for a specific spec
func ProcessSpec(spec string) bool {
	if c.Config.PostHook.URL == "" {
		return false
	}
	for _, e := range c.Config.PostHook.ExcludeSpec {
		if e == spec {
			return false
		}
	}
	return true
}

func Submit(ph *PostHookJob) {
	gauge.Queue(ph)
	wp.Submit(func() { ApplyPostHookJob(ph) })
	return
}

// ApplyPostHookJob applies the PostHookJob to all the configured URLs
func ApplyPostHookJob(ph *PostHookJob) {
	u := strings.TrimSuffix(c.Config.PostHook.URL, "/")
	url := fmt.Sprintf("%s/api/erfgoedbrabant/brabantcloud", u)
	err := ph.Post(url)
	if err != nil {
		gauge.Error(ph)
		log.Println(err)
		log.Printf("Unable to send %s to %s", ph.Subject, u)
		//} else {
		//log.Printf("stored: %s", ph.Subject)
		return
	}
	gauge.Done(ph)
	return
}

// Post sends json-ld to the specified endpointt
func (ph PostHookJob) Post(url string) error {

	request := gorequest.New()
	if ph.Deleted {
		log.Printf("Deleting via posthook: %s", ph.Subject)
		deleteURL := fmt.Sprintf("%s/delete", url)
		req := request.Delete(deleteURL).
			Query(fmt.Sprintf("id=%s&api_key=%s", ph.Subject, c.Config.PostHook.APIKey)).
			Retry(3, 5*time.Second, http.StatusBadRequest, http.StatusInternalServerError, http.StatusRequestTimeout)
		//log.Printf("%v", req)
		rsp, body, errs := req.End()
		if errs != nil {
			switch rsp.StatusCode {
			case http.StatusOK, http.StatusCreated:
			default:
				log.Printf("post-response: %#v -> %#v\n %#v", rsp, body, errs)
				log.Printf("Unable to delete: %#v", errs)
				return fmt.Errorf("Unable to save %s to endpoint %s", ph.Subject, url)
			}
		}
		//log.Printf("Deleted %s\n", ph.Subject)
		return nil
	}
	json, err := ph.String()

	if err != nil {
		return err
	}

	rsp, body, errs := request.Post(url).
		Set("Content-Type", "application/json-ld; charset=utf-8").
		Query(fmt.Sprintf("api_key=%s", c.Config.PostHook.APIKey)).
		Type("text").
		Send(json).
		Retry(3, 5*time.Second, http.StatusBadRequest, http.StatusInternalServerError, http.StatusRequestTimeout).
		End()
	//fmt.Printf("jsonld: %s\n", json)
	if errs != nil || rsp.StatusCode != http.StatusOK {
		log.Printf("post-response: %#v -> %#v\n %#v", rsp, body, errs)
		log.Printf("Unable to store: %#v\n", errs)
		log.Printf("JSON-LD: %s\n", json)
		return fmt.Errorf("Unable to save %s to endpoint %s", ph.Subject, url)
	}
	//log.Printf("Stored %s\n", ph.Subject)
	return nil
}

var (
	ns = struct {
		rdf, rdfs, acl, cert, foaf, stat, dc, dcterms, nave, rdagr2, edm ld.NS
	}{
		rdf:     ld.NewNS("http://www.w3.org/1999/02/22-rdf-syntax-ns#"),
		rdfs:    ld.NewNS("http://www.w3.org/2000/01/rdf-schema#"),
		acl:     ld.NewNS("http://www.w3.org/ns/auth/acl#"),
		cert:    ld.NewNS("http://www.w3.org/ns/auth/cert#"),
		foaf:    ld.NewNS("http://xmlns.com/foaf/0.1/"),
		stat:    ld.NewNS("http://www.w3.org/ns/posix/stat#"),
		dc:      ld.NewNS("http://purl.org/dc/elements/1.1/"),
		dcterms: ld.NewNS("http://purl.org/dc/terms/"),
		nave:    ld.NewNS("http://schemas.delving.eu/nave/terms/"),
		rdagr2:  ld.NewNS("http://rdvocab.info/ElementsGr2/"),
		edm:     ld.NewNS("http://www.europeana.eu/schemas/edm/"),
	}
)

var dateFields = []ld.Term{
	ns.dcterms.Get("created"),
	ns.dcterms.Get("issued"),
	ns.nave.Get("creatorBirthYear"),
	ns.nave.Get("creatorDeathYear"),
	ns.nave.Get("date"),
	ns.dc.Get("date"),
	ns.nave.Get("dateOfBurial"),
	ns.nave.Get("dateOfDeath"),
	ns.nave.Get("productionEnd"),
	ns.nave.Get("productionStart"),
	ns.nave.Get("productionPeriod"),
	ns.rdagr2.Get("dateOfBirth"),
	ns.rdagr2.Get("dateOfDeath"),
}

func cleanDates(g *fragments.SortedGraph, t *r.Triple) bool {
	for _, date := range dateFields {
		if t.Predicate.RawValue() == date.RawValue() {
			newTriple := r.NewTriple(
				t.Subject,
				r.NewResource(fmt.Sprintf("%sRaw", t.Predicate.RawValue())),
				t.Object,
			)
			g.Add(newTriple)
			return true
		}
	}
	return false
}

func cleanEbuCore(g *fragments.SortedGraph, t *r.Triple) bool {
	uri := t.Predicate.RawValue()
	if strings.HasPrefix(uri, "urn:ebu:metadata-schema:ebuCore_2014") {
		uri := strings.TrimLeft(uri, "urn:ebu:metadata-schema:ebuCore_2014")
		uri = strings.TrimLeft(uri, "/")
		uri = fmt.Sprintf("http://www.ebu.ch/metadata/ontologies/ebucore/ebucore#%s", uri)
		g.AddTriple(
			t.Subject,
			r.NewResource(uri),
			t.Object,
		)
		return true
	}
	return false
}

func cleanDateURI(uri string) string {
	return fmt.Sprintf("%sRaw", uri)
}

// cleanPostHookGraph applies post hook clean actions to the graph
func (ph *PostHookJob) cleanPostHookGraph() {
	cleanMap := []map[string]interface{}{}
	for _, rsc := range ph.jsonld {
		cleanEntry := make(map[string]interface{})
		for uri, v := range rsc {
			if strings.HasPrefix(uri, "urn:ebu:metadata-schema:ebuCore_2014") {
				uri = strings.TrimLeft(uri, "urn:ebu:metadata-schema:ebuCore_2014")
				uri = strings.TrimLeft(uri, "/")
				uri = fmt.Sprintf("http://www.ebu.ch/metadata/ontologies/ebucore/ebucore#%s", uri)
			}
			switch uri {
			case ns.dcterms.Get("created").RawValue():
				uri = cleanDateURI(uri)
			case ns.dcterms.Get("issued").RawValue():
				uri = cleanDateURI(uri)
			case ns.nave.Get("creatorBirthYear").RawValue():
				uri = cleanDateURI(uri)
			case ns.nave.Get("creatorDeathYear").RawValue():
				uri = cleanDateURI(uri)
			case ns.nave.Get("date").RawValue():
				uri = cleanDateURI(uri)
			case ns.dc.Get("date").RawValue():
				uri = cleanDateURI(uri)
			case ns.nave.Get("dateOfBurial").RawValue():
				uri = cleanDateURI(uri)
			case ns.nave.Get("dateOfDeath").RawValue():
				uri = cleanDateURI(uri)
			case ns.nave.Get("productionEnd").RawValue():
				uri = cleanDateURI(uri)
			case ns.nave.Get("productionStart").RawValue():
				uri = cleanDateURI(uri)
			case ns.nave.Get("productionPeriod").RawValue():
				uri = cleanDateURI(uri)
			case ns.rdagr2.Get("dateOfBirth").RawValue():
				uri = cleanDateURI(uri)
			case ns.rdagr2.Get("dateOfDeath").RawValue():
				uri = cleanDateURI(uri)

			}
			cleanEntry[uri] = v

		}
		cleanMap = append(cleanMap, cleanEntry)

	}
	ph.jsonld = cleanMap
}

// Bytes returns the PostHookJob as an JSON-LD bytes.Buffer
func (ph PostHookJob) Bytes() (bytes.Buffer, error) {
	var b bytes.Buffer
	b.WriteString(ph.Graph)
	return b, nil
}

// Bytes returns the PostHookJob as an JSON-LD string
func (ph PostHookJob) String() (string, error) {

	return ph.Graph, nil
}
