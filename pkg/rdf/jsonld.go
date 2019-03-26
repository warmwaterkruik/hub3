package rdf

import "io"

type jsonldDecoder struct{}

func newJsonLDEncoder(r io.Reader) *jsonldDecoder {
	return &jsonldDecoder{}
}
