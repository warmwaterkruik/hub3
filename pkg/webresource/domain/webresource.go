package domain

import (
	"fmt"
	"path/filepath"
	"strings"
)

// WebResourceType determines how and where a WebResource is stored
type WebResourceType string

const (
	THUMBNAIL WebResourceType = "thumbnail"
	DEEPZOOM                  = "deepzoom"
	SOURCE                    = "source"
)

// WebResource contains the elements to build a storage path for a digital object
// stored in the storage Repository
type WebResource struct {
	OrgID        string
	Spec         string
	SourceKey    string
	OriginalPath string
	SourcePath   string
	Kind         WebResourceType
	MimeType     string
	Body         []byte
	Public       bool
	Attrs        map[string]string
	Exif         map[string]string
	Operation    string
	StripSuffix  bool
}

// StoragePath returns the full storage path for the WebResource
func (wr WebResource) StoragePath() string {
	// TODO maybe return error

	base := fmt.Sprintf("/%s/%s/%s/%s",
		wr.Kind,
		wr.OrgID,
		wr.Spec,
		normalise(wr.SourceKey, wr.StripSuffix),
	)
	if wr.Operation != "" {
		return fmt.Sprintf("%s/%s", base, wr.Operation)
	}
	return base
}

// SourcePath returns the path to the source WebResource
//func (wr WebResource) SourcePath() string {
//if wr.SourcePath != "" {
//return wr.SourcePath
//}
//return fmt.Sprintf("/%s/%s/%s/%s",
//SOURCE,
//wr.OrgID,
//wr.Spec,
//normalise(wr.SourceKey, true),
//)
//}

// GetMetadata returns a map with all the custom metadata
func (wr WebResource) GetMetadata(prefix string) map[string]string {
	prefixed := func(name string) string { return fmt.Sprintf("%s%s", prefix, name) }
	attrs := map[string]string{
		prefixed("orgID"):    wr.OrgID,
		prefixed("spec"):     wr.Spec,
		prefixed("origName"): wr.SourceKey,
	}

	for k, v := range wr.Attrs {
		attrs[prefixed(k)] = v
	}

	return attrs
}

// normalise converts a path to a media to normalised form.
// this is applied both to the storage and retrieve key
// stripSuffix determines if the storage path is return with a file-type suffix
func normalise(path string, stripSuffix bool) string {
	path = strings.ToLower(path)
	ext := filepath.Ext(path)
	if ext != "" && stripSuffix {
		path = strings.TrimSuffix(path, ext)
	}
	return replacer.Replace(path)
}
