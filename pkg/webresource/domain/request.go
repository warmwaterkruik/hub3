package domain

import (
	"fmt"
	"strings"

	"github.com/cespare/xxhash"
)

var replacer = strings.NewReplacer(" ", "_")

const sep = "#"

// Request holds all information to store or find a Request
type Request struct {
	OrgID     string `json:"orgID"`
	Spec      string `json:"spec"`
	SourceKey string `json:"requestKey"` // can be both a path or an URL
	Remote    bool   `json:"remote"`
}

// NewRequest creates a new WebResource
func NewRequest(orgID, spec, sourceKey string, remote bool) *Request {
	return &Request{
		OrgID:     orgID,
		Spec:      spec,
		SourceKey: sourceKey,
		Remote:    remote,
	}
}

// StoreKey returns the normalised and formatted key for storage and retrieval
// Note that the key is made case-insensitive.
// The StoreKey is also going to be used al the IIIF identifier
func (wr *Request) StoreKey() string {
	key := normalise(wr.SourceKey, true)
	if wr.Remote {
		key = hash(key)
	}
	return fmt.Sprintf("%s#%s#%s", wr.OrgID, wr.Spec, key)
}

// DerivativeKey retrieves the key from the derivative store
func (wr *Request) DerivativeKey(mod string) string {
	return fmt.Sprintf("%s#%s", wr.StoreKey(), mod)
}

func hash(uri string) string {
	hash := xxhash.Sum64String(uri)
	return fmt.Sprintf("%016x", hash)
}
