package pixelmap

import (
	"github.com/elliotmr/tiled/tmx"
	"github.com/faiface/pixel"
	"github.com/pkg/errors"
)

type LayerDrawer struct {
	ts      *TileSets
	layer   *tmx.Layer
	mapData *tmx.Map
	drawers map[uint32]*pixel.Drawer
	mat     pixel.Matrix
}

func NewLayerDrawer(mapData *tmx.Map, layerIndex int, ts *TileSets) (*LayerDrawer, error) {
	if layerIndex >= len(mapData.Layers) {
		return nil, errors.Errorf("layer index out of range (%d >= %d)", layerIndex, len(mapData.Layers))
	}
	layer := mapData.Layers[layerIndex]
	ld := &LayerDrawer{
		ts:      ts,
		layer:   layer,
		mapData: mapData,
		drawers: make(map[uint32]*pixel.Drawer),
		mat:     pixel.IM,
	}

	for _, l := range mapData.TileSets {
		ld.drawers[l.FirstGID] = &pixel.Drawer{
			Triangles: &pixel.TrianglesData{},
			Picture:   ts.pics[l.FirstGID],
		}
	}

	return ld, nil
}

func (ld *LayerDrawer) Draw(t pixel.Target) error {
	iter, err := ld.layer.Data.Iter()
	if err != nil {
		return errors.Wrap(err, "unable to load layer iterator")
	}

	// TODO: color, opacity, offsets, draworder
	w := int(ld.mapData.Width)
	if ld.layer.Width != nil {
		w = int(*ld.layer.Width)
	}
	h := int(ld.mapData.Height)
	if ld.layer.Height != nil {
		h = int(*ld.layer.Height)
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
			vx := float64(cellIndex%w) * float64(ld.mapData.TileWidth)
			vy := float64(int(ld.mapData.Height)-cellIndex/h) * float64(ld.mapData.TileHeight)
			triangleSlice := drawer.Triangles.Slice((i-1)*6, i*6)
			ld.ts.FillTileAndMove(iter.Get().GID, pixel.V(vx, vy), triangleSlice)
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
