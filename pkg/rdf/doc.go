package rdf

//// Term abstracts away the difference between a URI reference and a Blank Node
//type Term interface {
//Type() string
//URI() string
//}

//// BlankNode is a blanknode reference
//type BlankNode struct {
//ID string
//}

//// Type returns the reference type of the BlankNode
//func (bn *BlankNode) Type() string {
//return bnode
//}

//// URI returns the formatted URI of the BlankNode
//func (bn BlankNode) URI() string {
//return fmt.Sprintf("_b:%s", bn.ID)
//}

//// Graph holds a list of Resources.
////
//// It also holds the order in which the Triples are added to the Graph.
//// By default Graph is not safe for concurrent access.
//type Graph struct {
//resources map[string]Resource
//sync.Mutex
//LockEnabled bool
//}

//// Get returns a resource for the requested subject.
//// When no subject is found it returns false
//func (g Graph) Get(subject Term) (*Resource, bool) {
//return nil, false
//}

//// Add stores an RDF triple in its Resource.
//// Entries are deduplicated
//func (g *Graph) Add(subject, predicate Term, object ...Entry) error {
////if g.LockEnabled {
////g.Lock()
////defer g.Unlock()
////}
////s, ok := g.Get(subject)
////if !ok {
////s := &Resource{
////Subject: subject,
////}
////}
////return s.Add(predicate, object)
//return nil
//}

//// Record is a group of resources that can be considered a Search Result.
//// It inlines up to 3 level of linked resources.
//type Record struct {
//PrimaryType Term
//Resource
//}

//// Resource is logical grouping of triples based on its Subject
//type Resource struct {
//Subject    Term
//NamedGraph Term
//Types      []Term
//Predicates map[Term]Entries
//Order      int
//}

//func (r *Resource) Add(predicate Term, object ...Entry) error {
////entries, ok := r.Predicates[predicate]
////if !ok {

////}
//return nil
//}

//type Parse interface {
//Parse(r io.Reader, f Format) (Graph, error)
//}

//type Writer interface {
//Write(w io.Writer) error
//}

//type IterTriples ie {
//Iter(ctx context.Context, triples chan Triple)
//}

//type Adder interface {
//Add(s Term, p Term, o Term) error
//}

//type TripleAdder interface {
//AddTriple(t Triple) error
//}

//type Remove interface {
//Remove(s Term, p Term, o Term) error
//}

//type TripleRemove interface {
//RemoveTriple(t Triple) error
//}

//type GetResource interface {
//Get(subject Term) (Resource, bool)
//}

//type Graph interface {
//Adder
//GetResource
//}
