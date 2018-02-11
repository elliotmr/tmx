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
	drawers map[uint32]pixel.Drawer
	mat     pixel.Matrix
}

func NewLayerDrawer(mapData *tmx.Map, layerIndex int, ts *TileSets) (*LayerDrawer, error) {
	if layerIndex >= len(mapData.Layers) {
		return nil, errors.Errorf("layer index out of range (%d >= %d)", layerIndex, len(mapData.Layers))
	}
	layer := mapData.Layers[layerIndex]
	ld := &LayerDrawer{
		ts: ts,
		layer: layer,
		mapData: mapData,
		drawers: make(map[uint32]pixel.Drawer),
		mat: pixel.IM,
	}

	// TODO: make this cleverer
	for _, l := range mapData.TileSets {
		ld.drawers[l.FirstGID] = pixel.Drawer{
			Triangles: &pixel.TrianglesData{},
			Picture: ts.pics[l.FirstGID],
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

	for iter.Next() {
		i := int(iter.GetIndex())
		vx := float64(i%w) * float64(ld.mapData.TileWidth)
		vy := float64(int(ld.mapData.Height)-i/h) * float64(ld.mapData.TileHeight)
		tse := ld.ts.entries[iter.Get().GID]
		drawer := ld.drawers[tse.firstGID]
		drawer.Triangles.SetLen(drawer.Triangles.Len() + 6)
		triangleSlice := drawer.Triangles.Slice(drawer.Triangles.Len() - 6, drawer.Triangles.Len())
		ld.ts.FillTileAndMove(iter.Get().GID, pixel.V(vx, vy), triangleSlice)
	}
	if iter.Error() != nil {
		return errors.Wrap(iter.Error(), "unable to iterate through layer")
	}
	for _, d := range ld.drawers {
		d.Draw(t)
	}
	return nil
}
