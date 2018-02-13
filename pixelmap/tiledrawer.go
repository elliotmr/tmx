package pixelmap

import (
	"github.com/elliotmr/tiled/tmx"
	"github.com/faiface/pixel"
	"github.com/pkg/errors"
)

type TileDrawer struct {
	ts      *TileSets
	layer   *tmx.Layer
	mapData *tmx.Map
	drawers map[uint32]*pixel.Drawer
}

func NewTileDrawer(mapData *tmx.Map, layerIndex int, ts *TileSets) (*TileDrawer, error) {
	if layerIndex >= len(mapData.Layers) {
		return nil, errors.Errorf("layer index out of range (%d >= %d)", layerIndex, len(mapData.Layers))
	}
	layer := mapData.Layers[layerIndex]
	ld := &TileDrawer{
		ts:      ts,
		layer:   layer,
		mapData: mapData,
		drawers: make(map[uint32]*pixel.Drawer),
	}

	for _, l := range mapData.TileSets {
		ld.drawers[l.FirstGID] = &pixel.Drawer{
			Triangles: &pixel.TrianglesData{},
			Picture:   ts.pics[l.FirstGID],
		}
	}

	return ld, nil
}

func (ld *TileDrawer) Draw(t pixel.Target) error {
	iter, err := ld.layer.Data.Iter()
	if err != nil {
		return errors.Wrap(err, "unable to load layer iterator")
	}

	// TODO: color, draworder, extract to function
	w := int(ld.mapData.Width)
	if ld.layer.Width != nil {
		w = int(*ld.layer.Width)
	}
	h := int(ld.mapData.Height)
	if ld.layer.Height != nil {
		h = int(*ld.layer.Height)
	}
	offX := float64(ld.mapData.TileWidth) * 0.5
	if ld.layer.OffsetX != nil {
		offX += *ld.layer.OffsetX
	}
	offY := float64(ld.mapData.TileHeight) * -0.5
	if ld.layer.OffsetY != nil {
		offY -= *ld.layer.OffsetY
	}

	opacity := 1.0
	if ld.layer.Opacity != nil {
		opacity = *ld.layer.Opacity
	}
	rgba := pixel.Alpha(opacity)

	if ld.layer.Visible != nil {
		if *ld.layer.Visible == 0 {
			rgba.A = 0.0
		}
	}

	for gid, drawer := range ld.drawers {
		i := 1
		for iter.Next() {
			tse := ld.ts.entries[iter.Get().GID]
			if tse.firstGID != gid {
				continue
			}
			if i*6 > drawer.Triangles.Len() {
				drawer.Triangles.SetLen(i * 6)
			}
			cellIndex := int(iter.GetIndex())
			vx := float64(cellIndex%w) * float64(ld.mapData.TileWidth) + offX
			vy := float64(int(ld.mapData.Height)-cellIndex/h) * float64(ld.mapData.TileHeight) + offY
			triangleSlice := drawer.Triangles.Slice((i-1)*6, i*6)
			ld.ts.FillTileAndMod(iter.Get().GID, pixel.V(vx, vy), rgba, triangleSlice)
			i++
		}
		drawer.Triangles.SetLen(i * 6)
	}

	if iter.Error() != nil {
		return errors.Wrap(iter.Error(), "unable to iterate through layer")
	}
	for _, d := range ld.drawers {
		d.Draw(t)
	}
	return nil
}
