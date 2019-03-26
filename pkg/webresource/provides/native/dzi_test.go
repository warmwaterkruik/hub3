// Copyright 2018 Delving B.V. All rights reserved.
// Use of this source code is governed by a
// license that can be found in the LICENSE file.

package native

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDZI(t *testing.T) {
	type args struct {
		o      DZOpts
		width  string
		height string
	}
	o := DZOpts{
		Format:   "jpg",
		Overlap:  2,
		TileSize: 256,
	}
	tests := []struct {
		name string
		args args
		want DZI
	}{
		{
			"simple",
			args{o, "500", "250"},
			DZI{
				XMLName: xml.Name{
					Space: "",
					Local: "",
				},
				AttrFormat:   "jpg",
				AttrOverlap:  "2",
				AttrTileSize: "256",
				Attrxmlns:    "http://schemas.microsoft.com/deepzoom/2008",
				DZISize: &DZISize{
					XMLName: xml.Name{
						Space: "",
						Local: "",
					},
					AttrHeight: "250",
					AttrWidth:  "500",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDZI(tt.args.o, tt.args.width, tt.args.height); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDZI() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestDZI_Bytes(t *testing.T) {
	type fields struct {
		XMLName      xml.Name
		AttrFormat   string
		AttrOverlap  string
		AttrTileSize string
		Attrxmlns    string
		DZISize      *DZISize
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			"simple",
			fields{
				XMLName: xml.Name{
					Space: "",
					Local: "",
				},
				AttrFormat:   "jpg",
				AttrOverlap:  "2",
				AttrTileSize: "256",
				Attrxmlns:    "http://schemas.microsoft.com/deepzoom/2008",
				DZISize: &DZISize{
					AttrHeight: "500",
					AttrWidth:  "220",
				},
			},
			`<?xml version="1.0" encoding="UTF-8"?>
    <Image Format="jpg" Overlap="2" TileSize="256" xmlns="http://schemas.microsoft.com/deepzoom/2008">
        <Size Height="500" Width="220"></Size>
    </Image>`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dzi := DZI{
				XMLName:      tt.fields.XMLName,
				AttrFormat:   tt.fields.AttrFormat,
				AttrOverlap:  tt.fields.AttrOverlap,
				AttrTileSize: tt.fields.AttrTileSize,
				Attrxmlns:    tt.fields.Attrxmlns,
				DZISize:      tt.fields.DZISize,
			}
			got, err := dzi.Bytes()
			if (err != nil) != tt.wantErr {
				t.Errorf("DZI.Bytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotString := strings.TrimSpace(fmt.Sprintf("%s", got))
			assert.Equal(t, gotString, tt.want, "dzi XML should be equal")
		})
	}
}
