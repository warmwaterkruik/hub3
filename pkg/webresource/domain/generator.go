package domain

// DerivativeGenerator is the interface that Generators need to implement in order
// to generate Derivatives of the source Digital Objects
type DerivativeGenerator interface {
	Thumbnail(b []byte, operation string) ([]byte, error)
	//GenerateThumbnail(operation string, param string, body *io.Reader) error
	//GenerateDeepZoom(body *io.Reader) error
}
