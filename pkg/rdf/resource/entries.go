package resource

import (
	"sort"

	"github.com/google/go-cmp/cmp"
)

const (
	literal EntryType = "Literal"
	uriRef            = "URIRef"
	bnode             = "Bnode"
)

// EntryType determines if an Object Entry is a Literal, UriRef or BNode
type EntryType string

// Entries holds a list of Object Entry linked to a single Predicate
type Entries struct {
	entries []Entry
}

// Add appends a non-duplicate Entry to Entries and makes sure that the Entry List
// is sorted by Entry.SortOrder.
func (e *Entries) Add(object ...Entry) error {
objectLoop:
	for _, o := range object {
		for _, entry := range e.entries {
			if cmp.Equal(entry, o) {
				continue objectLoop
			}
		}
		e.entries = append(e.entries, o)
	}

	// keep entries sorted
	sort.Slice(e.entries, func(i, j int) bool {
		return e.entries[i].SortOrder < e.entries[j].SortOrder
	})
	return nil
}

// Entry is the Object part of the triple
type Entry struct {
	ID        string    `json:"@id,omitempty"`
	Value     string    `json:"@value,omitempty"`
	Language  string    `json:"@language,omitempty"`
	DataType  string    `json:"@type,omitempty"`
	Type      EntryType `json:"entrytype"`
	Resolved  bool      `json:"resolved"`
	SortOrder int       `json:"order"`
	//Inline   *Resource `json:"inline"`
}
