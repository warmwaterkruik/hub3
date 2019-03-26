package native

import (
	"context"
	"fmt"
	"log"
	"math"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"gopkg.in/h2non/bimg.v1"
)

type Store interface {
	Write(path string, b []byte, mimeType string) error
}

// TileBound holds the dimensions for each tile
type TileBound struct {
	Level  int
	X      int
	Y      int
	Width  int
	Height int
}

// Tile holds the column and row information for a deepzoom tile
type Tile struct {
	Level  int
	Column int
	Row    int
}

type TileBuilder struct {
	Tile       Tile
	LevelImage []byte
	ImageSize  bimg.ImageSize
}

// DeepZoomDescriptor deals with all the DeepZoom calculations
type DeepZoomDescriptor struct {
	Opts      DZOpts
	Width     int
	Height    int
	numLevels int
	basePath  string
}

// NewDeepZoomDescriptor creates a DeepZoomDescriptor and sets the max number
// of zoom-levels
func NewDeepZoomDescriptor(o DZOpts, width, height int, basePath string) DeepZoomDescriptor {
	dzd := DeepZoomDescriptor{
		Opts:   o,
		Width:  width,
		Height: height,
	}
	dzd.numLevels = numLevels(width, height)
	dzd.basePath = basePath
	return dzd
}

// TilePath return the storage path for the Tile
func (dzd DeepZoomDescriptor) TilePath(tile Tile) string {
	return fmt.Sprintf(
		"%s_files/%d/%d_%d.%s",
		dzd.basePath,
		tile.Level,
		tile.Column,
		tile.Row,
		dzd.Opts.Format,
	)
}

// DZIPath returns the strorage path for DZI XML
func (dzd DeepZoomDescriptor) DZIPath() string {
	return fmt.Sprintf("%s.dzi", dzd.basePath)
}

// numLevels determines the number of zoom levels for an image based on its
// width and height
func numLevels(width, height int) int {
	maxDimension := math.Max(float64(width), float64(height))
	return int(math.Ceil(math.Log2(maxDimension)))
}

// scale determines the relative scale to the zoom-level
func scale(level, maxLevel int) float64 {
	return math.Pow(0.5, float64(maxLevel-level))
}

// dimensions calculates the width and height of the root-image for the zoom-level
func dimensions(level, width, height, maxLevel int) (int, int) {
	scale := scale(level, maxLevel)
	w := int(math.Ceil(float64(width) * scale))
	h := int(math.Ceil(float64(height) * scale))
	return w, h
}

// numTiles returns the number of tiles for the width and height of a given
// zoom-level
func (dzd DeepZoomDescriptor) numTiles(level int) Tile {
	w, h := dimensions(level, dzd.Width, dzd.Height, dzd.numLevels)
	tileSize := float64(dzd.Opts.TileSize)
	widthTiles := int(math.Ceil(float64(w) / tileSize))
	heightTiles := int(math.Ceil(float64(h) / tileSize))
	return Tile{level, widthTiles, heightTiles}
}

// Tiles returns an array of Tile for the given level.
func (dzd DeepZoomDescriptor) tiles(level int) []Tile {
	tileLevel := dzd.numTiles(level)
	tiles := []Tile{}
	for c := 0; c < tileLevel.Column; c++ {
		for r := 0; r < tileLevel.Row; r++ {
			tiles = append(tiles, Tile{level, c, r})
		}
	}
	return tiles
}

// tileBounds returns the bounds for each tile
func (dzd DeepZoomDescriptor) tileBounds(level, column, row int) TileBound {
	offsetX := 0
	offsetWidth := 1
	if column != 0 {
		offsetX = dzd.Opts.Overlap
		offsetWidth = 2
	}
	offsetY := 0
	offsetHeight := 1
	if row != 0 {
		offsetY = dzd.Opts.Overlap
		offsetHeight = 2
	}
	x := (column * dzd.Opts.TileSize) - offsetX
	y := (row * dzd.Opts.TileSize) - offsetY
	levelWidth, levelHeight := dimensions(level, dzd.Width, dzd.Height, dzd.numLevels)

	w := dzd.Opts.TileSize + (offsetWidth * dzd.Opts.Overlap)
	h := dzd.Opts.TileSize + (offsetHeight * dzd.Opts.Overlap)

	w = int(math.Min(float64(w), float64(levelWidth-x)))
	h = int(math.Min(float64(h), float64(levelHeight-y)))

	return TileBound{
		Level:  level,
		X:      x,
		Y:      y,
		Width:  w,
		Height: h,
	}
}

