package native

import (
	"willnorris.com/go/imageproxy"
)

type Generator struct{}

// Thumbnail Generates a Thumbnail of the source Digital Object
func (g Generator) Thumbnail(b []byte, options string) ([]byte, error) {
	opt := imageproxy.ParseOptions(options)
	return imageproxy.Transform(b, opt)
}
