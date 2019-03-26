package rest

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/delving/webresource/pkg/engine"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"willnorris.com/go/imageproxy"
)

const (
	orgIDParam   = "orgID"
	specParam    = "spec"
	localIDParam = "localID"
)

type service struct {
	ctx        context.Context
	client     *storage.Client
	bucket     *storage.BucketHandle
	bucketName string
}

func newService(bucketName, projectID string) (*service, error) {
	var err error
	ctx := context.Background()
	s := &service{
		ctx:        ctx,
		bucketName: bucketName,
	}

	// Creates a client.
	s.client, err = storage.NewClient(s.ctx)
	if err != nil {
		log.Printf("Failed to create client: %v", err)
		return nil, err
	}

	// Creates a Bucket instance.
	s.bucket = s.client.Bucket(bucketName)

	if _, err = s.bucket.Attrs(s.ctx); err != nil {
		// Creates the new bucket.
		attrs := &storage.BucketAttrs{
			//StorageClass: "MULTI-REGIONAL",
			Location: "europe-west1",
		}

		if err := s.bucket.Create(s.ctx, projectID, attrs); err != nil {
			log.Printf("Failed to create bucket: %v", err)
			return nil, errors.Wrap(err, "Failed to create bucket")
		}

		log.Printf("Bucket %v created.\n", bucketName)
	}

	return s, nil
}

// Handler returns a REST router
func Handler(oldService engine.Service) http.Handler {

	cors := cors.New(cors.Options{
		//AllowedOrigins: []string{"*"},
		AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	r := chi.NewRouter()
	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(10 * time.Second))

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	projectID := os.Getenv("PROJECT_ID")
	bucketName := os.Getenv("BUCKET_NAME")
	if projectID == "" || bucketName == "" {
		log.Fatal("PROJECT_ID and BUCKET_NAME must be set as environment variables")
	}
	service, err := newService(bucketName, projectID)
	if err != nil {
		log.Fatalf("Unable to start service: %#v", err)
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})

	r.Mount("/thumbnail", DerivativeResource{}.Routes(service))
	r.Mount("/deepzoom", DerivativeResource{}.Routes(service))
	// Setup the static directory
	//FileServer(r, "/static", getAbsolutePathToFileDir(s.))
	return r
}

// DerivativeResource is the router struct for all derivative links
type DerivativeResource struct{}

// Routes returns the chi.Router
func (rs DerivativeResource) Routes(service *service) chi.Router {
	r := chi.NewRouter()
	r.Get("/{orgID}/{spec}/{localID}*", renderDerivative(service))
	return r
}

func (s *service) getObject(path string) (*storage.ObjectHandle, string, error) {

	ext := filepath.Ext(path)
	if ext == "" {
		path = fmt.Sprintf("%s.jpg", path)
	}
	derivative := s.bucket.Object(path)
	_, err := derivative.Attrs(s.ctx)
	if err != nil {
		log.Printf("derivative not found: %s", path)
		return nil, "", err
	}
	return derivative, path, nil
}

func (s *service) createDerivative(webPath string) (string, error) {
	log.Printf("start creating derivative for: %s", webPath)
	parts := strings.Split(webPath, "/")
	rawPath := strings.Join(parts[:len(parts)-1], "/")
	thumbSource := fmt.Sprintf("%s/2000.jpg", rawPath)
	derivative, path, err := s.getObject(thumbSource)
	if err != nil || path == "" {
		return "", errors.Wrapf(err, "cannot find thumbnail source: %s", thumbSource)
	}
	ctx := context.Background()

	r, err := derivative.NewReader(ctx)
	if err != nil {
		return "", errors.Wrapf(err, "cannot read: %s", path)
	}
	options := parts[len(parts)-1]
	opt := imageproxy.ParseOptions(options)
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", errors.Wrap(err, "cannot convert reader to bytes")
	}

	thumb, err := imageproxy.Transform(b, opt)
	if err != nil {
		return "", errors.Wrapf(err, "cannot transform thumbnail %s with %s", path, options)
	}
	thumbPath := fmt.Sprintf("%s/%s.jpg", rawPath, options)
	w := s.bucket.Object(thumbPath).NewWriter(ctx)
	defer w.Close()
	_, err = w.Write(thumb)
	if err != nil {
		return "", errors.Wrapf(err, "cannot write derivative %s", thumbPath)
	}

	return thumbPath, nil
}

// publicURI returns the public URI for the derivative
func (s *service) publicURI(webPath string) (string, error) {

	_, path, err := s.getObject(webPath)
	if err != nil || path == "" {
		if filepath.Ext(webPath) != "" {
			return "", err
		}
		path, err = s.createDerivative(webPath)
		if err != nil {
			return "", err
		}
	}
	url := fmt.Sprintf("https://media.delving.io/%s", path)
	log.Printf("public URI: %s", url)
	return url, nil
}

