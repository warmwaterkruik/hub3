package rdf

// Supported RDF serialisation formats
const (
	FormatJSONLD   = "application/ld+json"   // https://www.w3.org/TR/json-ld/
	FormatNTriples = "application/n-triples" // https://www.w3.org/TR/n-triples/
	FormatNQuads   = "application/n-quads"   // https://www.w3.org/TR/n-quads/
	FormatRDFXML   = "application/rdf+xml"   // https://www.w3.org/TR/rdf-syntax-grammar/
	FormatTurtle   = "text/turtle"           // http://www.w3.org/TR/turtle/
	FormatTrig     = "application/trig"      // http://www.w3.org/TR/trig/
	FormatN3       = "text/n3"               // https://www.w3.org/TeamSubmission/n3/
)

// Format encapsulates a RDF serialisation format with support for parsing and writing
// source.
//type Format struct {
//MimeType  string
//Extension string
////Reader    Reader
////Writer    Writer
//}
