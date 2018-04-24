package ebitentmx

import (
	"image"
	_ "image/png" // This is required for the parsing png resource files
	"path/filepath"
	"os"

	"github.com/hajimehoshi/ebiten"
	"github.com/pkg/errors"
	"github.com/elliotmr/tmx"
)

type tileSetEntry struct {
	rect     *image.Rectangle
	firstGID uint32
	source   string
}

type Resources struct {
	path    string
	entries map[uint32]tileSetEntry
	images  map[string]*ebiten.Image
}

func (r *Resources) loadImage(source string) (string, error) {
	if filepath.IsAbs(source) {
		source = filepath.Clean(source)
	} else {
		source = filepath.Join(r.path, source)
	}
	imageFile, err := os.Open(source)
	if err != nil {
		return "", errors.Wrap(err, "unable to open tileset Image")
	}
	defer imageFile.Close()
	tilesetImg, _, err := image.Decode(imageFile)
	if err != nil {
		return "", errors.Wrap(err, "unable to decode tileset Image")
	}
	pic, _ := ebiten.NewImageFromImage(tilesetImg, ebiten.FilterNearest)
	r.images[source] = pic
	return source, nil
}

func (r *Resources) loadLayer(layer *tmx.Layer) error {
	// Load Images
	if layer.Image != nil {
		_, err := r.loadImage(layer.Image.Source)
		if err != nil {
			return err
		}
	}
	// TODO: Load Templates

	// walk the children recursively.
	for _, child := range layer.Layers {
		err := r.loadLayer(child)
		if err != nil {
			return err
		}
	}
	return nil
}

// LoadResources searches through the tmx map tree and loads any resources found. If
// the resources are located somewhere other than the current working directory, the
// location should be supplied in the path string.
func LoadResources(mapData *tmx.Map, path string) (*Resources, error) {
	// TODO: figure out how to abstract the file system (maybe use Afero?)
	if path == "" {
		path = "."
	}
	r := &Resources{
		path:    path,
		entries: make(map[uint32]tileSetEntry),
		images:  make(map[string]*ebiten.Image),
	}
	for _, set := range mapData.TileSets {
		source, err := r.loadImage(set.Image.Source)
		if err != nil {
			return nil, err
		}
		bounds := r.images[source].Bounds()
		for id := uint32(0); id < set.TileCount; id++ {
			row := id / set.Columns
			col := id % set.Columns
			minX := int(set.Margin + col*(set.TileWidth+set.Spacing))
			minY := int(set.Margin + row*(set.TileHeight+set.Spacing))
			maxX := int(set.Margin + col*(set.TileWidth+set.Spacing) + set.TileWidth)
			maxY := int(set.Margin + row*(set.TileHeight+set.Spacing) + set.TileHeight)
			if minX < bounds.Min.X || minY < bounds.Min.Y || maxX > bounds.Max.X || maxY > bounds.Max.Y {
				return nil, errors.Errorf("tile %d bounds outside of texture bounds (%f, %f, %f, %f)", id, minX, minY, maxX, maxY)
			}
			rect := image.Rect(minX, minY, maxX, maxY)
			r.entries[id+set.FirstGID] = tileSetEntry{
				rect:     &rect,
				firstGID: set.FirstGID,
				source:   source,
			}
		}
	}

	for _, l := range mapData.Layers {
		err := r.loadLayer(l)
		if err != nil {
			return nil, errors.Wrap(err, "unable to load resources")
		}
	}
	return r, nil
}
