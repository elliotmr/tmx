package pixeltmx

import (
	"path/filepath"

	"github.com/faiface/pixel"
	"github.com/pkg/errors"
)

type imageLayerDrawer struct {
	resources *Resources
	source    string
	info      *LayerInfo
	sprite    *pixel.Sprite
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

func (ild *imageLayerDrawer) Update() error {
	pic := ild.resources.images[ild.source]
	ild.sprite = pixel.NewSprite(pic, pic.Bounds())
	return nil
}

func (ild *imageLayerDrawer) Draw(t pixel.Target) {
	vec := ild.info.TMXToPixelRect(
		ild.info.offX,
		ild.info.offY,
		ild.sprite.Frame().W(),
		ild.sprite.Frame().H(),
	).Center()
	ild.sprite.DrawColorMask(t, pixel.IM.Moved(vec), ild.info.color)
}
