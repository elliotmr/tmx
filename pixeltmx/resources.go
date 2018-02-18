package pixeltmx

import (
	"image"
	"os"

	_ "image/png" // This is required for the parsing png resource files

	"github.com/elliotmr/tmx"
	"github.com/faiface/pixel"
	"github.com/pkg/errors"
	"path/filepath"
)

type tileSetEntry struct {
	data     *pixel.TrianglesData
	firstGID uint32
	source   string
}

// Resources holds all the raw resources required for rendering a TMX map. This
// includes tilesets and tileset pictures, raw images, object templates, etc.
type Resources struct {
	// TODO: add text atlas and template maps
	path    string
	entries map[uint32]tileSetEntry
	images  map[string]pixel.Picture
}

func (r *Resources) loadImage(source string) (string, error) {
	if filepath.IsAbs(source) {
		source = filepath.Clean(source)
	} else {
		source = filepath.Join(r.path, source)
	}
	imageFile, err := os.Open(source)
	if err != nil {
		return "", errors.Wrap(err, "unable to open tileset image")
	}
	defer imageFile.Close()
	tilesetImg, _, err := image.Decode(imageFile)
	if err != nil {
		return "", errors.Wrap(err, "unable to decode tileset image")
	}
	pic := pixel.PictureDataFromImage(tilesetImg)
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
		images:  make(map[string]pixel.Picture),
	}
	for _, set := range mapData.TileSets {
		source, err := r.loadImage(set.Image.Source)
		if err != nil {
			return nil, err
		}
		bounds := r.images[source].Bounds()
		// tmx convention right -> down (origin top left), pixel convetion right -> up (origin bottom left)
		// this means we have to flip the row index
		rows := set.TileCount / set.Columns

		for _, t := range set.Tiles {
			if t.ID >= set.TileCount {
				return nil, errors.Errorf("tile id greater than tilecount (%d > %d)", t.ID, set.TileCount-1)
			}
			row := rows - t.ID/set.Columns - 1
			col := t.ID % set.Columns
			minX := float64(set.Margin + col*(set.TileWidth+set.Spacing))
			minY := float64(set.Margin + row*(set.TileHeight+set.Spacing))
			maxX := float64(set.Margin + col*(set.TileWidth+set.Spacing) + set.TileWidth)
			maxY := float64(set.Margin + row*(set.TileHeight+set.Spacing) + set.TileHeight)
			if minX < bounds.Min.X || minY < bounds.Min.Y || maxX > bounds.Max.X || maxY > bounds.Max.Y {
				return nil, errors.Errorf("tile %d bounds outside of texture bounds (%f, %f, %f, %f)", t.ID, minX, minY, maxX, maxY)
			}
			r.entries[t.ID+set.FirstGID] = tileSetEntry{
				data:     createTriangleData(pixel.R(minX, minY, maxX, maxY)),
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

func (r *Resources) fillTileAndMod(id uint32, vec pixel.Vec, rbga pixel.RGBA, t pixel.Triangles) {
	_, exists := r.entries[id]
	if !exists {
		return
	}
	data, ok := r.entries[id].data.Copy().(*pixel.TrianglesData)
	if !ok {
		return
	}
	for i := range *data {
		(*data)[i].Position = (*data)[i].Position.Add(vec)
		(*data)[i].Color = rbga
	}
	t.Update(data)
}

func createTriangleData(r pixel.Rect) *pixel.TrianglesData {
	tri := pixel.MakeTrianglesData(6)
	halfWidthVec := pixel.V(r.W()/2, 0)
	halfHeightVec := pixel.V(0, r.H()/2)
	(*tri)[0].Position = pixel.Vec{}.Sub(halfWidthVec).Sub(halfHeightVec)
	(*tri)[1].Position = pixel.Vec{}.Add(halfWidthVec).Sub(halfHeightVec)
	(*tri)[2].Position = pixel.Vec{}.Add(halfWidthVec).Add(halfHeightVec)
	(*tri)[3].Position = pixel.Vec{}.Sub(halfWidthVec).Sub(halfHeightVec)
	(*tri)[4].Position = pixel.Vec{}.Add(halfWidthVec).Add(halfHeightVec)
	(*tri)[5].Position = pixel.Vec{}.Sub(halfWidthVec).Add(halfHeightVec)
	for i := range *tri {
		(*tri)[i].Color = pixel.Alpha(1)
		(*tri)[i].Picture = r.Center().Add((*tri)[i].Position)
		(*tri)[i].Intensity = 1
	}
	return tri
}
