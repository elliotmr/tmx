package ebitentmx

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/pkg/errors"
)

type tileLayerDrawer struct {
	resources *Resources
	info      *LayerInfo
	opts      *ebiten.DrawImageOptions
	image     *ebiten.Image
}

func newTileLayerDrawer(resources *Resources, info *LayerInfo) (*tileLayerDrawer, error) {
	img, _ := ebiten.NewImage(
		info.w*int(info.mapData.TileWidth),
		info.h*int(info.mapData.TileHeight),
		ebiten.FilterNearest,
	)
	bounds := img.Bounds()
	geom := ebiten.GeoM{}
	geom.Translate(info.offX, info.offY)
	opts := &ebiten.DrawImageOptions{
		SourceRect: &bounds,
		GeoM:       geom,
		Filter:     ebiten.FilterNearest,
	}
	ld := &tileLayerDrawer{
		resources: resources,
		info:      info,
		opts:      opts,
		image:     img,
	}

	return ld, ld.Update()
}

func (ld *tileLayerDrawer) Type() int {
	return TileLayerDrawer
}

func (ld *tileLayerDrawer) Info() *LayerInfo {
	return ld.info
}

func (ld *tileLayerDrawer) Image() *ebiten.Image {
	return ld.image
}

func (ld *tileLayerDrawer) Update() error {
	// TODO: draworder
	iter, err := ld.info.layer.Data.Iter()
	if err != nil {
		return errors.Wrap(err, "unable to load layer iterator")
	}
	for iter.Next() {
		tile := iter.Get()
		if tile.GID() == 0 {
			continue
		}
		tse, exists := ld.resources.entries[tile.GID()]
		if !exists {
			return errors.Errorf("tile with gid '%d' does not exist", tile.GID())
		}
		cellIndex := int(iter.GetIndex())
		rect, err := ld.info.TileRect(cellIndex)
		if err != nil {
			return errors.Wrap(err, "invalid tile")
		}

		opts := &ebiten.DrawImageOptions{
			SourceRect: tse.rect,
			ColorM:     ld.info.color,
			GeoM:       calcGeoM(tile, rect),
			Filter:     ebiten.FilterNearest,
		}

		srcImage, exists := ld.resources.images[tse.source]
		if !exists {
			return errors.Errorf("image source '%v' does not exist", tse)
		}

		ld.image.DrawImage(srcImage, opts)
	}
	return errors.Wrap(iter.Error(), "unable to iterate through layer")
}

func (ld *tileLayerDrawer) Draw(image *ebiten.Image) error {
	return errors.Wrap(image.DrawImage(ld.image, ld.opts), "unable to draw layer")
}
