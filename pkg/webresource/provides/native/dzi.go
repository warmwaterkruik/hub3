// Copyright 2018 Delving B.V. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package native

import (
	"encoding/xml"
	"fmt"

	"github.com/pkg/errors"
)

type DZI struct {
	XMLName      xml.Name `xml:"Image,omitempty" json:"Image,omitempty"`
	AttrFormat   string   `xml:"Format,attr"  json:",omitempty"`
	AttrOverlap  string   `xml:"Overlap,attr"  json:",omitempty"`
	AttrTileSize string   `xml:"TileSize,attr"  json:",omitempty"`
	Attrxmlns    string   `xml:"xmlns,attr"  json:",omitempty"`
	DZISize      *DZISize `xml:"http://schemas.microsoft.com/deepzoom/2008 Size,omitempty" json:"Size,omitempty"`
}

type DZISize struct {
	XMLName    xml.Name `xml:"Size,omitempty" json:"Size,omitempty"`
	AttrHeight string   `xml:"Height,attr"  json:",omitempty"`
	AttrWidth  string   `xml:"Width,attr"  json:",omitempty"`
}

type DZOpts struct {
	Format   string
	Overlap  int
	TileSize int
}

// NewDZI returns a DZI. This can be used to generate a DeepZoom Image descriptor.
func NewDZI(o DZOpts, width, height string) DZI {
	size := &DZISize{
		AttrHeight: height,
		AttrWidth:  width,
	}
	return DZI{
		AttrFormat:   o.Format,
		AttrOverlap:  fmt.Sprint(o.Overlap),
		AttrTileSize: fmt.Sprint(o.TileSize),
		Attrxmlns:    "http://schemas.microsoft.com/deepzoom/2008",
		DZISize:      size,
	}
}

// Bytes returns the DZI encoded as XML as a byte Array
func (dzi DZI) Bytes() ([]byte, error) {
	b, err := xml.MarshalIndent(dzi, "    ", "    ")
	if err != nil {
		return nil, errors.Wrap(err, "Unable to marschal the dzi to XML")
	}
	return append([]byte(xml.Header), b...), nil
}
