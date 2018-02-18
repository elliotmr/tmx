package pixeltmx

import (
	"github.com/faiface/pixel"
	"github.com/pkg/errors"
)

type tileLayerDrawer struct {
	resources *Resources
	info      *LayerInfo
	drawers   map[uint32]*pixel.Drawer
}

func newTileLayerDrawer(resources *Resources, info *LayerInfo) (*tileLayerDrawer, error) {
	ld := &tileLayerDrawer{
		resources: resources,
		info:      info,
		drawers:   make(map[uint32]*pixel.Drawer),
	}

	for _, entry := range resources.entries {
		_, exists := ld.drawers[entry.firstGID]
		if !exists {
			ld.drawers[entry.firstGID] = &pixel.Drawer{
				Triangles: &pixel.TrianglesData{},
				Picture:   ld.resources.images[entry.source],
			}
		}
	}

	return ld, ld.Update()
}

func (ld *tileLayerDrawer) Type() int {
	return TileLayerDrawer
}

func (ld *tileLayerDrawer) Info() *LayerInfo {
	return ld.info
}

func (ld *tileLayerDrawer) Update() error {
	// TODO: draworder
	iter, err := ld.info.layer.Data.Iter()
	if err != nil {
		return errors.Wrap(err, "unable to load layer iterator")
	}

	for gid, drawer := range ld.drawers {
		i := 1
		for iter.Next() {
			tse := ld.resources.entries[iter.Get().GID]
			if tse.firstGID != gid {
				continue
			}
			if i*6 > drawer.Triangles.Len() {
				drawer.Triangles.SetLen(i * 6)
			}
			cellIndex := int(iter.GetIndex())
			loc, _ := ld.info.TileRect(cellIndex)
			triangleSlice := drawer.Triangles.Slice((i-1)*6, i*6)
			ld.resources.fillTileAndMod(iter.Get().GID, loc.Center(), ld.info.color, triangleSlice)
			i++
		}
		drawer.Triangles.SetLen(i * 6)
	}
	return errors.Wrap(iter.Error(), "unable to iterate through layer")
}

func (ld *tileLayerDrawer) Draw(t pixel.Target) {
	for _, d := range ld.drawers {
		d.Draw(t)
	}
}
