package function

import (
	"reflect"
	"testing"
)

func TestNewStoreKey(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *StoreKey
		wantErr bool
	}{
		{
			"correct",
			args{"test-org/test-spec/123.jpg"},
			&StoreKey{
				Path:      "test-org/test-spec/123.jpg",
				OrgID:     "test-org",
				Spec:      "test-spec",
				LocalPath: "123.jpg",
				Operation: "",
			},
			false,
		},
		{
			"strip source from localPath",
			args{"test-org/test-spec/source/123.jpg"},
			&StoreKey{
				Path:      "test-org/test-spec/source/123.jpg",
				OrgID:     "test-org",
				Spec:      "test-spec",
				LocalPath: "123.jpg",
				Operation: "",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewStoreKey(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStoreKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStoreKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStoreKey_Key(t *testing.T) {
	type fields struct {
		Path      string
		OrgID     string
		Spec      string
		LocalPath string
		Operation string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"remove uppercase",
			fields{
				LocalPath: "123AB.jpg",
			},
			"123ab",
		},
		{"remove whitespace",
			fields{
				LocalPath: "123 AB.jpg",
			},
			"123_ab",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sk := StoreKey{
				Path:      tt.fields.Path,
				OrgID:     tt.fields.OrgID,
				Spec:      tt.fields.Spec,
				LocalPath: tt.fields.LocalPath,
				Operation: tt.fields.Operation,
			}
			got := sk.Key()
			if got != tt.want {
				t.Errorf("StoreKey.Key() = %v, want %v", got, tt.want)
			}
		})
	}
}
