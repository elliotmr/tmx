package pixeltmx

import (
	"image"
	"os"

	_ "image/png"

	"github.com/elliotmr/tmx"
	"github.com/faiface/pixel"
	"github.com/pkg/errors"
)

type tileSetEntry struct {
	data     *pixel.TrianglesData
	firstGID uint32
}

// Resources holds the
type Resources struct {
	// TODO: idea - change to resources and add text atlas and template maps
	entries map[uint32]tileSetEntry
	pics    map[uint32]pixel.Picture
}

func (ts *Resources) FillTileAndMod(id uint32, vec pixel.Vec, rbga pixel.RGBA, t pixel.Triangles) {
	_, exists := ts.entries[id]
	if !exists {
		return
	}
	data, ok := ts.entries[id].data.Copy().(*pixel.TrianglesData)
	if !ok {
		return
	}
	for i := range *data {
		(*data)[i].Position = (*data)[i].Position.Add(vec)
		(*data)[i].Color = rbga
	}
	t.Update(data)
}

func LoadResources(mapData *tmx.Map) (*Resources, error) {
	ts := &Resources{
		entries: make(map[uint32]tileSetEntry),
		pics:    make(map[uint32]pixel.Picture),
	}
	for _, set := range mapData.TileSets {
		imageFile, err := os.Open(set.Image.Source)
		if err != nil {
			return nil, errors.Wrap(err, "unable to open tileset image")
		}
		tilesetImg, _, err := image.Decode(imageFile)
		imageFile.Close()
		if err != nil {
			return nil, errors.Wrap(err, "unable to decode tileset image")
		}

		pic := pixel.PictureDataFromImage(tilesetImg)
		ts.pics[set.FirstGID] = pic
		bounds := pic.Bounds()
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
			ts.entries[t.ID+set.FirstGID] = tileSetEntry{
				data:     createTriangleData(pixel.R(minX, minY, maxX, maxY)),
				firstGID: set.FirstGID,
			}
		}
	}

	return ts, nil
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
