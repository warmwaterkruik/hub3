package posthook

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	c "github.com/delving/rapid-saas/config"
	"github.com/delving/rapid-saas/hub3/fragments"
	"github.com/delving/rapid-saas/hub3/index"
	"github.com/go-chi/chi"
	r "github.com/kiivihal/rdf2go"
	"github.com/pkg/errors"
	elastic "gopkg.in/olivere/elastic.v5"
)

type Hit struct {
	Source `json:"_source"`
}

type Source struct {
	Revision int `json:"revision"`
	System   `json:"system"`
}

type System struct {
	Spec        string `json:"spec"`
	Subject     string `json:"source_uri"`
	SourceGraph string `json:"source_graph"`
}

type AuthKey struct {
	Key string
	sync.Mutex
}

// PostHookResource is a struct for the Search routes
type postHookResource struct{}

// Routes returns the chi.Router
func (rs postHookResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/datasets", listDatasets)
	//r.Delete("/datasets/{spec}", listDatasets)
	r.Get("/input/{id}", showInput)
	r.Get("/output/{id}", showOutput)
	//r.Post("/{hudID}", storeJsonld)

	return r
}

func NewESPostHook(ctx context.Context, hubID string) (*PostHookJob, error) {
	source, err := getSource(ctx, hubID)
	if err != nil {
		return nil, err
	}

	ph := &PostHookJob{
		Graph:   &fragments.SortedGraph{},
		Spec:    source.Spec,
		Deleted: false,
		Subject: source.Subject,
	}

	g := r.NewGraph("")
	err = g.Parse(strings.NewReader(source.SourceGraph), "application/ld+json")
	if err != nil {
		return nil, err
	}

	for t := range g.IterTriples() {
		if !cleanDates(ph.Graph, t) && !cleanEbuCore(ph.Graph, t) {
			ph.Graph.Add(t)
		}
	}

	return ph, nil
}

func GetRoutes() chi.Router {
	return postHookResource{}.Routes()
}

var authKey AuthKey

func getSource(ctx context.Context, hubID string) (*System, error) {

	record, err := index.ESClient().Get().
		Index(c.Config.ElasticSearch.IndexName).
		Type("void_edmrecord").
		Id(hubID).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	var sourceGraph Source
	err = json.Unmarshal(*record.Source, &sourceGraph)
	if err != nil {
		return nil, err
	}
	//log.Printf("%#v", string(sourceGraph.SourceGraph))
	return &sourceGraph.System, nil
}

func showInput(w http.ResponseWriter, r *http.Request) {
	hudID := chi.URLParam(r, "id")
	source, err := getSource(r.Context(), hudID)
	if err != nil {
		if err.(*elastic.Error).Status == 404 {
			//http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		log.Printf("%#v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/ld+json")
	w.Write([]byte(source.SourceGraph))
	return
}

func showOutput(w http.ResponseWriter, r *http.Request) {
	hubID := chi.URLParam(r, "id")
	ph, err := NewESPostHook(r.Context(), hubID)
	if err != nil {
		if err.(*elastic.Error).Status == 404 {
			//http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		log.Printf("%#v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/ld+json")
	ph.Graph.SerializeFlatJSONLD(w)
	//w.Write(ph.Bytes)
	return
}

func storeJsonld(w http.ResponseWriter, r *http.Request) {
	return
}

func deleteDataset(w http.ResponseWriter, r *http.Request) {
	return
}

// ListDatasets returns a list of indexed datasets from the PostHook endpoint.
// It renews the authorisation key when this not valid.
func listDatasets(w http.ResponseWriter, r *http.Request) {
	key, err := getAuthKey()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	url := fmt.Sprintf("%s/api/erfgoedbrabant/brabantcloud", strings.TrimSuffix(c.Config.PostHook.URL, "/"))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Header.Set("Cookie", key.Key)

	req.Header.Set("Content-Type", "application/json")

	var netClient = &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := netClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println(key)

	//render.PlainText(w, r, key.Key)
	return
}

func getAuthKey() (AuthKey, error) {
	if authKey.Key != "" {
		return authKey, nil
	}
	return renewAuthKey()
}

func renewAuthKey() (AuthKey, error) {
	authEndpoint := fmt.Sprintf("%s/data/auth/login", strings.TrimSuffix(c.Config.PostHook.URL, "/"))
	payload := strings.NewReader(fmt.Sprintf(`{"username": "%s", "password": "%s"}`, c.Config.PostHook.UserName, c.Config.PostHook.Password))
	req, err := http.NewRequest("POST", authEndpoint, payload)
	if err != nil {
		return authKey, errors.Wrapf(err, "unable to build authentication request for posthook")
	}
	req.Header.Set("Content-Type", "application/json")

	var netClient = &http.Client{
		Timeout: time.Second * 5,
	}
	resp, err := netClient.Do(req)
	if err != nil {
		log.Printf("Error in posthook auth request: %s", err)
		return authKey, err
	}

	zSid := resp.Header.Get("Set-Cookie")
	if zSid == "" || resp.StatusCode != 200 {
		return authKey, errors.New("unable to get auth key from response")
	}
	authKey.Lock()
	authKey.Key = zSid
	authKey.Unlock()

	return authKey, nil
}