func renderDerivative(s *service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		orgID := chi.URLParam(r, orgIDParam)
		spec := chi.URLParam(r, specParam)
		localID := normalise(chi.URLParam(r, localIDParam))
		path := fmt.Sprintf("%s/%s/%s", orgID, spec, localID)
		url, err := s.publicURI(path)
		if err != nil {
			render.Status(r, http.StatusNotFound)
			render.PlainText(w, r, fmt.Sprintf("not found in %#v", err))
			return
		}
		http.Redirect(w, r, url, http.StatusFound)

		return
	}
}

var replacer = strings.NewReplacer(" ", "_")

// normalise converts a path to a media to normalised form.
// this is applied both to the storage and retrieve key
// stripSuffix determines if the storage path is return with a file-type suffix
func normalise(path string) string {
	path = strings.ToLower(path)
	return replacer.Replace(path)
}

//r.get("/{orgid}/{spec}/{localid}*", renderderivative(service))

//// thumbnailresource is the router struct for thumbnail links
//type thumbnailresource struct{}

//// routes returns the chi.router
//func (rs thumbnailresource) routes(service engine.service) chi.router {
//r := chi.newrouter()

//r.get("/{orgid}/{spec}/{localid}*", renderderivative(service))
//return r
//}
//func addBeer(s adding.Service) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
//return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

//func renderThumbnail(w http.ResponseWriter, r *http.Request) {

//wr := domain.WebResource{
//OrgID:        chi.URLParam(r, orgIDParam),
//Spec:         chi.URLParam(r, specParam),
//OriginalPath: chi.URLParam(r, localIDParam),
//Operation:    chi.URLParam(r, "operation"),
//}
//log.Printf("wr: %#v", wr)

//render.JSON(w, r, wr)
//return
//}

//func getAbsolutePathToFileDir(relativePath string) http.Dir {
//workDir, _ := os.Getwd()
//filesDir := filepath.Join(workDir, relativePath)
//return http.Dir(filesDir)
//}

//// FileServer conveniently sets up a http.FileServer handler to serve
//// static files from a http.FileSystem.
//func FileServer(r chi.Router, path string, root http.FileSystem) {
//if strings.ContainsAny(path, "{}*") {
//log.Fatalf("FileServer does not permit URL parameters: %s", path)
//}

//fs := http.StripPrefix(path, http.FileServer(root))

//if path != "/" && path[len(path)-1] != '/' {
//r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
//path += "/"
//}
//path += "*"

//r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//fs.ServeHTTP(w, r)
//}))
//}

//// WebResourceAPIResource is the router struct for webresource data
//type WebResourceAPIResource struct{}

//// Routes returns the chi.Router
//func (wra WebResourceAPIResource) Routes(service engine.Service) chi.Router {
//r := chi.NewRouter()

//r.Get("/{orgID}/{urn}*", listWebResource)
//return r
//}

//func listWebResource(w http.ResponseWriter, r *http.Request) {
//urn := chi.URLParam(r, "urn")
//if strings.HasSuffix(urn, "__") {
////path := filepath.Join(c.Config.WebResource.WebResourceDir, strings.TrimPrefix(urn, "urn:"))
////log.Printf(path)
////matches, err := filepath.Glob(fmt.Sprintf("%s*", path))
////if err != nil {
////log.Printf("%v", err)
////}
////log.Printf("matches: %s", matches)
//}
//log.Printf("urn: %s", urn)
//render.JSON(w, r, `{"type": "thumbnail"}`)
//return
//}

//// DeepZoomResource is the router struct for DeepZoom paths
//type DeepZoomResource struct{}

//// Routes returns the chi.Router
//func (rs DeepZoomResource) Routes(service engine.Service) chi.Router {
//r := chi.NewRouter()

//r.Get("/{orgId}/{spec}/{localId}.tif.dzi", renderDeepZoom)
//r.Get("/{orgId}/{spec}/{localId}.dzi", renderDeepZoom)
//r.Get("/{orgId}/{spec}/{localId}_files/{level}/{col}_{row}.{tile_format}", renderDeepZoomTiles)
//r.Get("/{orgId}/{spec}/{localId}_.tif.files/{level}/{col}_{row}.{tile_format}", renderDeepZoomTiles)
//return r
//}

//// ExploreResource is the router struct for DeepZoom paths
//type ExploreResource struct{}

//// Routes returns the chi.Router
//func (rs ExploreResource) Routes(service engine.Service) chi.Router {
//r := chi.NewRouter()
//r.Get("/", func(w http.ResponseWriter, r *http.Request) {
//render.PlainText(w, r, `{"type": "explore"}`)
//return
//})
//r.Get("/index", func(w http.ResponseWriter, r *http.Request) {
////err := mediamanager.IndexWebResources(s)
////if err != nil {
////log.Printf("Unable to index webresources: %s", err)
////}
//return
//})
//return r
//}

//func renderDeepZoom(w http.ResponseWriter, r *http.Request) {
//render.PlainText(w, r, `{"type": "deepzoom tiles"}`)
//return
//}

//func renderDeepZoomTiles(w http.ResponseWriter, r *http.Request) {
//render.PlainText(w, r, `{"type": "deepzoom tiles"}`)
//return
//}
