package hub3

import "fmt"

// Deleter is a generic interface for 'deleting' from the storage layer
type Deleter interface {
	Delete() error
}

// BulkProcessor is a generic interface to 'add' to the storage layer
// todo later change this for generic with FragmentGraph
type BulkProcessor interface {
	Add(record Record) error
}

// BulkableRequest is a generic interface to bulkable requests.
type BulkableRequest interface {
	fmt.Stringer
	Source() ([]string, error)
}

type HubConfig struct{}

// Tenant is the top data structure for the Hub3
// A Hub3  instance holds one or more Tenant tenants
// All Tenants are read-write and need to be created via the API
// Their information is stored in the Configuration storage
type Tenant struct {
	OrgID string
	//Collections []Collection
	//Config      TenantConfig
}

type Collection struct {
	Tenant
	Spec         string
	Revision     int64
	ErrorMessage string
	//Config   CollectionConfig
	//SourceType
}

type Record struct {
}

type Resource struct {
}

type Entry struct {
}

type RecordStore interface {
	Add(record Record)
	Drop(hubId string)
	DropOrphans(spec string)
	DropAll(spec string)
}

type QuadStore interface {
	Add(record Record)
	AddTriple(triple string, context string)
	Drop(hubId string)
	DropTriple(hubId string)
	DropGraph(context string)
	DropOrphans(spec string)
	DropAll(spec string)
}

/*

 */
