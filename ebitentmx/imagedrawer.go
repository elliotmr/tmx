package ebitentmx

import (
	"path/filepath"

	"github.com/hajimehoshi/ebiten"
	"github.com/pkg/errors"
	"image"
)

type imageLayerDrawer struct {
	resources *Resources
	source    string
	info      *LayerInfo
	opts      *ebiten.DrawImageOptions
	image     *ebiten.Image
}

func newImageLayerDriver(resources *Resources, info *LayerInfo) (*imageLayerDrawer, error) {
	if info.layer.Image == nil {
		return nil, errors.New("image not set for image layer")
	}

	ild := &imageLayerDrawer{
		resources: resources,
		info:      info,
	}

	if filepath.IsAbs(info.layer.Image.Source) {
		ild.source = filepath.Clean(info.layer.Image.Source)
	} else {
		ild.source = filepath.Join(resources.path, info.layer.Image.Source)
	}



	return ild, ild.Update()
}

func (ild *imageLayerDrawer) Type() int {
	return ImageLayerDrawer
}

func (ild *imageLayerDrawer) Info() *LayerInfo {
	return ild.info
}

func (ild *imageLayerDrawer) Bounds() image.Rectangle {
	return ild.image.Bounds()
}

func (ild *imageLayerDrawer) Update() error {
	var exists bool
	ild.image, exists = ild.resources.images[ild.source]
	if !exists {
		return errors.Errorf("image source '%s' not found", ild.source)
	}

	bounds := ild.image.Bounds()
	geom := ebiten.GeoM{}
	geom.Translate(ild.info.offX, ild.info.offY)
	ild.opts = &ebiten.DrawImageOptions{
		SourceRect: &bounds,
		GeoM:       geom,
		Filter:     ebiten.FilterNearest,
	}
	return nil
}

func (ild *imageLayerDrawer) Draw(image *ebiten.Image) error {
	return errors.Wrap(image.DrawImage(ild.image, ild.opts), "unable to draw image layer")
}
