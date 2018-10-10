package domain

// StoreRequest is the abstraction to store Nested Key-Value documents in the
// Storage Repository
type StoreRequest struct {
	IndexName string
	DocType   string
	DocID     string
	DataDoc   string
	PostHook  bool
}

// NewStoreRequest returns a StoreRequest that can be used as a Builder
func NewStoreRequest() *StoreRequest {
	return &StoreRequest{DocType: "doc"}
}

// Index sets the IndexName target for the Request
func (sr *StoreRequest) Index(name string) *StoreRequest {
	sr.IndexName = name
	return sr
}

// ID sets the unique identifier for the record being stored
func (sr *StoreRequest) ID(id string) *StoreRequest {
	sr.DocID = id
	return sr
}

// Type sets the docType for the StoreRequest
func (sr *StoreRequest) Type(docType string) *StoreRequest {
	sr.DocType = docType
	return sr
}

// Doc sets the data that needs to be stored.
// To optimize store performance the doc needs to be JSON marshalled to a string
func (sr *StoreRequest) Doc(doc string) *StoreRequest {
	sr.DataDoc = doc
	return sr
}