// getLevelImage resizes the source image to the dimensions of the deepzoom-level
func (dzd DeepZoomDescriptor) getLevelImage(image *bimg.Image, level int) ([]byte, error) {
	width, height := dimensions(level, dzd.Width, dzd.Height, dzd.numLevels)
	//opt := imageproxy.ParseOptions(fmt.Sprintf("%dx%d,fit", width, height))
	//return imageproxy.Transform(image, opt)
	options := bimg.Options{
		Width:   width,
		Height:  height,
		Crop:    false,
		Quality: 95,
	}
	return image.Process(options)
}

// getTileImage crops the levelImage by the TileBound
func getTileImage(levelImage *bimg.Image, size bimg.ImageSize, tb TileBound) ([]byte, error) {
	//func getTileImage(levelImage []byte, size bimg.ImageSize, tb TileBound) ([]byte, error) {
	//opt := imageproxy.Options{
	//Width:          0.0,
	//Height:         0.0,
	//Fit:            false,
	//Rotate:         0,
	//FlipVertical:   false,
	//FlipHorizontal: false,
	//Quality:        0,
	//Signature:      "",
	//ScaleUp:        false,
	//Format:         "jpeg",
	//CropX:          float64(tb.X),
	//CropY:          float64(tb.Y),
	//CropWidth:      float64(tb.Width),
	//CropHeight:     float64(tb.Height),
	//SmartCrop:      false,
	//}
	//log.Printf("tile bound: %#v\n", tb)
	//return imageproxy.Transform(levelImage, opt)
	//tileWidth := tb.Width - tb.X
	//if tileWidth > size.Width {
	//tileWidth = size.Width
	//}
	//tileHeight := tb.Height - tb.Y
	//if tileHeight > size.Height {
	//tileHeight = size.Height
	//}

	tile, err := levelImage.Extract(tb.Y, tb.X, tb.Width, tb.Height)
	if err != nil {
		log.Printf("unable to extract tile: %d; bounds: %#v; dimensions %#v; tile width %d; tile height: %d", tb.Level, tb, size, tb.Width, tb.Height)
		return nil, err
	}
	return tile, nil
}

// storeTile stores a single tile in the Service Repository
func storeTile(tile []byte, path string, s Store) error {
	return s.Write(path, tile, "image/jpeg")
}

// StoreTiles generates all tiles for each deepzoom level
// The engine.Service provides the store interface
func (dzd *DeepZoomDescriptor) StoreTiles(ctx context.Context, path string, image []byte, s Store) error {
	g, ctx := errgroup.WithContext(ctx)
	tiles := make(chan *TileBuilder)

	dzd.basePath = path
	nrTiles := 0

	g.Go(func() error {
		defer close(tiles)
		sourceImage := bimg.NewImage(image)
		for level := 0; level <= dzd.numLevels; level++ {
			levelImage, err := dzd.getLevelImage(sourceImage, level)
			if err != nil {
				return errors.Wrapf(err, "unable to resize source image for level: %d", level)
			}
			newImage := bimg.NewImage(levelImage)
			size, sizeErr := newImage.Size()
			if sizeErr != nil {
				return errors.Wrap(err, "cannot get size for image")
			}
			for _, tile := range dzd.tiles(level) {
				select {
				case tiles <- &TileBuilder{tile, levelImage, size}:
					nrTiles++
				case <-ctx.Done():
					return ctx.Err()

				}
			}
		}
		return nil

	})

	const numDigesters = 6
	for i := 0; i < numDigesters; i++ {
		g.Go(func() error {

			for tb := range tiles {
				tb := tb
				t := tb.Tile
				bounds := dzd.tileBounds(t.Level, t.Column, t.Row)
				levelImage := bimg.NewImage(tb.LevelImage)
				tileImage, err := getTileImage(levelImage, tb.ImageSize, bounds)
				if err != nil {
					return errors.Wrapf(err, "unable to crop level image for level: %d; bounds: %#v; dimensions %#v; path %#v; level size: %#v", t.Level, bounds, tb.ImageSize, dzd.TilePath(t), tb.ImageSize)
				}
				err = storeTile(tileImage, dzd.TilePath(t), s)
				if err != nil {
					return errors.Wrapf(err, "unable to store tile image in storage repository")
				}

			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

// storeDZI generates the deepzoom description in XML and stores it via the engine.Service
func (dzd *DeepZoomDescriptor) StoreDZI(path string, s Store) error {
	dzd.basePath = path
	dzi := NewDZI(dzd.Opts, fmt.Sprintf("%d", dzd.Width), fmt.Sprintf("%d", dzd.Height))
	b, err := dzi.Bytes()
	if err != nil {
		return errors.Wrapf(err, "Unable to generate dzi")
	}

	return s.Write(dzd.DZIPath(), b, "application/xml")

}
