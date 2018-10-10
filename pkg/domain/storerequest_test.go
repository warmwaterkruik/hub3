package domain

import (
	"reflect"
	"testing"
)

func TestNewStoreRequest(t *testing.T) {
	tests := []struct {
		name string
		want *StoreRequest
	}{
		{"simple create", &StoreRequest{
			IndexName: "",
			DocType:   "doc",
			DocID:     "",
			DataDoc:   "",
			PostHook:  false,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStoreRequest(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStoreRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
