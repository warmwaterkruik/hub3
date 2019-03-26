package rdf

import "errors"

// ErrFormatNotSupported is the error returned from a Decoder/Encoder when RDF format is not supported
var ErrFormatNotSupported = errors.New("format not supported")

// ErrEncoderClosed is the error returned from Encode() when the Triple/Quad-Encoder is closed
var ErrEncoderClosed = errors.New("Encoder is closed and cannot encode anymore")
