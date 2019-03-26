package vips

import "fmt"

type Generator struct{}

// Thumbnail Generates a Thumbnail of the source Digital Object
func (g Generator) Thumbnail(b []byte, w int) ([]byte, error) {
	return nil, fmt.Errorf("not implemented yet")
	//return bimg.NewImage(b).Thumbnail(w)
}
