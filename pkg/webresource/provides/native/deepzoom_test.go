package native

import (
	"reflect"
	"testing"
)

func Test_numLevels(t *testing.T) {
	type args struct {
		height int
		width  int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"default", args{447, 640}, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := numLevels(tt.args.width, tt.args.height); got != tt.want {
				t.Errorf("numLevels() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_scale(t *testing.T) {
	type args struct {
		level    int
		maxLevel int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"level 10", args{10, 10}, 1.0},
		{"level 1", args{1, 10}, 0.001953125},
		{"level 5", args{5, 10}, 0.03125},
		{"level 0", args{0, 10}, 0.0009765625},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := scale(tt.args.level, tt.args.maxLevel); got != tt.want {
				t.Errorf("scale() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dimensions(t *testing.T) {
	type args struct {
		level    int
		height   int
		width    int
		maxLevel int
	}
	tests := []struct {
		name  string
		args  args
		want1 int
		want  int
	}{
		{"level 0", args{0, 447, 640, 10}, 1, 1},
		{"level 1", args{1, 447, 640, 10}, 1, 2},
		{"level 5", args{5, 447, 640, 10}, 14, 20},
		{"level 10", args{10, 447, 640, 10}, 447, 640},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := dimensions(tt.args.level, tt.args.width, tt.args.height, tt.args.maxLevel)
			if got != tt.want {
				t.Errorf("dimensions() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("dimensions() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDeepZoomDescriptor_numTiles(t *testing.T) {
	type fields struct {
		Opts      DZOpts
		Width     int
		Height    int
		numLevels int
	}
	type args struct {
		level int
	}

	f := fields{
		Opts: DZOpts{
			Format:   "jpeg",
			Overlap:  1,
			TileSize: 254,
		},
		Width:     640,
		Height:    447,
		numLevels: 10,
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   Tile
	}{
		{"level 10", f, args{10}, Tile{10, 3, 2}},
		{"level 9", f, args{9}, Tile{9, 2, 1}},
		{"level 8", f, args{8}, Tile{8, 1, 1}},
		{"level 1", f, args{1}, Tile{1, 1, 1}},
		{"level 0", f, args{0}, Tile{0, 1, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dzd := DeepZoomDescriptor{
				Opts:      tt.fields.Opts,
				Width:     tt.fields.Width,
				Height:    tt.fields.Height,
				numLevels: tt.fields.numLevels,
			}
			tile := dzd.numTiles(tt.args.level)
			if got := tile; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDZI() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestDeepZoomDescriptor_tileBounds(t *testing.T) {
	type fields struct {
		Opts      DZOpts
		Width     int
		Height    int
		numLevels int
	}
	f := fields{
		Opts: DZOpts{
			Format:   "jpeg",
			Overlap:  1,
			TileSize: 254,
		},
		Width:     1955,
		Height:    2404,
		numLevels: 12,
	}
	type args struct {
		level  int
		column int
		row    int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   TileBound
	}{
		{"0/0_0", f, args{0, 0, 0}, TileBound{0, 0, 0, 1, 1}},
		{"1/0_0", f, args{1, 0, 0}, TileBound{1, 0, 0, 1, 2}},
		{"2/0_0", f, args{2, 0, 0}, TileBound{2, 0, 0, 2, 3}},
		{"3/0_0", f, args{3, 0, 0}, TileBound{3, 0, 0, 4, 5}},
		{"4/0_0", f, args{4, 0, 0}, TileBound{4, 0, 0, 8, 10}},
		{"5/0_0", f, args{5, 0, 0}, TileBound{5, 0, 0, 16, 19}},
		{"6/0_0", f, args{6, 0, 0}, TileBound{6, 0, 0, 31, 38}},
		{"7/0_0", f, args{7, 0, 0}, TileBound{7, 0, 0, 62, 76}},
		{"8/0_0", f, args{8, 0, 0}, TileBound{8, 0, 0, 123, 151}},
		{"9/0_0", f, args{9, 0, 0}, TileBound{9, 0, 0, 245, 255}},
		{"9/0_1", f, args{9, 0, 1}, TileBound{9, 0, 253, 245, 48}},
		{"10/0_0", f, args{10, 0, 0}, TileBound{10, 0, 0, 255, 255}},
		{"10/0_1", f, args{10, 0, 1}, TileBound{10, 0, 253, 255, 256}},
		{"10/1_0", f, args{10, 1, 0}, TileBound{10, 253, 0, 236, 255}},
		{"10/1_1", f, args{10, 1, 1}, TileBound{10, 253, 253, 236, 256}},
		{"10/0_2", f, args{10, 0, 2}, TileBound{10, 0, 507, 255, 94}},
		{"10/1_2", f, args{10, 1, 2}, TileBound{10, 253, 507, 236, 94}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dzd := DeepZoomDescriptor{
				Opts:      tt.fields.Opts,
				Width:     tt.fields.Width,
				Height:    tt.fields.Height,
				numLevels: tt.fields.numLevels,
			}
			if got := dzd.tileBounds(tt.args.level, tt.args.column, tt.args.row); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeepZoomDescriptor.tileBounds() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeepZoomDescriptor_tiles(t *testing.T) {
	type fields struct {
		Opts      DZOpts
		Width     int
		Height    int
		numLevels int
	}
	f := fields{
		Opts: DZOpts{
			Format:   "jpeg",
			Overlap:  1,
			TileSize: 254,
		},
		Width:     640,
		Height:    447,
		numLevels: 10,
	}
	type args struct {
		level int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Tile
	}{
		{"0", f, args{0}, []Tile{Tile{0, 0, 0}}},
		{"1", f, args{1}, []Tile{Tile{1, 0, 0}}},
		{"2", f, args{2}, []Tile{Tile{2, 0, 0}}},
		{"8", f, args{8}, []Tile{Tile{8, 0, 0}}},
		{"9", f, args{9}, []Tile{Tile{9, 0, 0}, Tile{9, 1, 0}}},
		{"10", f, args{10}, []Tile{
			Tile{10, 0, 0}, Tile{10, 0, 1},
			Tile{10, 1, 0}, Tile{10, 1, 1},
			Tile{10, 2, 0}, Tile{10, 2, 1},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dzd := DeepZoomDescriptor{
				Opts:      tt.fields.Opts,
				Width:     tt.fields.Width,
				Height:    tt.fields.Height,
				numLevels: tt.fields.numLevels,
			}
			if got := dzd.tiles(tt.args.level); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeepZoomDescriptor.tiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeepZoomDescriptor_TilePath(t *testing.T) {
	type fields struct {
		Opts      DZOpts
		Width     int
		Height    int
		numLevels int
		basePath  string
	}
	f := fields{
		Opts: DZOpts{
			Format:   "jpeg",
			Overlap:  1,
			TileSize: 25,
		},
		Width:     640,
		Height:    447,
		numLevels: 10,
		basePath:  "/deepzoom/orgID/spec/123",
	}
	type args struct {
		tile Tile
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"0/0_0", f, args{Tile{0, 0, 0}}, "/deepzoom/orgID/spec/123_files/0/0_0.jpeg"},
		{"10/1_2", f, args{Tile{10, 1, 2}}, "/deepzoom/orgID/spec/123_files/10/1_2.jpeg"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dzd := DeepZoomDescriptor{
				Opts:      tt.fields.Opts,
				Width:     tt.fields.Width,
				Height:    tt.fields.Height,
				numLevels: tt.fields.numLevels,
				basePath:  tt.fields.basePath,
			}
			if got := dzd.TilePath(tt.args.tile); got != tt.want {
				t.Errorf("DeepZoomDescriptor.TilePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeepZoomDescriptor_DZIPath(t *testing.T) {
	type fields struct {
		Opts      DZOpts
		Width     int
		Height    int
		numLevels int
		basePath  string
	}
	f := fields{
		Opts: DZOpts{
			Format:   "jpeg",
			Overlap:  1,
			TileSize: 254,
		},
		Width:     640,
		Height:    447,
		numLevels: 10,
		basePath:  "/orgID/spec/123",
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"dzi", f, "/orgID/spec/123.dzi"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dzd := DeepZoomDescriptor{
				Opts:      tt.fields.Opts,
				Width:     tt.fields.Width,
				Height:    tt.fields.Height,
				numLevels: tt.fields.numLevels,
				basePath:  tt.fields.basePath,
			}
			if got := dzd.DZIPath(); got != tt.want {
				t.Errorf("DeepZoomDescriptor.DZIPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
