package domain

import (
	"reflect"
	"testing"
)

func TestNewStoreRequest(t *testing.T) {
	type args struct {
		id      string
		docType string
		doc     interface{}
		index   string
	}
	tests := []struct {
		name    string
		args    args
		want    *StoreRequest
		wantErr bool
	}{
		{"empty doc", args{"123", "", nil, ""}, nil, true},
		{"empty id", args{"", "", map[string]interface{}{"key": "value"}, ""}, nil, true},
		{"empty id", args{"123", "", map[string]interface{}{"key": "value"}, "index"},
			&StoreRequest{"doc", "123", "index", map[string]interface{}{"key": "value"}, false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewStoreRequest(tt.args.id, tt.args.docType, tt.args.index, tt.args.doc)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStoreRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStoreRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
