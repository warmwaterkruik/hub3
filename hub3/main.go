package hub3

import "fmt"

// Deleter
type Deleter interface {
	Delete() error
}

// BulkProcessor is a generic interface to 'add' to the storage layer
// todo later change this for generic with FragmentGraph
type BulkProcessor interface {
	Add(request BulkableRequest) error
}

// BulkableRequest is a generic interface to bulkable requests.
type BulkableRequest interface {
	fmt.Stringer
	Source() ([]string, error)
}
