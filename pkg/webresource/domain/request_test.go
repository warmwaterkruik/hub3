package domain

import (
	"reflect"
	"testing"
)

func Test_normalise(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name        string
		args        args
		want        string
		wantErr     bool
		stripSuffix bool
	}{
		{"no conversion applied", args{path: "myimage"}, "myimage", false, true},
		{"white space conversion", args{path: "my image"}, "my_image", false, true},
		{"case conversion", args{path: "My Image"}, "my_image", false, true},
		{"strip file-type extension", args{path: "enba A.jpg"}, "enba_a", false, true},
		{"don't strip file-type extension", args{path: "enba A.jpg"}, "enba_a.jpg", false, false},
		{"path with forward slash should be preserved", args{path: "path/enba A.jpg"}, "path/enba_a", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalise(tt.args.path, tt.stripSuffix)
			if got != tt.want {
				t.Errorf("normalise() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkNormalise(b *testing.B) {
	for i := 0; i < b.N; i++ {
		normalise("My Image.jpg", true)
	}
}

func Test_hash(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test uri", args{"http://example.com/path/123.jpg"}, "78a4ba5bce25d177"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hash(tt.args.uri); got != tt.want {
				t.Errorf("hash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		hash("http://example.com/path/123.jpg")
	}
}

func TestWebResource_StoreKey(t *testing.T) {
	type fields struct {
		OrgID     string
		Spec      string
		SourceKey string
		Remote    bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"local uri", fields{"org", "spec", "123 A", false}, "org#spec#123_a"},
		{"remote uri", fields{"org", "spec", "http://example.com/123.jpg", true}, "org#spec#5342d62857bd66f8"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wr := &Request{
				OrgID:     tt.fields.OrgID,
				Spec:      tt.fields.Spec,
				SourceKey: tt.fields.SourceKey,
				Remote:    tt.fields.Remote,
			}
			if got := wr.StoreKey(); got != tt.want {
				t.Errorf("WebResource.StoreKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkStoreKey(b *testing.B) {
	wr := NewRequest("orgID", "spec", "123 a.jpg", false)
	for i := 0; i < b.N; i++ {
		wr.StoreKey()
	}
}

func TestNewWebResource(t *testing.T) {
	type args struct {
		orgID     string
		spec      string
		sourceKey string
		remote    bool
	}
	tests := []struct {
		name string
		args args
		want *Request
	}{
		{"simple init", args{"org", "spec", "123.jpg", false}, &Request{"org", "spec", "123.jpg", false}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRequest(
				tt.args.orgID,
				tt.args.spec,
				tt.args.sourceKey,
				tt.args.remote); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWebResource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebResource_DerivativeKey(t *testing.T) {
	type fields struct {
		OrgID     string
		Spec      string
		SourceKey string
		Remote    bool
	}
	type args struct {
		mod string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"thumbnail",
			fields{"orgID", "spec", "123 a.jpg", false},
			args{"200"},
			"orgID#spec#123_a#200",
		},
		{
			"dzi",
			fields{"orgID", "spec", "123 a.jpg", false},
			args{"dzi"},
			"orgID#spec#123_a#dzi",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wr := &Request{
				OrgID:     tt.fields.OrgID,
				Spec:      tt.fields.Spec,
				SourceKey: tt.fields.SourceKey,
				Remote:    tt.fields.Remote,
			}
			if got := wr.DerivativeKey(tt.args.mod); got != tt.want {
				t.Errorf("WebResource.DerivativeKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
