package resource

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEntries_Add(t *testing.T) {
	type fields struct {
		entries []Entry
	}
	type args struct {
		object []Entry
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"single add literal",
			fields{[]Entry{Entry{Value: "door", Language: "en"}}},
			args{[]Entry{Entry{Value: "door", Language: "en"}}},
			false,
		},
		{
			"duplicate add literal",
			fields{[]Entry{Entry{Value: "door", Language: "en"}}},
			args{[]Entry{
				Entry{Value: "door", Language: "en"},
				Entry{Value: "door", Language: "en"},
			}},
			false,
		},
		{
			"single add term",
			fields{[]Entry{Entry{ID: "http://ex.com/1"}}},
			args{[]Entry{Entry{ID: "http://ex.com/1"}}},
			false,
		},
		{
			"duplicate add term",
			fields{[]Entry{
				Entry{ID: "http://ex.com/2", SortOrder: 1},
				Entry{ID: "http://ex.com/1", SortOrder: 2},
			}},
			args{[]Entry{
				Entry{ID: "http://ex.com/1", SortOrder: 2},
				Entry{ID: "http://ex.com/2", SortOrder: 1},
				Entry{ID: "http://ex.com/1", SortOrder: 2},
			}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Entries{}
			if err := e.Add(tt.args.object...); (err != nil) != tt.wantErr {
				t.Errorf("Entries.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(e.entries) != len(tt.fields.entries) {
				t.Errorf("Entries.Add() are not the same length: got %v; want %v", len(e.entries), len(tt.fields.entries))
			}
			if !cmp.Equal(e.entries, tt.fields.entries) {
				t.Errorf("Entries are sorted: got %v; want %v", e.entries, tt.fields.entries)

			}
		})
	}
}

func BenchmarkEntriesAdd(b *testing.B) {
	for n := 0; n < b.N; n++ {
		e := &Entries{}
		e.Add(Entry{ID: "http://ex.com/" + string(n), SortOrder: n})
	}
}
