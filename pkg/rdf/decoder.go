package rdf

import (
	"fmt"
	"io"
)

// A ParseOption allows to customize the behaviour of a decoder.
type ParseOption int

// Options which can configure a decoder.
const (
	// Base IRI to resolve relative IRIs against (for formats that support
	// relative IRIs: Turtle, RDF/XML, TriG, JSON-LD)
	Base ParseOption = iota

	// Strict mode determines how the decoder responds to errors.
	// When true (the default), it will fail on any malformed input. When
	// false, it will try to continue parsing, discarding only the malformed
	// parts.
	// Strict

	// ErrOut
)

// TripleDecoder parses RDF documents (serializations of an RDF graph).
//
// For streaming parsing, use the Decode() method to decode a single Triple
// at a time. Or, if you want to read the whole document in one go, use DecodeAll().
//
// The decoder can be instructed with numerous options. Note that not all options
// are supported by all formats. Consult the following table:
//
//  Option      Description        Value      (default)       Format support
//  ------------------------------------------------------------------------------
//  Base        Base IRI           IRI        (empty IRI)     Turtle, RDF/XML
//  Strict      Strict mode        true/false (true)          TODO
//  ErrOut      Error output       io.Writer  (nil)           TODO
type TripleDecoder interface {
	// Decode parses a RDF document and return the next valid triple.
	// It returns io.EOF when the whole document is parsed.
	Decode() (Triple, error)

	// DecodeAll parses the entire RDF document and return all valid
	// triples, or an error.
	DecodeAll() ([]Triple, error)

	// SetOption sets a parsing option to the given value. Not all options
	// are supported by all serialization formats.
	SetOption(ParseOption, interface{}) error
}

// NewTripleDecoder returns a new TripleDecoder capable of parsing triples
// from the given io.Reader in the given serialization format.
func NewTripleDecoder(r io.Reader, f Format) TripleDecoder {
	switch f {
	case NTriples:
		return newNTDecoder(r)
	case RDFXML:
		return newRDFXMLDecoder(r)
	case Turtle:
		return newTTLDecoder(r)
	default:
		panic(fmt.Errorf("Decoder for serialization format %v not implemented", f))
	}
}
