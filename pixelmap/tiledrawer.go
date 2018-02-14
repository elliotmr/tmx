package pixelmap

import (
	"github.com/faiface/pixel"
	"github.com/pkg/errors"
)

type TileDrawer struct {
	ts      *TileSets
	li      *LayerInfo
	drawers map[uint32]*pixel.Drawer
}

func NewTileDrawer(li *LayerInfo, ts *TileSets) *TileDrawer {
	ld := &TileDrawer{
		ts:      ts,
		li:      li,
		drawers: make(map[uint32]*pixel.Drawer),
	}

	for gid, pic := range ts.pics {
		ld.drawers[gid] = &pixel.Drawer{
			Triangles: &pixel.TrianglesData{},
			Picture:   pic,
		}
	}
	return ld
}

func (ld *TileDrawer) Draw(t pixel.Target) error {
	// TODO: draworder
	iter, err := ld.li.layer.Data.Iter()
	if err != nil {
		return errors.Wrap(err, "unable to load layer iterator")
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
			vx, vy, _ := ld.li.CellCoordinates(cellIndex)
			triangleSlice := drawer.Triangles.Slice((i-1)*6, i*6)
			ld.ts.FillTileAndMod(iter.Get().GID, pixel.V(vx, vy), ld.li.color, triangleSlice)
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
